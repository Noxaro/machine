package oneandone_cloudserver_api

import (
	log "github.com/Sirupsen/logrus"
)

type DvdIso struct {
	withId
	withName
	OsFamily     string `json:"os_family"`
	Os           string `json:"os"`
	OsVersion    string `json:"os_version"`
	Architecture int    `json:"architecture"`
	Type         string `json:"type"`
	withApi
}

// GET /dvd_isos
func (api *API) GetDvdIsos() []DvdIso {
	log.Debug("requesting information about dvd isos")
	session := api.prepareSession()
	res := []DvdIso{}
	resp, _ := session.Get(createUrl(api, "dvd_isos"), nil, &res, nil)
	logResult(resp, 200)
	for index, _ := range res {
		res[index].api = api
	}
	return res
}

// GET /dvd_isos/{id}
func (api *API) GetDvdIso(Id string) DvdIso {
	log.Debug("requesting information about dvd iso", Id)
	session := api.prepareSession()
	res := DvdIso{}
	resp, _ := session.Get(createUrl(api, "dvd_isos", Id), nil, &res, nil)
	logResult(resp, 200)
	res.api = api
	return res
}
