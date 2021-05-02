package main

import (
	"sync/atomic"
)

type Stats struct {
	RequestCount RequestCount `json:"request_count"`
	Latency Latency `json:"latency"`
}

type Latency  struct {
	Average float64 `json:"average"`
	P95 float64 `json:"p_95"`
	P99 float64 `json:"p_99"`
}

type RequestCount struct{
	Success uint32 `json:"success"`
	Error uint32 `json:"error"`
}

func UpdateStatMap(success bool, executionTime int64){
	if success {
		atomic.AddUint32(&SuccessCalls, 1)
	} else {
		atomic.AddUint32(&FailedCalls, 1)
	}
	if executionTime > 0 {
		Digest.Add(float64(executionTime), 1)
	}
}


