package main

import (
	"testing"
)

func TestRegisterEndpoint(t *testing.T){
	router := registerEndpoints()
	if router == nil {
		t.Error("Gorilla Mux Router creation failed.")
	}
}

func TestDockerClientRegistration(t *testing.T){
	registerDockerServiceClient()
	if DockerServiceClient == nil {
		t.Error("Docker Service Client bootstrap failed.")
	}
}
