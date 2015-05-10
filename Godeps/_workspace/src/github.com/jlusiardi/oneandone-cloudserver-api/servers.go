package oneandone_cloudserver_api

import (
	log "github.com/Sirupsen/logrus"
)

type Server struct {
	withId
	withName
	withDescription
	Password string        `json:"first_password"`
	Status   Status        `json:"status"`
	Hardware Hardware      `json:"hardware"`
	Image    ImageInServer `json:"image"`
	withApi
}

type Hardware struct {
	Vcores            int   `json:"vcore"`
	CoresPerProcessor int   `json:"cores_per_processor"`
	Ram               int   `json:"ram"`
	Hdds              []Hdd `json:"hdds"`
}

type Hdd struct {
	withId
	Size   int  `json:"size"`
	IsMain bool `json:"is_main"`
}

type ImageInServer struct {
	withId
	withName
}

type IpInServer struct {
	withId
	withType
	Ip         string `json:"ip"`
	ReverseDns string `json:"reverse_dns"`
	FirewallId string `json:"firewall"`
}

type ServerCreateData struct {
	withName
	withDescription
	Hardware           Hardware `json:"appliance_id"`
	ApplianceId        string   `json:"password"`
	Password           string   `json:"password"`
	PowerOn            bool     `json:"power_on"`
	FirewallPolicyId   string   `json:"firewall_policy_id"`
	IpId               string   `json:"ip_id"`
	LoadBalancerId     string   `json:"load_balancer_id"`
	MonitoringPolicyId string   `json:"monitoring_policy_id"`
	PrivateNetworkId   string   `json:"private_network_id"`
}

// GET /servers
func (api *API) GetServers() []Server {
	log.Debug("requesting information about servers")
	session := api.prepareSession()
	res := []Server{}
	resp, _ := session.Get(createUrl(api, "servers"), nil, &res, nil)
	logResult(resp, 200)
	for index, _ := range res {
		res[index].api = api
	}
	return res
}

// POST /servers
func (api *API) CreateServer(configuration ServerCreateData) Server {
	log.Debug("requesting to create a new server")
	s := api.prepareSession()
	res := Server{}
	resp, _ := s.Post(createUrl(api, "servers"), configuration, &res, nil)
	logResult(resp, 200)
	res.api = api
	return res
}

// GET /servers/{id}
func (api *API) GetServer(Id string) Server {
	log.Debug("requesting to start server ", Id)
	session := api.prepareSession()
	res := Server{}
	resp, _ := session.Get(createUrl(api, "servers", Id), nil, &res, nil)
	logResult(resp, 200)
	res.api = api
	return res
}

// DELETE /servers/{id}
func (server *Server) Delete() Server {
	log.Debug("Requested to delete VM ", server.Id)
	session := server.api.prepareSession()
	res := Server{}
	resp, _ := session.Delete(createUrl(server.api, "servers", server.Id), &res, nil)
	logResult(resp, 200)
	res.api = server.api
	return res
}

// PUT /servers/{id}

// GET /servers/{id}/size

// PUT /servers/{id}/size

// GET /servers/{id}/size/hdds

// POST /servers/{id}/size/hdds

// GET /servers/{id}/size/hdds/{id}

// DELETE /servers/{id}/size/hdds/{id}

// PUT /servers/{id}/size/hdds/{id}

// GET /servers/{id}/image

// PUT /servers/{id}/image

// GET /servers/{id}/ips

// POST /servers/{id}/ips

// GET /servers/{id}/ips/{id}

// DELETE /servers/{id}/ips/{id}

// PUT /servers/{id}/ips/{id}

// GET /servers/{id}/status

// PUT /servers/{id}/status

// GET /servers/{id}/dvd

// DELETE /servers/{id}/dvd

// PUT /servers/{id}/dvd

// GET /servers/{id}/private_networks

// PUT /servers/{id}/private_networks

// GET /servers/{id}/private_networks/{id}

// DELETE /servers/{id}/private_networks/{id}

// GET /servers/{id}/snapshots

// POST /servers/{id}/snapshots

// GET /servers/{id}/snapshots/{id}

// DELETE /servers/{id}/snapshots/{id}

// PUT /servers/{id}/snapshots/{id}

// POST /servers/{server_id}/clone
func (server *Server) Clone(NewName string) Server {
	return Server{}
}
