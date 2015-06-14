package oneandone_cloudserver_api

import (
	"github.com/docker/machine/log"
	"net/http"
)

type Server struct {
	withId
	withName
	withDescription
	Password string        `json:"first_password"`
	Status   Status        `json:"status"`
	Hardware Hardware      `json:"hardware"`
	Image    ImageInServer `json:"image"`
	Ips      []IpInServer  `json:"ips"`
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
	Name               string   `json:"name"`
	Description        string   `json:"description"`
	Hardware           Hardware `json:"hardware"`
	ApplianceId        string   `json:"appliance_id"`
	Password           string   `json:"password"`
	PowerOn            bool     `json:"power_on"`
	FirewallPolicyId   string   `json:"firewall_policy_id"`
	IpId               string   `json:"ip_id"`
	LoadBalancerId     string   `json:"load_balancer_id"`
	MonitoringPolicyId string   `json:"monitoring_policy_id"`
	PrivateNetworkId   string   `json:"private_network_id"`
}

// GET /servers
func (api *API) GetServers() ([]Server, error) {
	log.Debug("requesting information about servers")
	result := []Server{}
	err := api.Client.Get(createUrl(api, "servers"), &result, http.StatusOK)
	if err != nil {
		return nil, err
	}
	for index, _ := range result {
		result[index].api = api
	}
	return result, nil
}

// POST /servers
func (api *API) CreateServer(configuration ServerCreateData) (*Server, error) {
	log.Debug("requesting to create a new server")
	result := new(Server)
	err := api.Client.Post(createUrl(api, "servers"), &configuration, &result, http.StatusAccepted)
	if err != nil {
		return nil, err
	}
	result.api = api
	return result, nil
}

// GET /servers/{id}
func (api *API) GetServer(Id string) (*Server, error) {
	log.Debug("requesting to about server ", Id)
	result := new(Server)
	err := api.Client.Get(createUrl(api, "servers", Id), &result, http.StatusOK)
	if err != nil {
		return nil, err
	}
	result.api = api
	return result, nil
}

// DELETE /servers/{id}
func (s *Server) Delete() (*Server, error) {
	log.Debug("Requested to delete VM ", s.Id)
	result := new(Server)
	err := s.api.Client.Delete(createUrl(s.api, "servers", s.Id), &result, http.StatusOK)
	if err != nil {
		return nil, err
	}
	result.api = s.api
	return result, nil
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
