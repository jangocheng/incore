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

package inapi

import (
	"errors"
	"regexp"
	"strings"
	"sync"

	"github.com/lessos/lessgo/types"
)

var (
	pod_op_mu sync.RWMutex
	pod_st_mu sync.RWMutex
	PodIdReg  = regexp.MustCompile("^[a-f0-9]{12,20}$")
)

type PodSpecBoxImageDriver string

const (
	PodSpecBoxImageDocker PodSpecBoxImageDriver = "docker"
	PodSpecBoxImageRockit PodSpecBoxImageDriver = "rkt"
)

const (
	SpecStatusActive  = "active"
	SpecStatusSuspend = "suspend"
)

type OpType string

const (
	OpStatusPending   OpType = "pending"
	OpStatusRunning   OpType = "running"
	OpStatusStopped   OpType = "stopped"
	OpStatusFailed    OpType = "failed"
	OpStatusDestroyed OpType = "destroyed"
	OpStatusUnknown   OpType = "unknown"
)

// Pod is a collection of boxes, used as either input (create, update) or as output (list, get).
type Pod struct {
	types.TypeMeta `json:",inline"`
	Meta           types.InnerObjectMeta `json:"meta,omitempty"`

	// Spec defines the behavior of a pod.
	Spec *PodSpecBound `json:"spec,omitempty"`

	//
	Operate PodOperate `json:"operate,omitempty"`

	// Apps represents the information about a collection of applications to deploy.
	// this is a module for App Engine
	Apps AppInstances `json:"apps,omitempty"`

	// Status represents the current information about a pod. This data may not be up
	// to date.
	Status *PodStatus `json:"status,omitempty"`

	//
	Payment *PodPayment `json:"payment,omitempty"`
}

type PodPayment struct {
	TimeStart uint32  `json:"time_start"`
	TimeClose uint32  `json:"time_close"`
	Prepay    float64 `json:"prepay"`
	Payout    float64 `json:"payout"`
}

func (pod *Pod) AppServicePorts() ServicePorts {

	var ports ServicePorts

	for _, va := range pod.Apps {

		for _, sp := range va.Spec.ServicePorts {
			ports.Sync(*sp)
		}
	}

	return ports
}

func (pod *Pod) OpRepName() string {
	if pod.Operate.Replica == nil {
		return NsZonePodOpRepKey(pod.Meta.ID, 0)
	}
	return NsZonePodOpRepKey(pod.Meta.ID, pod.Operate.Replica.Id)
}

// PodList is a list of Pods.
type PodList struct {
	types.TypeMeta `json:",inline"`
	Items          []Pod `json:"items"`
}

// PodSpecBound is a description of a bound spec based on PodSpecPlan
type PodSpecBound struct {
	Ref     ObjectReference    `json:"ref,omitempty"`
	Zone    string             `json:"zone,omitempty"`
	Cell    string             `json:"cell,omitempty"`
	Labels  types.Labels       `json:"labels,omitempty"`
	Volumes []PodSpecResVolume `json:"volumes,omitempty"`
	Boxes   []PodSpecBoxBound  `json:"boxes,omitempty"`
}

func (obj *PodSpecBound) Volume(name string) *PodSpecResVolume {

	for _, v := range obj.Volumes {

		if v.Name == name {
			return &v
		}
	}

	return nil
}

func (obj *PodSpecBound) ResComputeBound() *PodSpecBoxResComputeBound {

	rs := &PodSpecBoxResComputeBound{}

	for _, v := range obj.Boxes {

		if v.Resources == nil {
			continue
		}

		rs.CpuLimit += v.Resources.CpuLimit
		rs.MemLimit += v.Resources.MemLimit
	}

	return rs
}

type PodSpecBoundList struct {
	types.TypeMeta `json:",inline"`
	Items          []PodSpecBound `json:"items,omitempty"`
}

// ObjectReference contains enough information to let you inspect or modify the referred object
type ObjectReference struct {
	Id      string `json:"id,omitempty"`
	Name    string `json:"name,omitempty"`
	Version string `json:"version,omitempty"`
	Title   string `json:"title,omitempty"`
}

type PodSpecResVolume struct {
	Ref       ObjectReference `json:"ref,omitempty"`
	Name      string          `json:"name"`
	Labels    types.Labels    `json:"labels,omitempty"`
	SizeLimit int64           `json:"size_limit,omitempty"`
}

