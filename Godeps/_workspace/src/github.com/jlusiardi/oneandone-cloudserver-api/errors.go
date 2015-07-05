/*
 * Copyright 2015 1&1 Internet AG, http://1und1.de . All rights reserved. Licensed under the Apache v2 License.
 */

package oneandone_cloudserver_api

import (
	"fmt"
)

type ApiError struct {
	httpStatusCode int
	message        string
}

func (e ApiError) Error() string {
	return fmt.Sprintf("%d - %s", e.httpStatusCode, e.message)
}

func (e *ApiError) HttpStatusCode() int {
	return e.httpStatusCode
}

func (e *ApiError) Message() string {
	return e.message
}
