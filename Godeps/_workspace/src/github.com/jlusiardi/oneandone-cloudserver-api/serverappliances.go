/*
 * Copyright 2015 1&1 Internet AG, http://1und1.de . All rights reserved. Licensed under the Apache v2 License.
 */

package oneandone_cloudserver_api

import (
	"errors"
	"github.com/docker/machine/log"
	"net/http"
	"sort"
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

//Functions for the sort.Sort interface to sort the serverAppliance struct by OsVersion
type sortServerAppliance []ServerAppliance

func (s sortServerAppliance) Less(i, j int) bool {
	return s[i].OsVersion > s[j].OsVersion
}

func (s sortServerAppliance) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s sortServerAppliance) Len() int {
	return len(s)
}

// Function to get the available architectures for the given operating system
//
// Returns the available architectures. i.E. [32, 64]
func (api *API) ServerApplianceListArchitectures(family string, os string, osType string) ([]int, error) {
	apps, err := api.GetServerAppliances()
	if err != nil {
		return nil, err
	}
	architectures := make(map[int]int)
	for index, _ := range apps {
		if apps[index].OsFamily == family && apps[index].Os == os && apps[index].OsImageType == osType {
			log.Debug(apps[index])
			architectures[apps[index].Architecture] = 1
		}
	}
	if len(architectures) >= 1 {
		return getMapKeysInt(architectures), nil
	}
	return nil, errors.New("No entries found with given parameters")
}

// Function to get the available operating system type images
//
//Returns the available operating system types. i.E. [Minimal, Standard, ISO_OS]
func (api *API) ServerApplianceListTypes(family string, os string) ([]string, error) {
	apps, err := api.GetServerAppliances()
	if err != nil {
		return nil, err
	}
	osTypes := make(map[string]int)
	for index, _ := range apps {
		if apps[index].OsFamily == family && apps[index].Os == os {
			log.Debug(apps[index])
			osTypes[apps[index].OsImageType] = 1
		}
	}
	if len(osTypes) > 1 {
		return getMapKeysString(osTypes), nil
	}
	return nil, errors.New("No entries found with given parameters")
}

// Function to get the available operating system by the os family
//
// Returns all operating systems who are in the given family. i.E. Linux: [Ubuntu, Debian] and so on..
func (api *API) ServerApplianceListOperationSystems(family string) ([]string, error) {
	apps, err := api.GetServerAppliances()
	if err != nil {
		return nil, err
	}
	os := make(map[string]int)
	for index, _ := range apps {
		if apps[index].OsFamily == family {
			log.Debug(apps[index])
			os[apps[index].Os] = 1
		}
	}
	if len(os) >= 1 {
		return getMapKeysString(os), nil
	}
	return nil, errors.New("No entries found with given parameters")
}

// Function to get the available operating system families
//
// Returns the available operating system families. i.E. [Linux, Windows]
func (api *API) ServerApplianceListFamilies() ([]string, error) {
	apps, err := api.GetServerAppliances()
	if err != nil {
		return nil, err
	}
	osFamilies := make(map[string]int)
	for index, _ := range apps {
		log.Debug(apps[index])
		osFamilies[apps[index].OsFamily] = 1
	}
	if len(osFamilies) >= 1 {
		return getMapKeysString(osFamilies), nil
	}
	return nil, errors.New("No entries found")
}

// Function to get the newest operating system
//
// Returns the newest operating system as ServerAppliance object
func (api *API) ServerApplianceFindNewest(family string, os string, osType string, architecture int, autoInstall bool) (*ServerAppliance, error) {
	apps, err := api.GetServerAppliances()
	if err != nil {
		return nil, err
	}
	filteredApps := sortServerAppliance{}
	for index, _ := range apps {
		if apps[index].OsFamily == family && apps[index].Os == os && apps[index].OsImageType == osType &&
			apps[index].Architecture == architecture && apps[index].IsAutomaticInstall == autoInstall {
			log.Debugf("Appliance with id '%v' matches criteria", apps[index].Id)
			filteredApps = append(filteredApps, apps[index])
		}
	}
	sort.Sort(filteredApps)
	if len(filteredApps) >= 1 {
		return &filteredApps[0], nil
	}
	return nil, errors.New("No entries found with given parameters")
}