type PodSpecBoxBound struct {
	Name      string                     `json:"name"`
	Image     PodSpecBoxImageBound       `json:"image,omitempty"`
	Resources *PodSpecBoxResComputeBound `json:"resources,omitempty"`
	Mounts    VolumeMounts               `json:"mounts,omitempty"`
	Ports     Ports                      `json:"ports,omitempty"`
	Command   []string                   `json:"command,omitempty"`
	Updated   types.MetaTime             `json:"updated,omitempty"`
}

type PodSpecBoxImageBound struct {
	Ref    *ObjectReference      `json:"ref,omitempty"`
	Driver PodSpecBoxImageDriver `json:"driver,omitempty"`

	Options types.Labels `json:"options,omitempty"`

	// Distribution short name of the operating system.
	//  ex: el6, el7, deb7, ubu1404, ...
	OsDist string `json:"os_dist,omitempty"`

	// A human-readable description of the operating system.
	// OsName string `json:"os_name,omitempty"`

	// Architecture indicates the type of hardware.
	//  ex: amd64, armv6l, ...
	Arch string `json:"arch,omitempty"`
}

type PodSpecBoxResComputeBound struct {
	Ref      *ObjectReference `json:"ref,omitempty"`
	CpuLimit int64            `json:"cpu_limit,omitempty"`
	MemLimit int64            `json:"mem_limit,omitempty"`
}

type PodSpecBoxImage struct {
	types.TypeMeta `json:",inline"`
	Meta           types.InnerObjectMeta `json:"meta,omitempty"`

	Status string `json:"status,omitempty"`

	// Container type of the image.
	//  ex: docker, rkt, ...
	Driver PodSpecBoxImageDriver `json:"driver,omitempty"`

	// TODO
	AccessRoles string `json:"access_roles,omitempty"`

	// Options are name value pairs that representing extensional information,
	// usually be used in special system components, names must be unique within the list.
	// ex:
	//  {name: "docker/image/name", value: "centos/lastest"},
	//  {name: "example.com/spec/name", value: "hello"}, ...
	Options types.Labels `json:"options,omitempty"`

	// Annotations are name value pairs that representing additional information,
	// any extra metadata you wish may be added to the list.
	// ex:
	//  {name: "homepage", value: "http://example.com"}, ...
	Annotations types.Labels `json:"annotations,omitempty"`

	// Type name of the operating system.
	//  ex: linux, freebsd, darwin, ...
	OsType string `json:"os_type,omitempty"`

	// Distribution short name of the operating system.
	//  ex: el6, el7, deb7, ubu1404, ...
	OsDist string `json:"os_dist,omitempty"`

	// Version of the operating system.
	OsVersion string `json:"os_version,omitempty"`

	// A human-readable description of the operating system.
	OsName string `json:"os_name,omitempty"`

	// Architecture indicates the type of hardware.
	//  ex: amd64, armv6l, ...
	Arch string `json:"arch,omitempty"`
}

type PodSpecBoxImageList struct {
	types.TypeMeta `json:",inline"`
	Items          []PodSpecBoxImage `json:"items,omitempty"`
}

type PodSpecResourceCompute struct {
	types.TypeMeta `json:",inline"`
	Meta           types.InnerObjectMeta `json:"meta,omitempty"`

	Status string `json:"status,omitempty"`

	// Labels are name value pairs that representing extensional information,
	// usually be used in special system components, names must be unique within the list.
	// ex:
	//  {name: "plan/name", value: "general"}, ...
	Labels types.Labels `json:"labels,omitempty"`

	// CPU, in cores. (500 = .5 cores)
	CpuLimit int64 `json:"cpu_limit,omitempty"`

	// Memory, in bytes.
	MemoryLimit int64 `json:"memory_limit,omitempty"`
}

type PodSpecResourceComputes []*PodSpecResourceCompute

type PodSpecResourceComputeList struct {
	types.TypeMeta `json:",inline"`
	Items          PodSpecResourceComputes `json:"items,omitempty"`
}

func (s PodSpecResourceComputes) Len() int {
	return len(s)
}

func (s PodSpecResourceComputes) Less(i, j int) bool {

	if s[i].CpuLimit < s[j].CpuLimit {
		return true
	}

	return false
}

func (s PodSpecResourceComputes) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

type PodSpecResourceComputeCharge struct {
	Type   uint8   `json:"type"`
	Cycle  uint64  `json:"cycle"`  // default to 3600 seconds
	Cpu    float64 `json:"cpu"`    // value in Core
	Memory float64 `json:"memory"` // value in MiB
}

