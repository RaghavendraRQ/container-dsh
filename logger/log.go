package logger

import (
	"encoding/json"
	"os"
)

type ContainerMetrics struct {
	CpuUsage float64 `json:"cpu_usage"`
	MemUsage float64 `json:"mem_usage"`
	NetIO    float64 `json:"net_io"`
	DiskIO   float64 `json:"disk_io"`
}

type ContainerLog struct {
	Entry map[string]ContainerMetrics `json:"entry"`
}

func (cl *ContainerLog) Log(id string, cpuUsage, memUsage, netIO, diskIO float64) {
	if cl.Entry == nil {
		cl.Entry = make(map[string]ContainerMetrics)
	}
	containerMeterics := ContainerMetrics{
		CpuUsage: cpuUsage,
		MemUsage: memUsage,
		NetIO:    netIO,
		DiskIO:   diskIO,
	}
	cl.Entry[id] = containerMeterics
}

func (cl *ContainerLog) Dump(filepath string) {
	file, err := os.OpenFile(filepath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic("Error opening file: " + err.Error())
	}
	defer file.Close()

	if err := json.NewEncoder(file).Encode(cl.Entry); err != nil {
		panic("Error encoding JSON: " + err.Error())
	}
	cl.Entry = make(map[string]ContainerMetrics)
}
