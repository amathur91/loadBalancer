package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseConfig(t *testing.T){
	configResponse := ParseConfig("config.yml")
	assert.NotNil(t, configResponse)
	assert.NotNil(t, configResponse.Backends)
	assert.NotNil(t, configResponse.Routes)
}
