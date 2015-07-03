/*
 * Copyright 2015 1&1 Internet AG, http://1und1.de . All rights reserved. Licensed under the Apache v2 License.
 */

package oneandone_cloudserver_api

import (
	"github.com/docker/machine/log"
	"net/http"
	"time"
)

type FirewallPolicy struct {
	withId
	withName
	withDescription
	// should be fixed: Status        Status                   `json:"status"`
	Status        string                   `json:"state"`
	DefaultPolicy int                      `json:"default"`
	Rules         []FirewallPolicyRules    `json:"rules"`
	ServerIps     []FirewallPolicyServerIp `json:"server_ips"`
	withApi
}

type FirewallPolicyRules struct {
	withId
	Protocol string `json:"protocol"`
	PortFrom *int   `json:"port_from"`
	PortTo   *int   `json:"port_to"`
	SourceIp string `json:"source"`
}

type FirewallPolicyServerIp struct {
	withId
	Ip         string `json:"ip"`
	ServerName string `json:"server_name"`
}

type FirewallPolicyCreateData struct {
	Name        string                          `json:"name"`
	Description string                          `json:"description"`
	Rules       []FirewallPolicyRulesCreateData `json:"rules"`
}

type FirewallPolicyRulesCreateData struct {
	Protocol string `json:"protocol"`
	PortFrom *int   `json:"port_from"`
	PortTo   *int   `json:"port_to"`
	SourceIp string `json:"source"`
}

type FirewallPolicyAddIpsData struct {
	ServerIps []string `json:"server_ips"`
}

// GET /firewall_policies
func (api *API) GetFirewallPolicies() ([]FirewallPolicy, error) {
	log.Debug("requesting information about firewall policies")
	result := []FirewallPolicy{}
	err := api.Client.Get(createUrl(api, "firewall_policies"), &result, http.StatusOK)
	if err != nil {
		return nil, err
	}
	for index, _ := range result {
		result[index].api = api
	}
	return result, nil
}

// POST /firewall_policies
func (api *API) CreateFirewallPolicy(configuration FirewallPolicyCreateData) (*FirewallPolicy, error) {
	log.Debug("requesting to create a new firewall policy")
	result := new(FirewallPolicy)
	err := api.Client.Post(createUrl(api, "firewall_policies"), configuration, &result, http.StatusAccepted)
	if err != nil {
		return nil, err
	}
	result.api = api
	return result, nil
}

// GET /firewall_policies/{id}
func (api *API) GetFirewallPolicy(Id string) (*FirewallPolicy, error) {
	log.Debugf("requesting information about firewall policy: '%s'", Id)
	result := new(FirewallPolicy)
	err := api.Client.Get(createUrl(api, "firewall_policies", Id), &result, http.StatusOK)
	if err != nil {
		return nil, err
	}
	result.api = api
	return result, nil

}

// DELETE /firewall_policies/{id}
func (fwp *FirewallPolicy) Delete() (*FirewallPolicy, error) {
	log.Debug("Requested to delete firewall policy ", fwp.Id)
	result := new(FirewallPolicy)
	err := fwp.api.Client.Delete(createUrl(fwp.api, "firewall_policies", fwp.Id), &result, http.StatusAccepted)
	if err != nil {
		return nil, err
	}
	result.api = fwp.api
	return result, nil
}

// PUT /firewall_policies/{id}

// GET /firewall_policies/{id}/server_ips

// PUT /firewall_policies/{id}/server_ips
func (fwp *FirewallPolicy) AddServerIp(ipId string) (*FirewallPolicy, error) {
	log.Debugf("Requested to apply firewall policy '%v' to ip '%v'", fwp.Id, ipId)
	result := new(FirewallPolicy)
	request := FirewallPolicyAddIpsData{
		ServerIps: []string{ipId},
	}
	err := fwp.api.Client.Put(createUrl(fwp.api, "firewall_policies", fwp.Id, "server_ips"), request, result, http.StatusAccepted)
	if err != nil {
		return nil, err
	}
	result.api = fwp.api
	return result, nil
}

// GET /firewall_policies/{id}/server_ips/{id}

// DELETE /firewall_policies/{id}/server_ips/{id}
func (fwp *FirewallPolicy) DeleteServerIp(ipId string) (*FirewallPolicy, error) {
	log.Debugf("Requested to remove firewall policy '%v' from ip '%v'", fwp.Id, ipId)
	result := new(FirewallPolicy)
	err := fwp.api.Client.Delete(createUrl(fwp.api, "firewall_policies", fwp.Id, "server_ips", ipId), result, http.StatusAccepted)
	if err != nil {
		return nil, err
	}
	result.api = fwp.api
	return result, nil
}

// GET /firewall_policies/{id}/rules

// PUT /firewall_policies/{id}/rules

// GET /firewall_policies/{id}/rules/{id}

// DELETE /firewall_policies/{id}/rules/{id}

func (fwp *FirewallPolicy) exists() (bool, error) {
	_, err := fwp.api.GetFirewallPolicy(fwp.Id)
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

func (fwp *FirewallPolicy) WaitUntilDeleted() error {
	exists := true
	var err error
	for exists {
		exists, err = fwp.exists()
		if err != nil {
			return err
		}
		log.Debugf("Wait for firewall policy: '%s' to be deleted", fwp.Id)
		time.Sleep(5 * time.Second)
	}
	log.Infof("The firewall policy: '%s' is now deleted", fwp.Id)
	return nil
}

func (fwp *FirewallPolicy) WaitForState(State string) error {
	fw, err := fwp.api.GetFirewallPolicy(fwp.Id)
	if err != nil {
		return err
	}
	for fw.Status != State {
		time.Sleep(5 * time.Second)
		fw, err := fwp.api.GetFirewallPolicy(fwp.Id)
		if err != nil {
			return err
		}
		if fw.Status == State {
			log.Infof("The firewall policy is now in the expected state: '%s'", State)
			return nil
		} else {
			log.Debugf("Wait for expected status: '%s' current: '%s'", State, fw.Status)
		}
	}
	return nil
}
