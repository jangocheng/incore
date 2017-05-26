// Copyright 2015 Authors, All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package v1

/*
#cgo linux CFLAGS: -D__LINUX__
#cgo darwin CFLAGS: -D__DARWIN__
#cgo freebsd CFLAGS: -D__FREEBSD__
#cgo LDFLAGS: -lutil

#include <stdlib.h>
#include <sys/ioctl.h>

#ifdef __LINUX__
#include <pty.h>
#endif

#ifdef __DARWIN__
#include <util.h>
#endif

#ifdef __FREEBSD__
#include <sys/types.h>
#include <termios.h>
#include <libutil.h>
#endif

int goForkpty(int *amaster, struct winsize *winp) {
    return forkpty(amaster, NULL, NULL, winp);
}

int goChangeWinsz(int fd, struct winsize *winp) {
    return ioctl(fd, TIOCSWINSZ, winp);
}
*/
import "C"

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/user"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"unicode/utf8"

	"github.com/lessos/lessgo/deps/go.net/websocket"
	"github.com/lessos/lessgo/httpsrv"
	"github.com/lessos/lessgo/logger"
	"github.com/lessos/lessgo/utils"
)

type Terminal struct {
	*httpsrv.Controller
}

type ApiResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

type TerminalResponse struct {
	ApiResponse
	Data struct {
		Output string `json:"output"`
	} `json:"data"`
}

var (
	connections = 0
)

func readFull(r io.Reader, buf []byte) {
	_, err := io.ReadFull(r, buf)
	if err != nil {
		panic(fmt.Sprintf("Could not read fully from %v: %s", r, err))
	}
}

func readInt(r io.Reader) int {
	var buf [8]byte
	readFull(r, buf[:])
	result, err := strconv.Atoi(strings.TrimSpace(string(buf[:])))
	if err != nil {
		panic(fmt.Sprintf("Could not convert input from %v to int: %s", r, err))
	}
	return result
}

func setColsRows(winsz *C.struct_winsize, cols int, rows int) {
	winsz.ws_row = C.ushort(rows)
	winsz.ws_col = C.ushort(cols)
	winsz.ws_xpixel = C.ushort(cols * 9)
	winsz.ws_ypixel = C.ushort(rows * 16)
}

func redirToWs(fd int, ws *websocket.Conn) {

	defer func() {
		if r := recover(); r != nil {
			//fmt.Fprintf(os.Stderr, "Error occured: %s\n", r)
			runtime.Goexit()
		}
	}()

	var buf [8192]byte
	start, end, buflen := 0, 0, 0
	for {

		switch nr, _ := syscall.Read(fd, buf[start:]); {

		case nr < 0:
			//fmt.Fprintf(os.Stderr, "error reading from websocket %d with code %d\n", fd, er)
			return

		case nr == 0: // EOF
			return

		case nr > 0:
			buflen = start + nr
			for end = buflen - 1; end >= 0; end-- {
				if utf8.RuneStart(buf[end]) {
					ch, width := utf8.DecodeRune(buf[end:buflen])
					if ch != utf8.RuneError {
						end += width
					}
					break
				}

				if buflen-end >= 6 {
					//fmt.Fprintf(os.Stderr, "Invalid UTF-8 sequence in output")
					end = nr
					break
				}
			}

			runes := bytes.Runes(buf[0:end])
			buf_clean := []byte(string(runes))
			//fmt.Println(string(buf_clean), len(string(buf_clean)))

			nw, ew := ws.Write(buf_clean[:])
			if ew != nil {
				//fmt.Fprintf(os.Stderr, "error writing to websocket with code %s\n", ew)
				return
			}

			if nw != len(buf_clean) {
				//fmt.Fprintf(os.Stderr, "Written %d instead of expected %d\n", nw, end)
			}

			start = buflen - end

			if start > 0 {
				// copy remaning read bytes from the end to the beginning of a buffer
				// so that we will get normal bytes
				for i := 0; i < start; i++ {
					buf[i] = buf[end+i]
				}
			}
		}
	}
}

