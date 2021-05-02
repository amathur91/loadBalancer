package main

import "time"

type Service struct {
	name string
	path string
	servicePort int
	serviceStatus bool
	backendConfig *Backend
	lastUpdated time.Time
	ipAddress string
	port uint16
}
