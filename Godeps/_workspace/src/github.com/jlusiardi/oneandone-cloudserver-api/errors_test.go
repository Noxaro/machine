/*
 * Copyright 2015 1&1 Internet AG, http://1und1.de . All rights reserved. Licensed under the Apache v2 License.
 */

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
