package main

import (
	"flag"
	"fmt"
	"github.com/docker/docker/client"
	"github.com/gorilla/mux"
	"github.com/spenczar/tdigest"
	"log"
	"net/http"
	"os"
	"strconv"
)

var ServicePort int
var ServiceConfigFile string
var RegisteredServices = make(map[string]Service, 0)
var EndPointToServiceMap = make(map[string]*Service, 0)
var Configuration *Config
var DockerServiceClient *client.Client
var SuccessCalls uint32
var FailedCalls uint32
var P95 float64
var P99 float64
var Avg float64
var Digest *tdigest.TDigest


func startLoadBalancer(){
	InfoLogger.Println("Bootstrapping LoadBalancer")
	ParseFlags()
	setupExecutionStatsDigest()
	registerDockerServiceClient()
	Configuration = ParseConfig(ServiceConfigFile)
	updateRegisteredServices()
	route := registerEndpoints()
	go UpdateServices()
	log.Fatalln(http.ListenAndServe(":"+strconv.Itoa(ServicePort), route))
}

func setupExecutionStatsDigest() {
	Digest = tdigest.New()
	Digest.Add(float64(0), 1)
}


func registerEndpoints() *mux.Router {
	route := mux.NewRouter().StrictSlash(true)
	// Register all the endpoint path mentioned in the configuration.
	for _, service := range RegisteredServices {
		InfoLogger.Printf("Registering Service %s with Service path : %s \n", service.name, service.path)
		route.HandleFunc(service.path, DelegateBackendRequest)
	}
	route.HandleFunc(StatisticsEndpoint, GetAPIStatistics)
	return route
}

func registerDockerServiceClient(){
	client, err := client.NewClientWithOpts()
	if err != nil {
		panic(err)
	}
	DockerServiceClient = client
}

func updateRegisteredServices() {
	routingDetails := Configuration.Routes
	// Mapping Routes with Backend
	for serviceIndex := 0; serviceIndex < len(routingDetails); serviceIndex++ {
		service := Service{
			name: routingDetails[serviceIndex].Backend,
			path: routingDetails[serviceIndex].PathPrefix}
		RegisteredServices[routingDetails[serviceIndex].Backend] = service
		EndPointToServiceMap[routingDetails[serviceIndex].PathPrefix] = &service
	}

	// Mapping match labels with backend
	backendDetails := Configuration.Backends
	for backendIndex := 0; backendIndex < len(backendDetails); backendIndex++ {
		serviceName := backendDetails[backendIndex].Name
		service, ok := RegisteredServices[serviceName]
		if ok {
			service.backendConfig = &backendDetails[backendIndex]
			RegisteredServices[serviceName] = service
			InfoLogger.Printf("Service %s registered for Traffic Routing. \n", serviceName)
		}
	}
}

func ParseFlags() {
	port := flag.Int("port", -1, "API Gateway service port.")
	configFile := flag.String("config", "", "Configuration of backend services")
	flag.Parse()
	if port == nil || *port == -1 {
		fmt.Println("Port for API Gateway is not correctly specified. Exiting 1")
		os.Exit(1)
	}
	if configFile == nil || *configFile == "" {
		fmt.Println("Config File is not specified")
		os.Exit(1)
	}
	ServicePort, ServiceConfigFile = *port, *configFile
}
