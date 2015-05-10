package oneandone_cloudserver_api

import (
	log "github.com/Sirupsen/logrus"
)

type ServerAppliance struct {
	withId
	withName
	OsImageType        string   `json:"os_image_type"`
	OsFamily           string   `json:"os_family"`
	Os                 string   `json:"os"`
	OsVersion          string   `json:"os_version"`
	MinHddSize         int      `json:"min_hdd_size"`
	Architecture       int      `json:"architecture"`
	Licenses           []string `json:"licenses"`
	IsAutomaticInstall bool     `json:"automatic_installation"`
	Type               string   `json:"type"`
	withApi
}

// GET /server_appliances
func (api *API) GetServerAppliances() []ServerAppliance {
	log.Debug("requesting information about server appliances")
	session := api.prepareSession()
	res := []ServerAppliance{}
	resp, _ := session.Get(createUrl(api, "server_appliances"), nil, &res, nil)
	logResult(resp, 200)
	for index, _ := range res {
		res[index].api = api
	}
	return res
}

// GET /server_appliances/{id}
func (api *API) GetServerAppliance(Id string) ServerAppliance {
	log.Debug("requesting information about server appliance", Id)
	session := api.prepareSession()
	res := ServerAppliance{}
	resp, _ := session.Get(createUrl(api, "server_appliances"), nil, &res, nil)
	logResult(resp, 200)
	res.api = api
	return res
}
