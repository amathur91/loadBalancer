package main

import (
	"context"
	"github.com/docker/docker/api/types"
	"strings"
	"time"
)

func UpdateServices(){
	for ;; {
		InfoLogger.Println("Scanning Services for updates.")
		containers, err := DockerServiceClient.ContainerList(context.Background(), types.ContainerListOptions{})
		if err != nil {
			ErrorLogger.Println("Unable to connect to Docker Client. Service may not reflect correct state.")
			panic(err)
		}
		updateContainerServiceInfo(containers)
		// For the services whose container died we need to check for the last updated and set value false.
		// We are assuming that docker container scanning should finish under 1 min
		updateUnreachableServices()
		InfoLogger.Printf("Scanning for services complete. Sleeping %.1f minutes \n", ServiceAvailabilityThresholdInMinutes)
		time.Sleep(time.Duration(ServiceDiscoverySleepDurationMinutes) * time.Minute)
	}
}

func updateUnreachableServices() {
	for _, service := range RegisteredServices {
		lastUpdatedTimeStamp := service.lastUpdated
		currentTime := time.Now()
		if currentTime.Sub(lastUpdatedTimeStamp).Minutes() > ServiceAvailabilityThresholdInMinutes {
			service.serviceStatus = false
			RegisteredServices[service.name] = service
			EndPointToServiceMap[service.path] = &service
			InfoLogger.Printf("Service %s is offline.", service.name)
		}
	}
}

func updateContainerServiceInfo(containers []types.Container) {
	// For each running container update the service last timestamp for availability check
	for _, container := range containers {
		containerName := container.Names[0]
		containerName = strings.Split(containerName, "/")[1]
		service, present := RegisteredServices[containerName]
		if present {
			InfoLogger.Printf("Updating Service : %s \n", containerName)
			dockerLabels := container.Labels
			nameValue, appNameMatch := dockerLabels[AppNameLabel]
			envValue, envMatch := dockerLabels[EnvironmentLabel]
			if appNameMatch && envMatch && nameValue == service.backendConfig.MatchLabels.AppName &&
				envValue == service.backendConfig.MatchLabels.Env {
				service.lastUpdated = time.Now()
				service.serviceStatus = true
				service.ipAddress = container.Ports[0].IP
				service.port = container.Ports[0].PublicPort
				RegisteredServices[containerName] = service
				EndPointToServiceMap[service.path] = &service
				InfoLogger.Printf("Service %s is refreshed with updated info. \n", containerName)
			}
		}
	}
}