type PodSpecResourceVolume struct {
	types.TypeMeta `json:",inline"`
	Meta           types.InnerObjectMeta `json:"meta,omitempty"`

	Status string `json:"status,omitempty"`

	// Labels are name value pairs that representing extensional information,
	// usually be used in special system components, names must be unique within the list.
	// ex:
	//  {name: "plan/name", value: "general"}, ...
	Labels types.Labels `json:"labels,omitempty"`

	// Volume size, in bytes.
	Limit   int64 `json:"limit,omitempty"`   // max to 2000GB
	Request int64 `json:"request,omitempty"` // start from 100MB
	Step    int64 `json:"step,omitempty"`    // every step by 100MB
	Default int64 `json:"default,omitempty"` // default to 100MB
}

type PodSpecResourceVolumeList struct {
	types.TypeMeta `json:",inline"`
	Items          []PodSpecResourceVolume `json:"items,omitempty"`
}

type PodSpecResourceVolumeCharge struct {
	Type    uint8   `json:"type"`
	Cycle   uint64  `json:"cycle"`    // default to 3600 seconds
	CapSize float64 `json:"cap_size"` // value in MiB
}

// TODO
type PodSpecResourceNetwork struct {
	types.TypeMeta `json:",inline"`
	Meta           types.InnerObjectMeta `json:"meta,omitempty"`
}

// TODO
type PodSpecResourceNetworkList struct {
	types.TypeMeta `json:",inline"`
	Items          []PodSpecResourceNetwork `json:"items,omitempty"`
}

type PodSpecPlanResourceCharge struct {
	Type  uint8  `json:"type"`
	Cycle uint64 `json:"cycle"` // default to 3600 seconds
}

type PodSpecPlan struct {
	types.TypeMeta `json:",inline"`
	Meta           types.InnerObjectMeta `json:"meta,omitempty"`

	Status string            `json:"status,omitempty"`
	Zones  []PodSpecPlanZone `json:"zones,omitempty"`

	Labels      types.Labels `json:"labels,omitempty"`
	Annotations types.Labels `json:"annotations,omitempty"`

	Images       []PodSpecBoxImage `json:"images,omitempty"`
	ImageDefault string            `json:"image_default,omitempty"`

	ResourceComputes       PodSpecResourceComputes      `json:"res_computes,omitempty"`
	ResourceComputeDefault string                       `json:"res_compute_default,omitempty"`
	ResourceComputeCharge  PodSpecResourceComputeCharge `json:"res_compute_charge,omitempty"`

	ResourceVolumes       []PodSpecResourceVolume     `json:"res_volumes,omitempty"`
	ResourceVolumeDefault string                      `json:"res_volume_default,omitempty"`
	ResourceVolumeCharge  PodSpecResourceVolumeCharge `json:"res_volume_charge,omitempty"`

	ResourceCharge PodSpecPlanResourceCharge `json:"res_charge"`
}

func (s *PodSpecPlan) ChargeFix() {

	s.ResourceCharge.Cycle = 3600

	s.ResourceComputeCharge.Cpu = 0.1
	s.ResourceComputeCharge.Memory = 0.0001

	s.ResourceVolumeCharge.CapSize = 0.000004
}

func (s PodSpecPlan) Image(id string) *PodSpecBoxImage {

	for _, v := range s.Images {

		if v.Meta.ID == id {
			return &v
		}
	}

	return nil
}

func (s PodSpecPlan) ResCompute(id string) *PodSpecResourceCompute {

	for _, v := range s.ResourceComputes {

		if v.Meta.ID == id {
			return v
		}
	}

	return nil
}

func (s PodSpecPlan) ResVolume(id string) *PodSpecResourceVolume {

	for _, v := range s.ResourceVolumes {

		if v.Meta.ID == id {
			return &v
		}
	}

	return nil
}

type PodSpecPlanList struct {
	types.TypeMeta `json:",inline"`
	Items          []PodSpecPlan `json:"items,omitempty"`
}

type PodSpecPlanZone struct {
	Name  string   `json:"name,omitempty"`
	Cells []string `json:"cells,omitempty"`
}

type PodSpecPlanSetup struct {
	types.TypeMeta     `json:",inline"`
	Pod                string                `json:"pod,omitempty"`
	Name               string                `json:"name"`
	Plan               string                `json:"plan"`
	Zone               string                `json:"zone"`
	Cell               string                `json:"cell"`
	ResourceVolume     string                `json:"res_volume"`
	ResourceVolumeSize int64                 `json:"res_volume_size"`
	Boxes              []PodSpecPlanSetupBox `json:"boxes"`
}

