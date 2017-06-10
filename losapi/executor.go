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

package losapi

import (
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/lessos/lessgo/types"
)

var (
	executor_mu sync.Mutex
)

type ExecutorSetup struct {
	types.TypeMeta `json:",inline"`
	Spec           types.NameIdentifier `json:"spec,omitempty"`
	PodId          string               `json:"pod_id,omitempty"`
	AppId          string               `json:"app_id,omitempty"`
	Executor       Executor             `json:"executor,omitempty"`
}

type SpecExecutor struct {
	types.TypeMeta `json:",inline"`
	Meta           types.InnerObjectMeta `json:"meta"`
	Description    string                `json:"description,omitempty"`
	Labels         types.Labels          `json:"labels,omitempty"`
	Executor       Executor              `json:"executor,omitempty"`
}

type SpecExecutorList struct {
	types.TypeMeta `json:",inline"`
	Items          []SpecExecutor `json:"items,omitempty"`
}

type Executor struct {
	Name      types.NameIdentifier `json:"name"`
	Updated   types.MetaTime       `json:"updated,omitempty"`
	Vendor    string               `json:"vendor,omitempty"`
	Dir       string               `json:"dir,omitempty"`   // /home/action/apps/demo
	User      string               `json:"user,omitempty"`  // default: action
	Group     string               `json:"group,omitempty"` // default: action
	ExecStart string               `json:"exec_start,omitempty"`
	ExecStop  string               `json:"exec_stop,omitempty"`
	Priority  uint8                `json:"priority,omitempty"`
	Plan      ExecPlanner          `json:"plan,omitempty"`
	Options   types.Labels         `json:"options,omitempty"`
	Status    *ExecutorStatus      `json:"status,omitempty"`
}

type ExecutorList struct {
	types.TypeMeta `json:",inline"`
	Items          Executors `json:"items,omitempty"`
}

type Executors []Executor

func (ls *Executors) Sync(item2 Executor) {

	executor_mu.Lock()
	defer executor_mu.Unlock()

	for i, v := range *ls {

		if v.Name == item2.Name {

			if item2.Updated > v.Updated {
				(*ls)[i] = item2
				if v.Status != nil {
					(*ls)[i].Status = v.Status
				}
			}

			return
		}
	}

	*ls = append((*ls), item2)
}

func (ls *Executors) Remove(name types.NameIdentifier) {

	executor_mu.Lock()
	defer executor_mu.Unlock()

	for i, v := range *ls {

		if v.Name == name {
			*ls = append((*ls)[0:i], (*ls)[i+1:]...)
			break
		}
	}
}

type ExecutorAction uint64

const (
	ExecutorActionStart   ExecutorAction = 1 << 1
	ExecutorActionStarted ExecutorAction = 1 << 2
	ExecutorActionStop    ExecutorAction = 1 << 3
	ExecutorActionStopped ExecutorAction = 1 << 4
	ExecutorActionPending ExecutorAction = 1 << 10
	ExecutorActionFailed  ExecutorAction = 1 << 11
)

func (a ExecutorAction) Allow(as ExecutorAction) bool {
	return (as & a) == as
}

func (a *ExecutorAction) Remove(as ExecutorAction) {
	*a = (*a | as) - (as)
}

func (a *ExecutorAction) Append(as ExecutorAction) {
	*a = (*a | as)
}

func (a ExecutorAction) String() string {

	as := []string{}

	if a.Allow(ExecutorActionPending) {
		as = append(as, "pending")
	}

	if a.Allow(ExecutorActionStarted) {
		as = append(as, "started")
	}

	if a.Allow(ExecutorActionStopped) {
		as = append(as, "stopped")
	}

	if a.Allow(ExecutorActionFailed) {
		as = append(as, "failed")
	}

	return strings.Join(as, ",")
}

type ExecPlanTimer string

func (pt ExecPlanTimer) Seconds() int64 {

	if t, err := time.ParseDuration(string(pt)); err == nil {
		return int64(t.Seconds())
	}

	return 0
}

type ExecPlanner struct {
	OnBoot     bool              `json:"on_boot,omitempty"`
	OnTick     uint32            `json:"on_tick,omitempty"`
	OnCalendar *ExecPlanTimer    `json:"on_calendar,omitempty"`
	OnFailed   *ExecPlanOnFailed `json:"on_failed,omitempty"`
}

// ExecPlanOnFailed describes how the executor should be re-executed.
type ExecPlanOnFailed struct {
	RetrySec ExecPlanTimer `json:"retry_sec,omitempty"`
	RetryMax int           `json:"retry_max,omitempty"`
}

type ExecutorStatusPlanner struct {
	Updated          types.MetaTime `json:"updated,omitempty"`
	OnFailedRetryNum int            `json:"on_failed_retry_num,omitempty"`
}

type ExecutorStatus struct {
	Name    types.NameIdentifier  `json:"name"`
	Created types.MetaTime        `json:"created,omitempty"`
	Updated types.MetaTime        `json:"updated,omitempty"`
	Vendor  string                `json:"vendor,omitempty"`
	Action  ExecutorAction        `json:"action,omitempty"`
	Plan    ExecutorStatusPlanner `json:"plan,omitempty"`
	Cmd     *exec.Cmd             `json:"-"`
}

type ExecutorStatuses []*ExecutorStatus

func (es *ExecutorStatuses) Get(name types.NameIdentifier) *ExecutorStatus {

	executor_mu.Lock()
	defer executor_mu.Unlock()

	for _, v := range *es {

		if v.Name == name {
			return v
		}
	}

	return nil
}

func (es *ExecutorStatuses) Sync(item *ExecutorStatus) {

	executor_mu.Lock()
	defer executor_mu.Unlock()

	for i, v := range *es {

		if v.Name == item.Name {
			(*es)[i] = item
			return
		}
	}

	*es = append(*es, item)
}
