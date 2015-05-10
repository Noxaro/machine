package oneandone_cloudserver_api

import (
	log "github.com/Sirupsen/logrus"
)

type PublicIp struct {
	withId
	IpAddress    string      `json:"ip"`
	withType
	AssignedTo   ServerForIp `json:"assigned_to"`
	ReverseDns   string      `json:"reverse_dns"`
	IsDhcp       bool        `json:"dhcp"`
	State        string      `json:"state"`
	CreationDate string      `json:"creation_date"`
	withApi
}

type ServerForIp struct {
	withId
	withName
	withType
}

type PublicIpCreateData struct {
	ReverseDns string `json:"reverse_dns"`
}

// GET /public_ips
func (api *API) GetPublicIps() []PublicIp {
	log.Debug("requesting information about public ips")
	session := api.prepareSession()
	res := []PublicIp{}
	resp, _ := session.Get(createUrl(api, "public_ips"), nil, &res, nil)
	logResult(resp, 200)
	for index, _ := range res {
		res[index].api = api
	}
	return res
}

// POST /public_ips
func (api *API) CreatePublicIp(configuration PublicIpCreateData) PublicIp {
	return PublicIp{}
}

// GET /public_ips/{id}
func (api *API) GetPublicIp(Id string) PublicIp {
	return PublicIp{}
}

// DELETE /public_ips/{id}
func (ip *PublicIp) Delete() PublicIp {
	return PublicIp{}
}

// PUT /public_ips/{id}
func (ip *PublicIp) UpdateReverseDns(NewReverseDns string) PublicIp {
	return PublicIp{}
}
