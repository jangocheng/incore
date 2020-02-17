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

package injob

import (
	"sync"

	"github.com/hooto/hlog4g/hlog"
)

type Job interface {
	Name() string
	Run(ctx *Context) error
}

type JobEntry struct {
	mu      sync.Mutex
	sch     *Schedule
	status  int
	job     Job
	running bool
}

func NewJobEntry(job Job, sch *Schedule) *JobEntry {

	if sch == nil {
		sch = NewSchedule()
	}

	return &JobEntry{
		sch:    sch,
		status: StatusStop,
		job:    job,
	}
}

func (it *JobEntry) exec(ctx *Context) {

	it.mu.Lock()
	if it.running {
		it.mu.Unlock()
		return
	}
	it.running = true
	it.mu.Unlock()

	defer func() {
		it.running = false
	}()

	if err := it.job.Run(ctx); err != nil {
		hlog.Printf("warn", "job %s err %s", it.job.Name(), err.Error())
	}
}

func (it *JobEntry) Schedule() *Schedule {
	return it.sch
}

func (it *JobEntry) Commit() *JobEntry {
	it.status = StatusOK
	return it
}