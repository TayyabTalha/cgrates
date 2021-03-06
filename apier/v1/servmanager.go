/*
Real-time Online/Offline Charging System (OCS) for Telecom & ISP environments
Copyright (C) ITsysCOM GmbH

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>
*/

package v1

import (
	"github.com/cgrates/cgrates/servmanager"
	"github.com/cgrates/cgrates/utils"
)

func NewServManagerV1(sm *servmanager.ServiceManager) *ServManagerV1 {
	return &ServManagerV1{sm: sm}
}

type ServManagerV1 struct {
	sm *servmanager.ServiceManager // Need to have them capitalize so we can export in V2
}

func (servManager *ServManagerV1) StartService(args servmanager.ArgStartService, reply *string) (err error) {
	return servManager.sm.V1StartService(args, reply)
}

func (servManager *ServManagerV1) StopService(args servmanager.ArgStartService, reply *string) (err error) {
	return servManager.sm.V1StopService(args, reply)
}

func (servManager *ServManagerV1) ServiceStatus(args servmanager.ArgStartService, reply *string) (err error) {
	return servManager.sm.V1ServiceStatus(args, reply)
}

// Ping return pong if the service is active
func (servManager *ServManagerV1) Ping(ign *utils.CGREventWithArgDispatcher, reply *string) error {
	*reply = utils.Pong
	return nil
}

// Call implements rpcclient.RpcClientConnection interface for internal RPC
func (servManager *ServManagerV1) Call(serviceMethod string,
	args interface{}, reply interface{}) error {
	return utils.APIerRPCCall(servManager, serviceMethod, args, reply)
}