func redirFromWs(fd int, ws *websocket.Conn, pid int, winsz *C.struct_winsize) {

	defer func() {
		if r := recover(); r != nil {
			//fmt.Fprintf(os.Stderr, "Error occured: %s\n", r)
			syscall.Kill(pid, syscall.SIGHUP)
			runtime.Goexit()
		}
	}()

	var buf [2048]byte

	for {
		/*
		   communication protocol:

		   1 byte   cmd

		   if cmd = i // input
		       8 byte        length (ascii)
		       length bytes  the actual input

		   if cmd = w // window size changed
		       8 byte        cols (ascii)
		       8 byte        rows (ascii)
		*/

		readFull(ws, buf[0:1])

		switch buf[0] {

		case 'i':

			length := readInt(ws)

			switch nr, _ := io.ReadFull(ws, buf[0:length]); {

			case nr < 0:
				//fmt.Fprintf(os.Stderr, "error reading from websocket with code %s\n", er)
				return

			case nr == 0: // EOF
				//fmt.Fprintf(os.Stderr, "connection closed, sending SIGHUP to %d\n")
				syscall.Kill(pid, syscall.SIGHUP)
				return

			case nr > 0:
				nw, _ := syscall.Write(fd, buf[0:nr])
				if nw != nr {
					//fmt.Fprintf(os.Stderr, "error writing to fd = %d with code %d\n", fd, ew)
					return
				}
			}

		case 'w':
			cols, rows := readInt(ws), readInt(ws)
			setColsRows(winsz, cols, rows)
			C.goChangeWinsz(C.int(fd), winsz)

		default:
			panic("Unknown command " + string(buf[0]))
		}
	}
}

func TerminalWsOpenAction(wsconn *websocket.Conn) {

	logger.Printf("info", "ws://terminal/open")

	var err error
	var rsp TerminalResponse

	defer func() {
		logger.Printf("info", "ws://terminal/close")
		wsconn.Close()
	}()

	var msg string
	if err := websocket.Message.Receive(wsconn, &msg); err != nil {
		return
	}

	//fmt.Println("msg", msg)
	var req struct {
		AccessToken string `json:"access_token"`
		Data        struct {
			Command string `json:"command"`
		} `json:"data"`
	}
	err = utils.JsonDecode(msg, &req)
	if err != nil {
		return
	}

	//
	// sess := this.Session.Instance(req.AccessToken)
	// if sess.Uid == "0" || sess.Uid == "" {
	// 	rsp.Status = 401
	// 	rsp.Message = "Unauthorized"
	// 	return
	// }
	osuser := "action"
	//fmt.Println("login as user:", osuser)

	//
	u, e := user.Lookup(osuser)
	if e != nil {
		rsp.Status = 401
		rsp.Message = "Unauthorized"
		return
	}
	uuid, _ := strconv.Atoi(u.Uid)
	ugid, _ := strconv.Atoi(u.Gid)
	fmt.Println(uuid, ugid)

	/* syscall.Setgid(ugid)
	   syscall.Setuid(uuid)
	   syscall.Chdir(this.Cfg.Prefix + "/" + sess.Uname)
	*/
	connections++
	defer func() {
		connections--

		if r := recover(); r != nil {
			//fmt.Fprintf(os.Stderr, "Error occured: %s\n", r)
			runtime.Goexit()
		}
	}()

	cols, rows := readInt(wsconn), readInt(wsconn)

	var winsz = new(C.struct_winsize)
	setColsRows(winsz, cols, rows)

	cpttyno := C.int(0)
	pid := int(C.goForkpty(&cpttyno, winsz))
	// pttyno := int(cpttyno)
	// defer syscall.Close(pttyno) // forgot to close on errors too

	// fmt.Println(cols, rows, pid)

	if pid == 0 {

		bashargs := []string{"bash", "--rcfile", "/home/action/.bashrc"}

		basepath, err := exec.LookPath("bash")
		if err != nil {
			//fmt.Fprintf(os.Stderr, "Could not find bash: %s\n", err)
			basepath = "/bin/sh"
			bashargs = []string{"sh"}
		}

		// basepath = "/usr/bin/bash"

		// fmt.Println(basepath, bashargs)

		// syscall.Setgid(ugid)
		// syscall.Setuid(uuid)
		// syscall.Chdir("/opt/lessosu/fa63f52b0afe/")
		// syscall.Chroot("/opt/lessosu/fa63f52b0afe/")
		envs := os.Environ()
		envs = append(envs, "TMOUT=3600")

		err = syscall.Exec(basepath, bashargs, envs)
		// fmt.Println(err)
		panic("unreachable code")
	}

	pttyno := int(cpttyno)

	//fmt.Println("Pid is", pid, " ptty number is", pttyno)
	go redirFromWs(pttyno, wsconn, pid, winsz)
	go redirToWs(pttyno, wsconn)

	syscall.Wait4(pid, nil, 0, nil)
}
