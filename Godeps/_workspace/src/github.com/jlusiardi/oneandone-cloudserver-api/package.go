package oneandone_cloudserver_api

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/jmcvetta/napping"
	"net/http"
)

type API struct {
	AuthToken string
	Endpoint  string
}

type withApi struct {
	api *API
}

type withId struct {
	Id string `json:"id"`
}

type withType struct {
	Type string `json:"type"`
}

type withName struct {
	Name string `json:"name"`
}

type withDescription struct {
	Description string `json:"description"`
}

type Status struct {
	State   string `json:"state"`
	Percent string `json:"percent"`
}

func New(token string, url string) *API {
	api := new(API)
	api.AuthToken = token
	api.Endpoint = url
	return api
}

func (api *API) prepareSession() *napping.Session {
	s := new(napping.Session)
	h := &http.Header{}
	h.Set("X_TOKEN", api.AuthToken)
	s.Header = h
	return s
}

func logResult(response *napping.Response, expectedStatus int) {
	if response != nil {
		log.Debug("response Status:", response.Status())
		if response.Status() != expectedStatus {
			log.Debug("response:", response.RawText())
			log.Debug("response Headers:", response.HttpResponse().Header)
		}
	} else {
		log.Info("No response was given")
	}
}

func createUrl(api *API, sections ...interface{}) string {
	url := api.Endpoint
	for _, section := range sections {
		url += "/" + fmt.Sprint(section)
	}
	return url
}
