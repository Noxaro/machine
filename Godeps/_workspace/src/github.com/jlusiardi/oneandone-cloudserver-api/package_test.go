package oneandone_cloudserver_api

import (
	assert "github.com/stretchr/testify/assert"
	"testing"
)

func TestNew(t *testing.T) {
	api := New("134", "abc")
	
	assert.Equal(t, "134", api.AuthToken)
	assert.Equal(t, "abc", api.Endpoint)
}

func TestPrepareSession(t *testing.T) {
	api := New("134", "abc")
	session := api.prepareSession()
	assert.Equal(t, "134", session.Header.Get("X_TOKEN"))
}

func TestCreateUrl_1(t *testing.T) {
	api := New("token", "http://test.de/v1")
	result := createUrl(api)
	assert.Equal(t, "http://test.de/v1", result)
}

func TestCreateUrl_2(t *testing.T) {
	api := New("token", "http://test.de/v1")

	result := createUrl(api, "servers")
	assert.Equal(t, "http://test.de/v1/servers", result)
}

func TestCreateUrl_3(t *testing.T) {
	api := New("token", "http://test.de/v1")

	result := createUrl(api, "servers", 1)
	assert.Equal(t, "http://test.de/v1/servers/1", result)
}
