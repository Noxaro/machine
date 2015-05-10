package oneandone_cloudserver_api

import ()

type PrivateNetwork struct {
	withId
	withName
	CloudpanelId string `json:"cloudpanel_id"`
	withDescription
	NetworkAddress string   `json:"network_address"`
	SubnetMask     string   `json:"subnet_mask"`
	State          string   `json:"state"`
	CreationDate   string   `json:"creation_date"`
	Servers        []ReferencedServer `json:"servers"`
}

type ReferencedServer struct {
	withId
	withName
}

type PrivateNetworkCreateData struct {
}

// GET /private_networks
func (api *API) GetPrivateNetworks() []PrivateNetwork {
	return []PrivateNetwork{}
}

// POST /private_networks
func (api *API) CreatePrivateNetwork(configuration PrivateNetworkCreateData) PrivateNetwork {
	return PrivateNetwork{}
}

// GET /private_networks/{id}
func (api *API) GetPrivateNetwork(Id string) PrivateNetwork {
	return PrivateNetwork{}
}

// DELETE /private_networks/{id}
func (fwp *PrivateNetwork) Delete() PrivateNetwork {
	return PrivateNetwork{}
}

// PUT /private_networks/{id}

// GET /private_networks/{id}/servers

// PUT /private_networks/{id}/servers

// GET /private_networks/{id}/servers/{id}

// DELETE /private_networks/{id}/servers/{id}

