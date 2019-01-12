// Copyright 2018 Eryx <evorui аt gmail dοt com>, All rights reserved.
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

package filerender

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"text/template"
)

func RenderString(src string, dstFile string, perm os.FileMode, sets map[string]interface{}) error {

	//
	tpl, err := template.New("s").Parse(src)
	if err != nil {
		return err
	}

	var bsdst bytes.Buffer
	if err := tpl.Execute(&bsdst, sets); err != nil {
		return err
	}

	if _, err := os.Stat(dstFile); err != nil {
		os.MkdirAll(filepath.Dir(dstFile), 0755)
	}
	fpdst, err := os.OpenFile(dstFile, os.O_RDWR|os.O_CREATE, perm)
	if err != nil {
		return err
	}
	defer fpdst.Close()

	fpdst.Seek(0, 0)
	fpdst.Truncate(0)

	_, err = fpdst.Write(bsdst.Bytes())

	return err
}

func Render(srcFile, dstFile string, perm os.FileMode, sets map[string]interface{}) error {

	fpsrc, err := os.Open(srcFile)
	if err != nil {
		return err
	}
	defer fpsrc.Close()

	src, err := ioutil.ReadAll(fpsrc)
	if err != nil {
		return err
	}

	return RenderString(string(src), dstFile, perm, sets)
}
