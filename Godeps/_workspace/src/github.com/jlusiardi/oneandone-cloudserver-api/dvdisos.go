/*
 * Copyright 2015 1&1 Internet AG, http://1und1.de . All rights reserved. Licensed under the Apache v2 License.
 */

package oneandone_cloudserver_api

import (
	"github.com/docker/machine/log"
	"net/http"
)

// Struct to describe a ISO image that can be used to boot a server.
//
// Values of this type describe ISO images that can be inserted into the servers virtual DVD drive.
//
//
type DvdIso struct {
	Id           string `json:"id"`
	Name         string `json:"name"`
	OsFamily     string `json:"os_family"`
	Os           string `json:"os"`
	OsVersion    string `json:"os_version"`
	Architecture int    `json:"architecture"`
	Type         string `json:"type"`
	withApi
}

// GET /dvd_isos
func (api *API) GetDvdIsos() ([]DvdIso, error) {
	log.Debug("requesting information about dvd isos")
	result := []DvdIso{}
	err := api.Client.Get(createUrl(api, "dvd_isos"), &result, http.StatusOK)
	if err != nil {
		return nil, err
	}
	for index, _ := range result {
		result[index].api = api
	}
	return result, nil
}

// GET /dvd_isos/{id}
func (api *API) GetDvdIso(Id string) (*DvdIso, error) {
	log.Debug("requesting information about dvd iso", Id)
	result := new(DvdIso)
	err := api.Client.Get(createUrl(api, "dvd_isos", Id), &result, http.StatusOK)
	if err != nil {
		return nil, err
	}
	result.api = api
	return result, nil
}
