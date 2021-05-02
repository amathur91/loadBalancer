package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"
)

func DelegateBackendRequest(w http.ResponseWriter, req *http.Request){
	service, available := EndPointToServiceMap[req.RequestURI]
	InfoLogger.Printf("Received API request for URL : %s \n", req.RequestURI)
	if available && service.serviceStatus {
		req.Host = service.ipAddress + ":" + strconv.Itoa(int(service.port))
		req.URL.Host = service.ipAddress + ":" + strconv.Itoa(int(service.port))
		req.URL.Scheme = "http"
		req.URL.Path = ""
		req.RequestURI = ""
		startTime := time.Now()
		response, err := http.DefaultClient.Do(req)
		executionTime := time.Since(startTime).Milliseconds()
		if err != nil {
			ErrorLogger.Printf("Error Received from Backend Service %s \n", service.name)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, err)
			go UpdateStatMap(false, executionTime)
			return
		}
		w.WriteHeader(response.StatusCode)
		io.Copy(w, response.Body)
		go UpdateStatMap(true, executionTime)

	} else {
		InfoLogger.Printf("Service %s is unavailable. \n", service.name)
		w.WriteHeader(Configuration.DefaultResponse.StatusCode)
		fmt.Fprintf(w, Configuration.DefaultResponse.Body)
		go UpdateStatMap(false, -1)
	}
}

func GetAPIStatistics(w http.ResponseWriter, r *http.Request){
	stat := Stats{RequestCount: RequestCount{Success: SuccessCalls, Error: FailedCalls},
		Latency: Latency{Average: Digest.Quantile(0.5), P95: Digest.Quantile(0.95), P99: Digest.Quantile(0.99)}}
	json.NewEncoder(w).Encode(stat)
}

