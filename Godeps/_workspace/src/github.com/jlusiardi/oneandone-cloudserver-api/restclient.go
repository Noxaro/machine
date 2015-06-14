package oneandone_cloudserver_api

import (
	"bytes"
	"encoding/json"
	"github.com/docker/machine/log"
	"io"
	"io/ioutil"
	"net/http"
	"fmt"
)

type RestClient struct {
	token string
}

func NewRestClient(token string) *RestClient {
	restClient := new(RestClient)
	restClient.token = token
	return restClient
}

func (c *RestClient) Get(url string, result interface{}, expectedStatus int) (error) {
	return c.doRequest(url, "GET", nil, result, expectedStatus)
}

func (c *RestClient) Delete(url string, result interface{}, expectedStatus int) (error) {
	return c.doRequest(url, "DELETE", nil, result, expectedStatus)
}

func (c *RestClient) Post(url string, requestBody interface{}, result interface{}, expectedStatus int) (error) {
	return c.doRequest(url, "POST", requestBody, result, expectedStatus)
}

func (c *RestClient) Put(url string, requestBody interface{}, result interface{}, expectedStatus int) (error) {
	return c.doRequest(url, "PUT", requestBody, result, expectedStatus)
}

func (c *RestClient) doRequest(url string, method string, requestBody interface{}, result interface{}, expectedStatus int) (error) {
	var bodyData io.Reader
	if requestBody != nil {
		data, _ := json.Marshal(requestBody)
		bodyData = bytes.NewBuffer(data)
	}

	request, err := http.NewRequest(method, url, bodyData)
	if err != nil {
		return err
	}

	request.Header.Add("X-Token", c.token)
	request.Header.Add("Content-Type", "application/json")
	client := http.Client{}
	response, err := client.Do(request)
	if err = isError(response, expectedStatus, err); err != nil {
		return err
	}

	body, err := ioutil.ReadAll(response.Body)
	//log.Debug(string(body[:]))
	response.Body.Close()
	if err != nil {
		return err
	}
	err = json.Unmarshal(body, result)
	if err != nil {
		return err
	}

	return nil
}

func isError(response *http.Response, expectedStatus int, err error) error {
	if err != nil {
		return err
	}
	if response != nil {
		if response.StatusCode == expectedStatus {
			// we got a response with the expected HTTP status code, hence no error
			return nil
		}
		body, _ := ioutil.ReadAll(response.Body)
		// extract the API's error message to be returned later
		errorResponse := new(errorResponse)
		err = json.Unmarshal(body, errorResponse)
		if err != nil {
			log.Debug("JSON decode failed: ", err)
		}

		return &ApiError{response.StatusCode, errorResponse.Message}
	} else {
		// no response from API means generic error
	}
	return nil
}

func createUrl(api *API, sections ...interface{}) string {
	url := api.Endpoint
	for _, section := range sections {
		url += "/" + fmt.Sprint(section)
	}
	return url
}
