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

import (
	"strings"
	"time"

	"code.hooto.com/lessos/iam/iamapi"
	"code.hooto.com/lessos/iam/iamclient"
	"code.hooto.com/lynkdb/iomix/skv"
	"github.com/lessos/lessgo/crypto/idhash"
	"github.com/lessos/lessgo/httpsrv"
	"github.com/lessos/lessgo/types"

	"code.hooto.com/lessos/loscore/data"
	"code.hooto.com/lessos/loscore/losapi"
)

type Pod struct {
	*httpsrv.Controller
	us iamapi.UserSession
}

func (c *Pod) Init() int {

	//
	c.us, _ = iamclient.SessionInstance(c.Session)

	if !c.us.IsLogin() {
		c.Response.Out.WriteHeader(401)
		c.RenderJson(types.NewTypeErrorMeta(iamapi.ErrCodeUnauthorized, "Unauthorized"))
		return 1
	}

	return 0
}

func (c Pod) ListAction() {

	var (
		ls losapi.PodList
	)

	defer c.RenderJson(&ls)

	// TODO pager
	var rs *skv.Result
	if zone_id := c.Params.Get("zone_id"); zone_id != "" {
		rs = data.ZoneMaster.PvScan(losapi.NsZonePodInstance(zone_id, ""), "", "", 1000)
	} else {
		rs = data.ZoneMaster.PvScan(losapi.NsGlobalPodInstance(""), "", "", 1000)
	}
	rss := rs.KvList()

	var fields types.ArrayPathTree
	if fns := c.Params.Get("fields"); fns != "" {
		fields.Set(fns)
		fields.Sort()
	}

	for _, v := range rss {

		var pod losapi.Pod

		if err := v.Decode(&pod); err == nil {

			// TOPO
			if pod.Meta.User != c.us.UserName {
				continue
			}

			// pod.Spec.Boxes[0].Image.Ref.Id = "6c34df9b5f760180"
			// data.ZoneMaster.PvSet(losapi.NsGlobalPodInstance(pod.Meta.ID), pod, &losapi.ObjectWriteOptions{
			// 	Force: true,
			// })

			if len(fields) > 0 {

				podfs := losapi.Pod{
					Meta: types.InnerObjectMeta{
						ID: pod.Meta.ID,
					},
				}

				if fields.Has("meta/name") {
					podfs.Meta.Name = pod.Meta.Name
				}

				if fields.Has("meta/updated") {
					podfs.Meta.Updated = pod.Meta.Updated
				}

				if fields.Has("spec") && pod.Spec != nil {
					podfs.Spec = &losapi.PodSpecBound{}

					if fields.Has("spec/ref/name") {
						podfs.Spec.Ref.Name = pod.Spec.Ref.Name
					}

					if fields.Has("spec/zone") {
						podfs.Spec.Zone = pod.Spec.Zone
					}

					if fields.Has("spec/cell") {
						podfs.Spec.Cell = pod.Spec.Cell
					}
				}

				if fields.Has("apps") {

					for _, a := range pod.Apps {

						afs := losapi.AppInstance{}

						if fields.Has("apps/meta/id") {
							afs.Meta.ID = a.Meta.ID
						}

						podfs.Apps = append(podfs.Apps, afs)
					}
				}

				if fields.Has("operate") {

					if fields.Has("operate/action") {
						podfs.Operate.Action = pod.Operate.Action
					}

					if fields.Has("operate/node") {
						podfs.Operate.Node = pod.Operate.Node
					}

					if fields.Has("operate/ports") {
						podfs.Operate.Ports = pod.Operate.Ports
					}
				}

				ls.Items = append(ls.Items, podfs)
			} else {
				ls.Items = append(ls.Items, pod)
			}
		}
	}

	ls.Kind = "PodList"
}

func (c Pod) EntryAction() {

	var (
		id      = c.Params.Get("id")
		fields  = c.Params.Get("fields")
		zone_id = c.Params.Get("zone_id")
		set     losapi.Pod
	)

	defer c.RenderJson(&set)

	if zone_id == "" {
		if obj := data.ZoneMaster.PvGet(losapi.NsGlobalPodInstance(id)); obj.OK() {
			obj.Decode(&set)
		}
	} else {
		if obj := data.ZoneMaster.PvGet(losapi.NsZonePodInstance(zone_id, id)); obj.OK() {
			obj.Decode(&set)
		}
	}
	if set.Meta.ID == "" || set.Meta.User != c.us.UserName {
		set.Error = types.NewErrorMeta("404", "Pod Not Found")
		return
	}

	if fields == "status" {
		if v := pod_status(set.Meta.ID, c.us.UserName); v.Phase != "" {
			set.Status = &v
		} else {
			set.Status = nil
		}
	}

	set.Kind = "Pod"
}

