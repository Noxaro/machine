package oneandone_cloudserver_api

import (

)

type FirewallPolicy struct {
}

type FirewallPolicyCreateData struct {
}

// GET /firewall_policies
func (api *API) GetFirewallPolicies() []FirewallPolicy {
	return []FirewallPolicy{}
}

// POST /firewall_policies
func (api *API) CreateFirewallPolicy(configuration FirewallPolicyCreateData) FirewallPolicy {
	return FirewallPolicy{}
}

// GET /firewall_policies/{id}
func (api *API) GetFirewallPolicy(Id string) FirewallPolicy {
	return FirewallPolicy{}
}

// DELETE /firewall_policies/{id}
func (fwp *FirewallPolicy) Delete() FirewallPolicy {
	return FirewallPolicy{}
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