type PodSpecPlanSetupBox struct {
	Name                    string `json:"name"`
	Image                   string `json:"image"`
	ResourceCompute         string `json:"res_compute"`
	ResourceComputeCpuLimit int64  `json:"res_compute_cpu_limit,omitempty"`
	ResourceComputeMemLimit int64  `json:"res_compute_mem_limit,omitempty"`
}

func (s *PodSpecPlanSetup) Valid(plan PodSpecPlan) error {

	if s.Name == "" {
		return errors.New("No Name Found")
	}

	hit := false

	for _, zone := range plan.Zones {

		if zone.Name != s.Zone {
			continue
		}

		for _, cell := range zone.Cells {

			if cell == s.Cell {
				hit = true
				break
			}
		}

		break
	}

	if !hit {
		return errors.New("No Zone/Cell Found")
	}
	hit = false

	if s.ResourceVolumeSize < 100*ByteMB {
		return errors.New("Invalid ResourceVolumeSize")
	}

	if s.ResourceVolumeSize >= ByteGB {

		if fix := s.ResourceVolumeSize % ByteGB; fix > 0 {
			s.ResourceVolumeSize -= fix
		}

	} else if fix := s.ResourceVolumeSize % ByteMB; fix > 0 {
		s.ResourceVolumeSize -= fix
	}

	for _, vol := range plan.ResourceVolumes {

		if vol.Meta.ID != s.ResourceVolume {
			continue
		}

		if s.ResourceVolumeSize > vol.Limit ||
			s.ResourceVolumeSize < vol.Request {
			return errors.New("Invalid ResourceVolumeSize")
		}

		hit = true
		break
	}

	if !hit {
		return errors.New("No ResourceVolume Found")
	}

	for i, box := range s.Boxes {

		hit = false

		for _, rv := range plan.ResourceComputes {

			if rv.Meta.ID != box.ResourceCompute {
				continue
			}

			s.Boxes[i].ResourceComputeCpuLimit = rv.CpuLimit
			s.Boxes[i].ResourceComputeMemLimit = rv.MemoryLimit

			hit = true
			break
		}

		if !hit {
			return errors.New("Invalid ResourceCompute")
		}
	}

	return nil
}

type PodOperate struct {
	Action     uint32             `json:"action,omitempty"`
	Version    uint32             `json:"version,omitempty"`
	ReplicaCap int                `json:"replica_cap,omitempty"`
	Replicas   PodOperateReplicas `json:"replicas,omitempty"`
	Replica    *PodOperateReplica `json:"replica,omitempty"`
}

type PodOperateReplica struct {
	Id    uint16       `json:"id"`
	Node  string       `json:"node,omitempty"`
	Ports ServicePorts `json:"ports,omitempty"`
}

type PodOperateReplicas []*PodOperateReplica

func (ls *PodOperateReplicas) CapacitySet(n int) {

	pod_op_mu.Lock()
	defer pod_op_mu.Unlock()

	if len(*ls) > 0 {
		return
	}

	if n < 1 {
		n = 1
	} else if n > 4096 {
		n = 4096
	}

	for id := uint16(0); id < uint16(n); id++ {
		*ls = append(*ls, &PodOperateReplica{
			Id:   id,
			Node: "",
		})
	}
}

func (ls *PodOperateReplicas) Set(set PodOperateReplica) error {

	pod_op_mu.Lock()
	defer pod_op_mu.Unlock()

	for _, v := range *ls {
		if v.Id == set.Id {
			if set.Node != "" && set.Node != v.Node {
				v.Node = set.Node
			}
			return nil
		}
	}

	return errors.New("No Replica Found")
}

func (ls *PodOperateReplicas) Get(rep_id uint16) *PodOperateReplica {

	pod_op_mu.RLock()
	defer pod_op_mu.RUnlock()

	for _, v := range *ls {
		if v.Id == rep_id {
			return v
		}
	}

	return nil
}

// PodStatus represents information about the status of a pod. Status may trail the actual
// state of a system.
type PodStatus struct {
	types.TypeMeta `json:",inline"`
	Phase          OpType            `json:"phase,omitempty"`
	Replicas       PodStatusReplicas `json:"replicas,omitempty"`
	Updated        types.MetaTime    `json:"updated,omitempty"`
}

func (it *Pod) StatusRefresh() {
	it.Status.Refresh(it.Operate.ReplicaCap)
}