func (c Pod) NewAction() {

	var (
		set  losapi.PodSpecPlanSetup
		spec losapi.PodSpecPlan
	)

	defer c.RenderJson(&set)

	if err := c.Request.JsonDecode(&set); err != nil {
		set.Error = types.NewErrorMeta("400", err.Error())
		return
	}

	set.Name = strings.TrimSpace(set.Name)
	if set.Name == "" {
		set.Error = types.NewErrorMeta("400", "Name can not be null")
		return
	}

	//
	if obj := data.ZoneMaster.PvGet(losapi.NsGlobalPodSpec("plan", set.Plan)); obj.OK() {
		obj.Decode(&spec)
	}
	if spec.Meta.ID == "" || spec.Meta.ID != set.Plan {
		set.Error = types.NewErrorMeta("400", "Spec Not Found")
		return
	}

	//
	if err := set.Valid(spec); err != nil {
		set.Error = types.NewErrorMeta("400", err.Error())
		return
	}

	//
	var zone losapi.ResZone
	if obj := data.ZoneMaster.PvGet(losapi.NsGlobalSysZone(set.Zone)); obj.OK() {
		obj.Decode(&zone)
	}
	if zone.Meta.Id == "" {
		set.Error = types.NewErrorMeta("400", "Zone Not Found")
		return
	}

	//
	var cell losapi.ResCell
	if obj := data.ZoneMaster.PvGet(losapi.NsGlobalSysCell(set.Zone, set.Cell)); obj.OK() {
		obj.Decode(&cell)
	}
	if cell.Meta.Id == "" {
		set.Error = types.NewErrorMeta("400", "Cell Not Found")
		return
	}

	res_vol := spec.ResVolume(set.ResourceVolume)
	if res_vol == nil {
		set.Error = types.NewErrorMeta("400", "No ResourceVolume Found")
		return
	}

	pod := losapi.Pod{
		Meta: types.InnerObjectMeta{
			ID:      idhash.RandHexString(16),
			Name:    set.Name,
			User:    c.us.UserName,
			Created: types.MetaTimeNow(),
			Updated: types.MetaTimeNow(),
		},
		Spec: &losapi.PodSpecBound{
			Ref: losapi.ObjectReference{
				Id:      spec.Meta.ID,
				Name:    spec.Meta.Name,
				Version: spec.Meta.Version,
			},
			Zone:   set.Zone,
			Cell:   set.Cell,
			Labels: spec.Labels,
			Volumes: []losapi.PodSpecResVolume{
				{
					Ref: losapi.ObjectReference{
						Name:    res_vol.Meta.Name,
						Version: res_vol.Meta.Version,
					},
					Name:      "system",
					SizeLimit: set.ResourceVolumeSize,
				},
			},
		},
		Operate: losapi.PodOperate{
			Action: losapi.OpActionStart,
		},
	}

	//
	for _, v := range set.Boxes {

		img := spec.Image(v.Image)
		if img == nil {
			set.Error = types.NewErrorMeta("400", "No Image Found")
			return
		}

		res := spec.ResCompute(v.ResourceCompute)
		if res == nil {
			set.Error = types.NewErrorMeta("400", "No ResourceCompute Found")
			return
		}

		pod.Spec.Boxes = append(pod.Spec.Boxes, losapi.PodSpecBoxBound{
			Name:    v.Name,
			Updated: types.MetaTimeNow(),
			Image: losapi.PodSpecBoxImageBound{
				Ref: &losapi.ObjectReference{
					Id:      img.Meta.ID,
					Name:    img.Meta.Name,
					Version: img.Meta.Version,
				},
				Driver:  img.Driver,
				OsDist:  img.OsDist,
				Arch:    img.Arch,
				Options: img.Options,
			},
			Resources: &losapi.PodSpecBoxResComputeBound{
				Ref: &losapi.ObjectReference{
					Id:      res.Meta.ID,
					Name:    res.Meta.Name,
					Version: res.Meta.Version,
				},
				CpuLimit: v.ResourceComputeCpuLimit,
				MemLimit: v.ResourceComputeMemLimit,
			},
		})
	}

	//
	if rs := data.ZoneMaster.PvNew(losapi.NsGlobalPodInstance(pod.Meta.ID), pod, nil); !rs.OK() {
		set.Error = types.NewErrorMeta("500", rs.Bytex().String())
		return
	}

	// Pod Map to Cell Queue
	qmpath := losapi.NsZonePodSetQueue(pod.Spec.Zone, pod.Spec.Cell, pod.Meta.ID)
	data.ZoneMaster.PvNew(qmpath, pod, nil)

	set.Kind = "PodInstance"
}

