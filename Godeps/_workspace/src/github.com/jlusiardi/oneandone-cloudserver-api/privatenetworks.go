package oneandone_cloudserver_api

import (
	"github.com/docker/machine/log"
	"net/http"
)

type PrivateNetwork struct {
	withId
	withName
	CloudpanelId string `json:"cloudpanel_id"`
	withDescription
	NetworkAddress      string                 `json:"network_address"`
	SubnetMask          string                 `json:"subnet_mask"`
	State               string                 `json:"state"`
	CreationDate        string                 `json:"creation_date"`
	Servers             []PrivateNetworkServer `json:"servers"`
	withApi
}

type PrivateNetworkServer struct {
	withId
	withName
	pn *PrivateNetwork `json:"omitempty"`
}

type PrivateNetworkSettings struct {
	Name                string              `json:"name"`
	Description         string              `json:"description"`
	NetworkAddress      string              `json:"network_address"`
	SubnetMask          string              `json:"subnet_mask"`
}

type PrivateNetworkAddServers struct {
	Servers []string `json:"servers"`
}

// GET /private_networks
func (api *API) GetPrivateNetworks() ([]PrivateNetwork, error) {
	log.Debug("Requesting informations about private networks")
	result := []PrivateNetwork{}
	err := api.Client.Get(createUrl(api, PrivateNetworksPathSegment), &result, http.StatusOK)
	if err != nil {
		return nil, err
	}
	for index, _ := range result {
		result[index].api = api
	}
	return result, nil
}

// POST /private_networks
func (api *API) CreatePrivateNetwork(configuration PrivateNetworkSettings) (*PrivateNetwork, error) {
	log.Debugf("Creating a new private network: '%s'", configuration.Name)
	result := new(PrivateNetwork)
	err := api.Client.Post(createUrl(api, PrivateNetworksPathSegment), &configuration, &result, http.StatusAccepted)
	if err != nil {
		return nil, err
	}
	result.api = api
	return result, nil
}

// GET /private_networks/{id}
func (api *API) GetPrivateNetwork(Id string) (*PrivateNetwork, error) {
	log.Debugf("Requesting informations about the private network: '%s'", Id)
	result := new(PrivateNetwork)
	err := api.Client.Get(createUrl(api, PrivateNetworksPathSegment, Id), &result, http.StatusOK)
	if err != nil {
		return nil, err
	}
	result.api = api
	return result, nil
}

// DELETE /private_networks/{id}
func (pn *PrivateNetwork) Delete() (*PrivateNetwork, error) {
	log.Debugf("Deleting private network: '%s'", pn.Id)
	result := new(PrivateNetwork)
	err := pn.api.Client.Delete(createUrl(pn.api, PrivateNetworksPathSegment, pn.Id), &result, http.StatusAccepted)
	if err != nil {
		return nil, err
	}
	result.api = pn.api
	return result, nil
}

// PUT /private_networks/{id}
func (pn *PrivateNetwork) Update(configuration PrivateNetworkSettings) (*PrivateNetwork, error) {
	log.Debugf("Updateing private network: '%s'", pn.Id)
	result := new(PrivateNetwork)
	err := pn.api.Client.Put(createUrl(pn.api, PrivateNetworksPathSegment, pn.Id), &configuration, &result, http.StatusOK)
	if err != nil {
		return nil, err
	}
	result.api = pn.api
	return result, nil
}

// GET /private_networks/{id}/servers
func (pn *PrivateNetwork) GetServers() ([]PrivateNetworkServer, error) {
	log.Debugf("Requesting servers witch are attached on the private network: '%s'", pn.Id)
	result := []PrivateNetworkServer{}
	err := pn.api.Client.Get(createUrl(pn.api, PrivateNetworksPathSegment, pn.Id, "servers"), &result, http.StatusOK)
	if err != nil {
		return nil, err
	}
	for index, _ := range result {
		result[index].pn = pn
	}
	return result, nil
}

// PUT /private_networks/{id}/servers
func (pn *PrivateNetwork) AddServers(privateNetworkAddServers PrivateNetworkAddServers) (*PrivateNetwork, error) {
	log.Debugf("Adding servers to the private network: '%s'", pn.Id)
	result := new(PrivateNetwork)
	err := pn.api.Client.Put(createUrl(pn.api, PrivateNetworksPathSegment, pn.Id, "servers"), &privateNetworkAddServers, &result, http.StatusAccepted)
	if err != nil {
		return nil, err
	}
	result.api = pn.api
	return result, nil
}

// GET /private_networks/{id}/servers/{id}
func (pn *PrivateNetwork) GetServer(id string) (*PrivateNetworkServer, error) {
	log.Debugf("Requesting server: '%s' from the private network: '%s'", id, pn.Id)
	result := new(PrivateNetworkServer)
	err := pn.api.Client.Get(createUrl(pn.api, PrivateNetworksPathSegment, pn.Id, "servers", id), &result, http.StatusOK)
	if err != nil {
		return nil, err
	}
	result.pn = pn
	return result, nil
}

// DELETE /private_networks/{id}/servers/{id}
func (pns *PrivateNetworkServer) Delete() (*PrivateNetwork, error) {
	log.Debugf("Deleting server: '%s' from the private network: '%s'", pns.Id, pns.pn.Id)
	result := new(PrivateNetwork)
	err := pns.pn.api.Client.Delete(createUrl(pns.pn.api, PrivateNetworksPathSegment, pns.pn.Id, "servers", pns.Id), &result, http.StatusAccepted)
	if err != nil {
		return nil, err
	}
	result.api = pns.pn.api
	return result, nil
}