func (it *PodStatus) Refresh(rep_cap int) {

	if rep_cap != len(it.Replicas) {

		it.Phase = OpStatusPending

	} else {

		s_diff := map[OpType]int{}

		for _, rep := range it.Replicas {

			if rep.Updated > it.Updated {
				it.Updated = rep.Updated
			}

			switch rep.Phase {

			case OpStatusRunning,
				OpStatusStopped,
				OpStatusFailed,
				OpStatusDestroyed:
				s_diff[rep.Phase]++

			default:
				s_diff[OpStatusPending]++
			}
		}

		if len(s_diff) == 1 {

			for k := range s_diff {
				it.Phase = k
				break
			}

		} else if _, ok := s_diff[OpStatusFailed]; ok {
			it.Phase = OpStatusFailed
		} else {
			it.Phase = OpStatusPending
		}
	}
}

type PodStatusList struct {
	types.TypeMeta `json:",inline"`
	Items          []PodStatus `json:"items"`
}

type PodStatusReplica struct {
	Id      uint16         `json:"id"`
	Phase   OpType         `json:"phase,omitempty"`
	Boxes   []PodBoxStatus `json:"boxes,omitempty"`
	Operate OpStatus       `json:"operate"`
	Updated types.MetaTime `json:"updated,omitempty"`
}

type PodStatusReplicas []*PodStatusReplica

func (ls *PodStatusReplicas) Set(set PodStatusReplica) {

	pod_st_mu.Lock()
	defer pod_st_mu.Unlock()

	for _, v := range *ls {
		if v.Id == set.Id {
			v = &set
			return
		}
	}

	*ls = append(*ls, &set)
}

type PodExecutorStatus struct {
	types.TypeMeta `json:",inline"`
	Items          ExecutorStatuses `json:"items"`
}

// PodBoxStatus represents the current information about a box
type PodBoxStatus struct {
	Name      string                 `json:"name,omitempty"`
	Image     PodBoxStatusImage      `json:"image,omitempty"`
	Resources PodBoxStatusResCompute `json:"resources,omitempty"`
	Mounts    VolumeMounts           `json:"mounts,omitempty"`
	Ports     ServicePorts           `json:"ports,omitempty"`
	Command   []string               `json:"command,omitempty"`
	Executors PodBoxStatusExecutors  `json:"executors,omitempty"`
	Phase     OpType                 `json:"phase,omitempty"`
	Started   types.MetaTime         `json:"started,omitempty"`
	Updated   types.MetaTime         `json:"updated,omitempty"`
}

type PodBoxStatusResCompute struct {
	CpuLimit int64 `json:"cpu_limit,omitempty"`
	MemLimit int64 `json:"mem_limit,omitempty"`
}

type PodBoxStatusImage struct {
	Driver  PodSpecBoxImageDriver `json:"driver,omitempty"`
	Options types.Labels          `json:"options,omitempty"`
}

type PodBoxStatusExecutors []PodBoxStatusExecutor

type PodBoxStatusExecutor struct {
	Name         string `json:"name,omitempty"`
	Phase        OpType `json:"phase,omitempty"`
	Retry        int    `json:"retry,omitempty"`
	ErrorCode    int    `json:"error_code,omitempty"`
	ErrorMessage string `json:"error_message,omitempty"`
}

func (bs *PodBoxStatus) Sync(item PodBoxStatus) (changed bool) {

	changed = false

	if len(item.Name) < 1 {
		return changed
	}

	//
	if bs.Name != item.Name {
		bs.Name = item.Name
		changed = true
	}

	//
	if bs.Started != item.Started {
		bs.Started = item.Started
		changed = true
	}

	//
	if bs.Phase != item.Phase {
		bs.Phase = item.Phase
		changed = true
	}

	//
	if bs.Resources.CpuLimit != item.Resources.CpuLimit ||
		bs.Resources.MemLimit != item.Resources.MemLimit {

		bs.Resources.CpuLimit = item.Resources.CpuLimit
		bs.Resources.MemLimit = item.Resources.MemLimit
		changed = true
	}

	//
	if !bs.Ports.Equal(item.Ports) {
		bs.Ports = item.Ports
		changed = true
	}

	//
	if bs.Image.Driver != item.Image.Driver {
		bs.Image.Driver = item.Image.Driver
		changed = true
	}
	if !bs.Image.Options.Equal(item.Image.Options) {
		bs.Image.Options = item.Image.Options
		changed = true
	}

	//
	if !bs.Mounts.Equal(item.Mounts) {
		bs.Mounts = item.Mounts
		changed = true
	}

	//
	if len(bs.Command) != len(item.Command) ||
		strings.Join(bs.Command, " ") != strings.Join(item.Command, " ") {
		bs.Command = item.Command
		changed = true
	}

	//
	if changed {
		bs.Updated = types.MetaTimeNow()
	}

	return changed
}
