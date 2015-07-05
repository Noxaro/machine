/*
 * Copyright 2015 1&1 Internet AG, http://1und1.de . All rights reserved. Licensed under the Apache v2 License.
 */

package oneandone_cloudserver_api

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"strconv"
)

func TestNew(t *testing.T) {
	api := New("134", "abc")

	assert.Equal(t, "abc", api.Endpoint)
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

func TestInt2Pointer_1(t *testing.T) {
	result := Int2Pointer(42)

	assert.Equal(t, 42, *result)
}

func TestGetMapKeysString(t *testing.T) {
	tMap := make(map[string]int)
	for i := 0; i <= 10; i++ {
		tMap["test_" + strconv.Itoa(i)] = 1
	}
	keys := getMapKeysString(tMap)
	assert.Equal(t, len(tMap), len(keys))

	for _, value := range keys {
		assert.Equal(t, 1, tMap[value])
		delete(tMap, value)
	}

	assert.Equal(t, 0, len(tMap))
}

func TestGetMapKeysInt(t *testing.T) {
	tMap := make(map[int]int)
	for i := 0; i <= 10; i++ {
		tMap[i] = 1
	}
	keys := getMapKeysInt(tMap)
	assert.Equal(t, len(tMap), len(keys))

	for _, value := range keys {
		assert.Equal(t, 1, tMap[value])
		delete(tMap, value)
	}

	assert.Equal(t, 0, len(tMap))
}

func TestGetMapKeysStringNil(t *testing.T) {
	result := getMapKeysString(nil)
	assert.Equal(t, []string{}, result)
}

func TestGetMapKeysIntNil(t *testing.T) {
	result := getMapKeysInt(nil)
	assert.Equal(t, []int{}, result)
}

func TestGetMapKeysStringEmpty(t *testing.T) {
	result := getMapKeysString(make(map[string]int))
	assert.Equal(t, []string{}, result)
}

func TestGetMapKeysIntEmpty(t *testing.T) {
	result := getMapKeysInt(make(map[int]int))
	assert.Equal(t, []int{}, result)
}

