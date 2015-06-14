package oneandone_cloudserver_api

import (
	"github.com/docker/machine/log"
	"net/http"
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
	err := api.Client.Post(createUrl(api, "firewall_policies"), configuration, &result, http.StatusCreated)
	if err != nil {
		return nil, err
	}
	result.api = api
	return result, nil
}

// GET /firewall_policies/{id}
func (api *API) GetFirewallPolicy(Id string) (*FirewallPolicy, error) {
	log.Debug("requesting to about firewall policy ", Id)
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
	err := fwp.api.Client.Delete(createUrl(fwp.api, "firewall_policies", fwp.Id), &result, http.StatusOK)
	if err != nil {
		return nil, err
	}
	result.api = fwp.api
	return result, nil
}

// PUT /firewall_policies/{id}

// GET /firewall_policies/{id}/server_ips

// PUT /firewall_policies/{id}/server_ips

// GET /firewall_policies/{id}/server_ips/{id}

// DELETE /firewall_policies/{id}/server_ips/{id}

// GET /firewall_policies/{id}/rules

// PUT /firewall_policies/{id}/rules

// GET /firewall_policies/{id}/rules/{id}

// DELETE /firewall_policies/{id}/rules/{id}
