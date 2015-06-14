package oneandone_cloudserver_api

import (
	"github.com/docker/machine/log"
	"net/http"
)

type ServerAppliance struct {
	withId
	withName
	OsImageType        string                   `json:"os_image_type"`
	OsFamily           string                   `json:"os_family"`
	Os                 string                   `json:"os"`
	OsVersion          string                   `json:"os_version"`
	MinHddSize         int                      `json:"min_hdd_size"`
	Architecture       int                      `json:"architecture"`
	Licenses           []ServerApplianceLicence `json:"licenses"`
	IsAutomaticInstall bool                     `json:"automatic_installation"`
	Type               string                   `json:"type"`
	withApi
}

type ServerApplianceLicence struct {
	withName
}

// GET /server_appliances
func (api *API) GetServerAppliances() ([]ServerAppliance, error) {
	log.Debug("requesting information about server appliances")
	res := []ServerAppliance{}
	err := api.Client.Get(createUrl(api, "server_appliances"), &res, http.StatusOK)
	if err != nil {
		return nil, err
	}
	for index, _ := range res {
		res[index].api = api
	}
	return res, nil
}

// GET /server_appliances/{id}
func (api *API) GetServerAppliance(Id string) (*ServerAppliance, error) {
	log.Debug("requesting information about server appliance", Id)
	res := new(ServerAppliance)
	err := api.Client.Get(createUrl(api, "server_appliances", Id), &res, http.StatusOK)
	if err != nil {
		return nil, err
	}
	res.api = api
	return res, nil
}
