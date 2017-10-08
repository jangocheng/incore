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

package ops

import (
	"fmt"
	"strings"

	"github.com/lessos/lessgo/types"

	"github.com/sysinner/incore/data"
	"github.com/sysinner/incore/inapi"
)

func (c Host) ZoneListAction() {

	ls := inapi.GeneralObjectList{}
	defer c.RenderJson(&ls)

	//
	rs := data.ZoneMaster.PvScan(inapi.NsGlobalSysZone(""), "", "", 100).KvList()

	for _, v := range rs {

		var zone inapi.ResZone

		if err := v.Decode(&zone); err != nil || zone.Meta.Id == "" {
			continue
		}

		if c.Params.Get("fields") == "cells" {

			rs2 := data.ZoneMaster.PvScan(inapi.NsGlobalSysCell(zone.Meta.Id, ""), "", "", 100).KvList()

			for _, v2 := range rs2 {

				var cell inapi.ResCell

				if err := v2.Decode(&cell); err == nil {
					zone.Cells = append(zone.Cells, &cell)
				}
			}
		}

		ls.Items = append(ls.Items, zone)
	}

	ls.Kind = "HostZoneList"
}

func (c Host) ZoneEntryAction() {

	var set struct {
		inapi.GeneralObject
		inapi.ResZone
	}

	defer c.RenderJson(&set)

	if obj := data.ZoneMaster.PvGet(inapi.NsGlobalSysZone(c.Params.Get("id"))); obj.OK() {

		if err := obj.Decode(&set.ResZone); err != nil {
			set.Error = &types.ErrorMeta{"400", err.Error()}
		} else {

			if c.Params.Get("fields") == "cells" {

				rs2 := data.ZoneMaster.PvScan(inapi.NsGlobalSysCell(set.Meta.Id, ""), "", "", 100).KvList()

				for _, v2 := range rs2 {

					var cell inapi.ResCell

					if err := v2.Decode(&cell); err == nil {
						set.Cells = append(set.Cells, &cell)
					}
				}
			}

		}
	}

	if set.Meta.Id != "" {
		set.Kind = "HostZone"
	} else {
		set.Error = &types.ErrorMeta{"404", "Item Not Found"}
	}
}

func (c Host) ZoneSetAction() {

	var set struct {
		inapi.GeneralObject
		inapi.ResZone
	}

	defer c.RenderJson(&set)

	if err := c.Request.JsonDecode(&set.ResZone); err != nil {
		set.Error = &types.ErrorMeta{"400", err.Error()}
		return
	}

	set.Meta.Id = strings.ToLower(set.Meta.Id)

	if mat := inapi.ResSysZoneIdReg.MatchString(set.Meta.Id); !mat {
		set.Error = &types.ErrorMeta{
			Code:    "400",
			Message: "Zone Id must consist of letters or numbers, and begin with a letter",
		}
		return
	}

	for _, addr := range set.WanAddrs {

		if !inapi.HostNodeAddress(addr).Valid() {
			set.Error = &types.ErrorMeta{
				Code:    "400",
				Message: fmt.Sprintf("Invalid Address (%s)", addr),
			}
			return
		}
	}

	for _, addr := range set.LanAddrs {

		if !inapi.HostNodeAddress(addr).Valid() {
			set.Error = &types.ErrorMeta{
				Code:    "400",
				Message: fmt.Sprintf("Invalid Address (%s)", addr),
			}
			return
		}
	}

	if obj := data.ZoneMaster.PvGet(inapi.NsGlobalSysZone(set.Meta.Id)); obj.OK() {

		var prev inapi.ResZone
		if err := obj.Decode(&prev); err == nil {

			if prev.Meta.Created > 0 {
				set.Meta.Created = prev.Meta.Created
			}
		}
	}

	//
	/*
		if obj := data.ZoneMaster.PvGet(inapi.NsGlobalDataBucket(set.Meta.Id)); obj.OK() {

			bucket_ns := btapi.BucketNameService{
				Phase: 1,
				Bound: []btapi.BucketBoundZone{
					{
						ZoneId:  set.Meta.Id,
						PermSrv: btapi.OpAll,
					},
				},
			}

			bucket_instance := btapi.Bucket{
				Meta: types.ObjectMeta{
					Id:      set.Meta.Id,
					Name:    set.Meta.Name,
					Created: utilx.TimeNow("atom"),
					Updated: utilx.TimeNow("atom"),
				},
				Phase:   bucket_ns.Phase,
				Bound:   bucket_ns.Bound,
				Summary: "",
			}

			data.ZoneMaster.PvSet("/ns/bt/buckets/"+set.Meta.Id, bucket_ns, nil)
			data.ZoneMaster.PvSet("/global/bt/buckets/instances/"+set.Meta.Id, bucket_instance, nil)

			zone_ns := inapi.ResZoneNameService{
				Phase:     1,
				Endpoints: set.WANEndpoints,
			}
			data.ZoneMaster.PvSet("/ns/bt/zones/"+set.Meta.Id, zone_ns, nil)
		}
	*/

	/*
		if rs := data.ZoneMaster.PvScan(inapi.NsGlobalSysCell(set.Meta.Id), "", "", 100); rs.OK() {

			cell := btapi.HostCell{
				Meta: types.ObjectMeta{
					Id:      utils.StringNewRand(8),
					Name:    "Default Cell",
					Created: utilx.TimeNow("atom"),
					Updated: utilx.TimeNow("atom"),
				},
				ZoneId:      set.Meta.Id,
				Status:      1,
				Description: "Default Cell",
			}

			data.ZoneMaster.PvSet(fmt.Sprintf("/global/bt/cells/%s/%s", set.Meta.Id, cell.Meta.Id), cell, nil)
		}
	*/

	if set.Meta.Created == 0 {
		set.Meta.Created = uint64(types.MetaTimeNow())
	}

	set.Meta.Updated = uint64(types.MetaTimeNow())

	data.ZoneMaster.PvPut(inapi.NsGlobalSysZone(set.Meta.Id), set.ResZone, nil)

	set.Kind = "HostZone"
}
