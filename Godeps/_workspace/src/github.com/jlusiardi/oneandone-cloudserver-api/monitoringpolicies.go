package oneandone_cloudserver_api

import ()

type MonitoringPolicy struct {
	withId
	withName
	LastStatus []MonitoringStatus `json:"last_status"`
}

type MonitoringStatus struct {
	Resource string          `json:"resource"`
	Date     string          `json:"date"`
	State    string          `json:"state"`
	Value    MonitoringValue `json:"value"`
}

type MonitoringValue struct {
	Unit   string           `json:"unit"`
	Value  string           `json:"value"`
	Limits MonitoringLimits `json:"limits"`
}

type MonitoringLimits struct {
	Upper int `json:"upper"`
	Lower int `json:"lower"`
}

type MonitoringPolicyCreateData struct {
}

// GET /monitoring_policies
func (api *API) GetMonitoringPolicies() []MonitoringPolicy {
	return []MonitoringPolicy{}
}

// POST /monitoring_policies
func (api *API) CreateMonitoringPolicy(configuration MonitoringPolicyCreateData) MonitoringPolicy {
	return MonitoringPolicy{}
}

// GET /monitoring_policies/{id}
func (api *API) GetMonitoringPolicy(Id string) MonitoringPolicy {
	return MonitoringPolicy{}
}

// DELETE /monitoring_policies/{id}
func (mp *MonitoringPolicy) Delete() MonitoringPolicy {
	return MonitoringPolicy{}
}

// PUT /monitoring_policies/{id}

// GET /monitoring_policies/{id}/ports

// PUT /monitoring_policies/{id}/ports

// GET /monitoring_policies/{id}/ports/{id}

// DELETE /monitoring_policies/{id}/ports/{id}

// PUT /monitoring_policies/{id}/ports/{id}

// GET /monitoring_policies/{id}/processes

// PUT /monitoring_policies/{id}/processes

// GET /monitoring_policies/{id}/processes/{id}

// DELETE /monitoring_policies/{id}/processes/{id}

// PUT /monitoring_policies/{id}/processes/{id}

// GET /monitoring_policies/{id}/servers

// PUT /monitoring_policies/{id}/servers

// GET /monitoring_policies/{id}/servers/{id}

// DELETE /monitoring_policies/{id}/servers/{id}
