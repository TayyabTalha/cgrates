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

package agents

import (
	"fmt"
	"strings"

	"github.com/cgrates/cgrates/config"
	"github.com/cgrates/cgrates/engine"
	"github.com/cgrates/cgrates/utils"
	"github.com/cgrates/rpcclient"
	"github.com/miekg/dns"
)

// NewDNSAgent is the constructor for DNSAgent
func NewDNSAgent(cgrCfg *config.CGRConfig, fltrS *engine.FilterS,
	sS rpcclient.RpcClientConnection) (da *DNSAgent, err error) {
	da = &DNSAgent{cgrCfg: cgrCfg, fltrS: fltrS, sS: sS}
	return
}

// DNSAgent translates DNS requests towards CGRateS infrastructure
type DNSAgent struct {
	cgrCfg *config.CGRConfig             // loaded CGRateS configuration
	fltrS  *engine.FilterS               // connection towards FilterS
	sS     rpcclient.RpcClientConnection // connection towards CGR-SessionS component
}

// ListenAndServe will run the DNS handler doing also the connection to listen address
func (da *DNSAgent) ListenAndServe() error {
	utils.Logger.Info(fmt.Sprintf("<%s> start listening on <%s:%s>",
		utils.DNSAgent, da.cgrCfg.DNSAgentCfg().ListenNet, da.cgrCfg.DNSAgentCfg().Listen))
	if strings.HasSuffix(da.cgrCfg.DNSAgentCfg().ListenNet, utils.TLSNoCaps) {
		return dns.ListenAndServeTLS(
			da.cgrCfg.DNSAgentCfg().Listen,
			da.cgrCfg.TlsCfg().ServerCerificate,
			da.cgrCfg.TlsCfg().ServerKey,
			dns.HandlerFunc(
				func(w dns.ResponseWriter, m *dns.Msg) {
					go da.handleMessage(w, m)
				}),
		)
	}
	return dns.ListenAndServe(
		da.cgrCfg.DNSAgentCfg().Listen,
		da.cgrCfg.DNSAgentCfg().ListenNet,
		dns.HandlerFunc(
			func(w dns.ResponseWriter, m *dns.Msg) {
				go da.handleMessage(w, m)
			}),
	)
}

// handleMessage is the entry point of all DNS requests
// requests are reaching here asynchronously
func (da *DNSAgent) handleMessage(w dns.ResponseWriter, req *dns.Msg) {
	fmt.Printf("got message: %+v\n", req)
	/*rply := new(dns.Msg)
	rply.SetReply(req)
	switch req.Question[0].Qtype {
	case dns.TypeA:
		rply.Authoritative = true
		if req.Question[0].Name == "cgrates.org." {
			rply.Answer = append(rply.Answer,
				&dns.A{
					Hdr: dns.RR_Header{
						Name:   req.Question[0].Name,
						Rrtype: dns.TypeA,
						Class:  dns.ClassINET,
						Ttl:    60},
					A: net.ParseIP("195.201.167.179")},
			)
		}
	}
	w.WriteMsg(rply)
	*/
	reqVars := make(map[string]interface{})
	reqVars[QueryType] = dns.TypeToString[req.Question[0].Qtype]
	if req.Question[0].Qtype == dns.TypeNAPTR {
		e164, err := e164FromNAPTR(req.Question[0].Name)
		if err != nil {
			utils.Logger.Warning(
				fmt.Sprintf("<%s> decoding NAPTR query: <%s>, err: %s",
					utils.DNSAgent, req.Question[0].Name, err.Error()))
		}
		reqVars[E164Address] = e164
	}
}
