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

package zonemaster

import (
	"errors"
	"sync"

	"github.com/lessos/lessgo/logger"
	"golang.org/x/net/context"

	"code.hooto.com/lessos/loscore/auth"
	"code.hooto.com/lessos/loscore/config"
	"code.hooto.com/lessos/loscore/data"
	"code.hooto.com/lessos/loscore/losapi"
	"code.hooto.com/lessos/loscore/status"
)

var (
	zm_mu sync.Mutex
)

type ApiZoneMaster struct{}

func (s *ApiZoneMaster) HostStatusSync(
	ctx context.Context,
	opts *losapi.ResHost,
) (*losapi.ResZoneMasterList, error) {

	// fmt.Println("host status sync", opts.Meta.Id, opts.Status.Uptime)
	if !status.IsZoneMasterLeader() {
		return nil, errors.New("No ZoneMasterLeader Found")
	}

	if opts == nil || opts.Meta == nil {
		return nil, errors.New("BadArgs")
	}

	if err := auth.TokenValid(ctx); err != nil {
		return nil, err
	}

	//
	host := status.ZoneHostList.Item(opts.Meta.Id)
	if host == nil || host.Meta == nil {
		return nil, errors.New("BadArgs No Host Found")
	}

	//
	if opts.Spec.PeerLanAddr != "" && opts.Spec.PeerLanAddr != host.Spec.PeerLanAddr {
		zm_host_addr_change(opts, host.Spec.PeerLanAddr)
	}

	//
	if host.SyncStatus(*opts) {
		data.ZoneMaster.PvPut(losapi.NsZoneSysHost(status.ZoneId, opts.Meta.Id), host, nil)
		logger.Printf("info", "zone-master/host %s updated", opts.Meta.Id)
	}

	return &status.ZoneMasterList, nil
}

func zm_host_addr_change(host *losapi.ResHost, addr_prev string) {

	zm_mu.Lock()
	defer zm_mu.Unlock()

	if status.ZoneId == "" {
		return
	}

	for _, v := range status.ZoneHostList.Items {

		if host.Spec.PeerLanAddr == v.Spec.PeerLanAddr &&
			host.Meta.Id != v.Meta.Id {
			return
		}
	}

	if status.Host.Meta.Id == host.Spec.PeerLanAddr {
		status.Host.Spec.PeerLanAddr = host.Spec.PeerLanAddr
	}

	for i, p := range status.Zone.LanAddrs {

		if addr_prev == p {

			status.Zone.LanAddrs[i] = host.Spec.PeerLanAddr

			logger.Printf("warn", "ZoneMaster status.Zone.LanAddrs %s->%s",
				addr_prev, host.Spec.PeerLanAddr)

			break
		}
	}

	//
	masters := []losapi.HostNodeAddress{}
	for i, v := range status.ZoneMasterList.Items {

		if v.Addr == addr_prev {

			status.ZoneMasterList.Items[i].Addr = host.Spec.PeerLanAddr

			logger.Printf("warn", "ZoneMaster status.ZoneMasterList.Items %s->%s",
				addr_prev, host.Spec.PeerLanAddr)
		}

		masters = append(masters, losapi.HostNodeAddress(status.ZoneMasterList.Items[i].Addr))
	}
	if len(masters) > 0 {
		config.Config.Masters = masters
		config.Config.Sync()
	}

	//
	for i, v := range status.ZoneHostList.Items {

		if v.Spec.PeerLanAddr == addr_prev {

			status.ZoneHostList.Items[i].Spec.PeerLanAddr = host.Spec.PeerLanAddr

			logger.Printf("warn", "ZoneMaster status.ZoneHostList.Items %s->%s",
				addr_prev, host.Spec.PeerLanAddr)
		}
	}
	for i, v := range status.LocalZoneMasterList.Items {

		if v.Addr == addr_prev {

			status.LocalZoneMasterList.Items[i].Addr = host.Spec.PeerLanAddr

			logger.Printf("warn", "ZoneMaster status.LocalZoneMasterList.Items %s->%s",
				addr_prev, host.Spec.PeerLanAddr)
			break
		}
	}

	//
	/*
		if host.Spec.PeerLanAddr != addr_prev {

			data.ZoneMaster.PvPut(losapi.NsZoneSysHost(status.ZoneId, host.Meta.Id), host, nil)

			logger.Printf("warn", "ZoneMaster NsZoneSysHost %s->%s",
				addr_prev, host.Spec.PeerLanAddr)
		}
	*/

	if rs := data.ZoneMaster.PvGet(losapi.NsGlobalSysZone(status.ZoneId)); rs.OK() {

		var zone losapi.ResZone
		if err := rs.Decode(&zone); err == nil {

			for i, p := range zone.LanAddrs {

				if addr_prev == p {

					zone.LanAddrs[i] = host.Spec.PeerLanAddr

					logger.Printf("warn", "ZoneMaster NsGlobalSysZone %s->%s",
						addr_prev, host.Spec.PeerLanAddr)

					data.ZoneMaster.PvPut(losapi.NsGlobalSysZone(status.ZoneId), zone, nil)
					data.ZoneMaster.PvPut(losapi.NsZoneSysInfo(status.ZoneId), zone, nil)

					break
				}
			}
		}
	}

	if rs := data.ZoneMaster.PvGet(losapi.NsZoneSysMasterNode(status.ZoneId, host.Meta.Id)); rs.OK() {

		var obj losapi.ResZoneMasterNode
		if err := rs.Decode(&obj); err == nil {

			if obj.Addr == addr_prev {

				data.ZoneMaster.PvPut(losapi.NsZoneSysMasterNode(status.ZoneId, host.Meta.Id), losapi.ResZoneMasterNode{
					Id:     host.Meta.Id,
					Addr:   host.Spec.PeerLanAddr,
					Action: 1,
				}, nil)

				logger.Printf("warn", "ZoneMaster NsZoneSysMasterNode %s->%s",
					addr_prev, host.Spec.PeerLanAddr)
			}
		}
	}

	logger.Printf("warn", "ZoneMaster %s->%s", addr_prev, host.Spec.PeerLanAddr)
}