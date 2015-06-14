package oneandone_cloudserver_api

import (
	assert "github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateError(t *testing.T) {
	err := ApiError{httpStatusCode: 404, message: "Not found"}

	assert.Equal(t, 404, err.HttpStatusCode())
	assert.Equal(t, "Not found", err.Message())
}
