/*
 * Copyright 2015 1&1 Internet AG, http://1und1.de . All rights reserved. Licensed under the Apache v2 License.
 */

package oneandone_cloudserver_api

import (
	"github.com/docker/machine/log"
	"net/http"
	"time"
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
	Ip         string               `json:"ip"`
	ReverseDns string               `json:"reverse_dns"`
	Firewall   FirewallInIpInServer `json:"firewall_policy"`
}

type FirewallInIpInServer struct {
	withId
	withName
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

type ServerAction struct {
	Action string `json:"action"`
	Method string `json:"method"`
}

type FixedInstanceInformation struct {
	withName
	withId
	Hardware Hardware `json:"hardware"`
	withApi
}

type ServerRenameData struct {
	Name        string `json:"name"`
	Description string `json:"description"`
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
	log.Debug("requesting information about server ", Id)
	result := new(Server)
	err := api.Client.Get(createUrl(api, "servers", Id), &result, http.StatusOK)
	if err != nil {
		return nil, err
	}
	result.api = api
	return result, nil
}

func (api *API) GetFixedInstanceSizes() ([]FixedInstanceInformation, error) {
	log.Debug("requesting information about fixed instance sizes")
	result := []FixedInstanceInformation{}
	err := api.Client.Get(createUrl(api, "servers", "fixed_instance_sizes"), &result, http.StatusOK)
	if err != nil {
		return nil, err
	}
	for index, _ := range result {
		result[index].api = api
	}
	return result, nil
}

// DELETE /servers/{id}
func (s *Server) Delete() (*Server, error) {
	log.Debugf("Requested to delete server '%v' ", s.Id)
	result := new(Server)
	err := s.api.Client.Delete(createUrl(s.api, "servers", s.Id), &result, http.StatusAccepted)
	if err != nil {
		return nil, err
	}
	result.api = s.api
	return result, nil
}

// PUT /servers/{id}
func (s *Server) RenameServer(data ServerRenameData) (*Server, error) {
	log.Debugf("Requested to rename server '%v'", s.Id)
	result := new(Server)
	err := s.api.Client.Put(createUrl(s.api, "servers", s.Id), data, &result, http.StatusOK)
	if err != nil {
		return nil, err
	}
	result.api = s.api
	return result, nil
}

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
func (s *Server) GetStatus() (*Status, error) {
	log.Debugf("Requesting server status for server: '%s'", s.Id)
	result := new(Status)
	err := s.api.Client.Get(createUrl(s.api, "servers", s.Id, "status"), &result, http.StatusOK)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// PUT /servers/{id}/status/action
func (s *Server) Reboot(hardware bool) (*Server, error) {
	log.Debugf("Requested to reboot Server '%v'. Hardware: '%t'", s.Id, hardware)
	result := new(Server)
	request := ServerAction{}
	request.Action = "REBOOT"
	if hardware {
		request.Method = "HARDWARE"
	} else {
		request.Method = "SOFTWARE"
	}
	err := s.api.Client.Put(createUrl(s.api, "servers", s.Id, "status", "action"), request, &result, http.StatusAccepted)
	if err != nil {
		return nil, err
	}
	result.api = s.api
	return result, nil
}

func (s *Server) Shutdown(hardware bool) (*Server, error) {
	log.Debugf("Requested to shutdown Server '%v'. Hardware: '%t'", s.Id, hardware)
	result := new(Server)
	request := ServerAction{}
	request.Action = "POWER_OFF"
	if hardware {
		request.Method = "HARDWARE"
	} else {
		request.Method = "SOFTWARE"
	}
	err := s.api.Client.Put(createUrl(s.api, "servers", s.Id, "status", "action"), request, &result, http.StatusAccepted)
	if err != nil {
		return nil, err
	}
	result.api = s.api
	return result, nil
}

func (s *Server) Start() (*Server, error) {
	log.Debugf("Requested to start Server '%v'.", s.Id)
	result := new(Server)
	request := ServerAction{}
	request.Action = "POWER_ON"
	err := s.api.Client.Put(createUrl(s.api, "servers", s.Id, "status", "action"), request, &result, http.StatusAccepted)
	if err != nil {
		return nil, err
	}
	result.api = s.api
	return result, nil
}

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

func (server *Server) exists() (bool, error) {
	_, err := server.api.GetServer(server.Id)
	if err == nil {
		return true, nil
	} else {
		if apiError, ok := err.(ApiError); ok && apiError.httpStatusCode == http.StatusNotFound {
			return false, nil
		} else {
			return false, err
		}
	}
}

func (server *Server) WaitUntilDeleted() error {
	exists := true
	var err error
	for exists {
		exists, err = server.exists()
		if err != nil {
			return err
		}
		log.Debugf("Wait for server: '%s' to be deleted", server.Id)
		time.Sleep(5 * time.Second)
	}
	log.Infof("The server: '%s' is now deleted", server.Id)
	return nil
}

// Function to perform busy-wating for a certain server state.
//
// This function queries the server with the given id every 5s until the server's state equals the given state.
func (server *Server) WaitForState(State string) error {
	server, err := server.api.GetServer(server.Id)
	if err != nil {
		return err
	}
	status := server.Status
	log.Infof("Wait for expected status: '%s' current: '%s' %d%%", State, status.State, status.Percent)
	for status.State != State {
		time.Sleep(5 * time.Second)
		status, err := server.GetStatus()
		if err != nil {
			return err
		}
		if status.State == State {
			log.Infof("The server is now in the expected state: '%s'", State)
			return nil
		} else {
			log.Debugf("Wait for expected status: '%s' current: '%s' %d%%", State, status.State, status.Percent)
		}
	}
	return nil
}