func (c Pod) SetInfoAction() {

	var (
		set          losapi.Pod
		prev         losapi.Pod
		prev_version uint64
	)

	defer c.RenderJson(&set)

	if err := c.Request.JsonDecode(&set); err != nil {
		set.Error = types.NewErrorMeta("400", err.Error())
		return
	}

	if rs := data.ZoneMaster.PvGet(losapi.NsGlobalPodInstance(set.Meta.ID)); !rs.OK() {

		set.Error = types.NewErrorMeta("400", "Prev Pod Not Found")
		return
	} else {
		rs.Decode(&prev)
		prev_version = rs.Meta().Version
	}

	if prev.Meta.ID != set.Meta.ID {
		set.Error = types.NewErrorMeta("400", "Prev Pod Not Found")
		return
	}

	if prev.Meta.User != c.us.UserName {
		set.Error = types.NewErrorMeta(iamapi.ErrCodeAccessDenied, "Access Denied")
		return
	}

	//
	if prev.Meta.Name == set.Meta.Name &&
		prev.Operate.Action == set.Operate.Action {
		set.Kind = "PodInstance"
		return
	}

	prev.Meta.Name = set.Meta.Name
	prev.Operate.Action = set.Operate.Action

	//
	prev.Meta.Updated = types.MetaTimeNow()

	data.ZoneMaster.PvPut(losapi.NsGlobalPodInstance(prev.Meta.ID), prev, &skv.PathWriteOptions{
		PrevVersion: prev_version,
	})

	// Pod Map to Cell Queue
	qstr := losapi.NsZonePodSetQueue(prev.Spec.Zone, prev.Spec.Cell, prev.Meta.ID)

	if rs := data.ZoneMaster.PvGet(qstr); rs.OK() {
		set.Error = types.NewErrorMeta(losapi.ErrCodeBadArgument, "ObjectAlreadyExists")
		return
	}

	data.ZoneMaster.PvPut(qstr, prev, &skv.PathWriteOptions{
		Force: true,
	})

	set.Kind = "PodInstance"
}

func (c Pod) StatusAction() {

	rsp := losapi.PodStatus{}

	defer c.RenderJson(&rsp)

	rsp = pod_status(c.Params.Get("id"), c.us.UserName)
}

/*
func (c Pod) StatusListAction() {

	ls := losapi.PodStatusList{}
	defer c.RenderJson(&ls)

	ids := strings.Split(c.Params.Get("ids"), ",")
	for i, id := range ids {

		if i > 30 {
			break
		}

		if set := pod_status(id, c.us.UserName); set.Error == nil && set.Kind == "PodStatus" {
			ls.Items = append(ls.Items, set)
		}
	}

	ls.Kind = "PodStatusList"
}
*/

func pod_status(pod_id, user_name string) losapi.PodStatus {

	var (
		pod    losapi.Pod
		status losapi.PodStatus
	)

	//
	if obj := data.ZoneMaster.PvGet(losapi.NsGlobalPodInstance(pod_id)); obj.OK() {
		obj.Decode(&pod)
	}

	if pod.Meta.ID == "" {
		status.Error = types.NewErrorMeta("404.01", "Pod Not Found")
		return status
	}

	if user_name != pod.Meta.User {
		status.Error = types.NewErrorMeta(iamapi.ErrCodeAccessDenied, "Access Denied")
		return status
	}

	//
	if obj := data.ZoneMaster.PvGet(
		losapi.NsZonePodInstance(pod.Spec.Zone, pod_id),
	); obj.OK() {
		obj.Decode(&pod)
	}

	if pod.Meta.ID == "" {
		status.Error = types.NewErrorMeta("404.02", "Pod Not Found")
		return status
	}

	//
	if obj := data.ZoneMaster.PvGet(
		losapi.NsZoneHostBoundPodStatus(pod.Spec.Zone, pod.Operate.Node, pod.Meta.ID),
	); obj.OK() {
		obj.Decode(&status)
	}

	if status.Phase == "" {
		status.Error = types.NewErrorMeta("404.03", "Pod Not Found")
	} else {

		if (time.Now().UTC().Unix() - status.Updated.Time().Unix()) > 600 {
			status.Phase = losapi.OpStatusUnknown
		}

		status.Kind = "PodStatus"
	}

	return status
}