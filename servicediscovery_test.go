package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestUpdateUnreachableServices(t *testing.T){
	currentTime := time.Now()
	oldTime := currentTime.AddDate(0,0,-5)
	service1 := Service{name: "demo", serviceStatus: true, lastUpdated: oldTime}
	RegisteredServices["demo"] = service1
	updateUnreachableServices()
	assert.False(t, RegisteredServices["demo"].serviceStatus)
}

func TestUpdateReachableServices(t *testing.T){
	currentTime := time.Now()
	service1 := Service{name: "demo", serviceStatus: true, lastUpdated: currentTime}
	RegisteredServices["demo"] = service1
	updateUnreachableServices()
	assert.True(t, RegisteredServices["demo"].serviceStatus)
}