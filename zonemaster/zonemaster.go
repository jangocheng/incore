// Copyright 2015 Eryx <evorui аt gmаil dοt cοm>, All rights reserved.
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

package zonemaster

import (
	"sync"
	"time"

	"github.com/hooto/hlog4g/hlog"
	iamWorker "github.com/hooto/iam/worker"

	"github.com/sysinner/incore/injob"
	"github.com/sysinner/incore/module/mail_queue"
	"github.com/sysinner/incore/module/pod_status_mail"
	"github.com/sysinner/incore/status"
)

var (
	worker_mu sync.Mutex
	running   = false
	jobDaemon *injob.Daemon
	err       error
)

func Start() error {

	worker_mu.Lock()
	defer worker_mu.Unlock()

	if running {
		return nil
	}

	hlog.Print("info", "zonemaster/worker started")

	jobDaemon, err = injob.NewDaemon(status.JobContextRefresh)
	if err != nil {
		return err
	}

	jobDaemon.Commit(pod_status_mail.NewPodStatusMailJobEntry())
	jobDaemon.Commit(mail_queue.NewMailQueueJobEntry())

	go jobDaemon.Start()

	go func() {

		for {
			time.Sleep(3e9)

			if !status.IsZoneMaster() {
				continue
			}

			zoneTracker()

			zmWorkerMasterLeaderRefresh()

			if status.IsZoneMasterLeader() {

				if err := scheduleAction(); err != nil {
					hlog.Printf("warn", "zm/scheduler err:%s", err.Error())
				}

				if err := podChargeRefresh(); err != nil {
					hlog.Printf("warn", "zm/pod/charge err:%s", err.Error())
				}

				status.ZoneScheduled = time.Now().Unix()

				go iamWorker.AccountChargeCloseRefresh()
				go iamWorker.MsgQueueRefresh()
			}
		}
	}()

	return nil
}
