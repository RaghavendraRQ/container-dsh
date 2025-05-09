package aggr

import (
	"container-dsh/internal/container"
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"
)

type MetricAggregate struct {
	Sum   float64 `json:"sum"`
	Avg   float64 `json:"avg"`
	Min   float64 `json:"min"`
	Max   float64 `json:"max"`
	Count int64   `json:"count"`
}

type AggregateResult struct {
	Timestamp   time.Time       `json:"timestamp"`
	ContainerID string          `json:"container_id"`
	CPU         MetricAggregate `json:"cpu_usage"`
	Memory      MetricAggregate `json:"memory_usage"`
	Network     MetricAggregate `json:"net_io"`
	Disk        MetricAggregate `json:"disk_io"`
}

type ContainerAggregator struct {
	containerID string
	windowSec   int
	writeSec    int
	inputCh     chan container.Stats
	quitCh      chan struct{}
	mu          sync.Mutex

	count                           int64
	sumCPU, sumMem, sumNet, sumDisk float64
	minCPU, minMem, minNet, minDisk float64
	maxCPU, maxMem, maxNet, maxDisk float64

	lastResult AggregateResult
}

func NewContainerAggregator(containerID string, windowSec, writeSec int) *ContainerAggregator {
	ca := &ContainerAggregator{
		containerID: containerID,
		windowSec:   windowSec,
		writeSec:    writeSec,
		inputCh:     make(chan container.Stats),
		quitCh:      make(chan struct{}),
		minCPU:      1e308,
		minMem:      1e308,
		minNet:      1e308,
		minDisk:     1e308,
	}
	go ca.run()
	return ca
}

func (ca *ContainerAggregator) run() {
	aggTicker := time.NewTicker(time.Duration(ca.windowSec) * time.Second)
	writeTicker := time.NewTicker(time.Duration(ca.writeSec) * time.Second)
	defer aggTicker.Stop()
	defer writeTicker.Stop()

	for {
		select {
		case stat := <-ca.inputCh:
			ca.mu.Lock()
			ca.count++
			ca.sumCPU += stat.CpuUsage
			ca.sumMem += stat.MemUsage
			ca.sumNet += stat.NetIO
			ca.sumDisk += stat.DiskIO

			if stat.CpuUsage < ca.minCPU {
				ca.minCPU = stat.CpuUsage
			}
			if stat.CpuUsage > ca.maxCPU {
				ca.maxCPU = stat.CpuUsage
			}

			if stat.MemUsage < ca.minMem {
				ca.minMem = stat.MemUsage
			}
			if stat.MemUsage > ca.maxMem {
				ca.maxMem = stat.MemUsage
			}

			if stat.NetIO < ca.minNet {
				ca.minNet = stat.NetIO
			}
			if stat.NetIO > ca.maxNet {
				ca.maxNet = stat.NetIO
			}

			if stat.DiskIO < ca.minDisk {
				ca.minDisk = stat.DiskIO
			}
			if stat.DiskIO > ca.maxDisk {
				ca.maxDisk = stat.DiskIO
			}

			ca.mu.Unlock()

		case <-aggTicker.C:
			ca.computeAggregates()

		// TODO: Change writeJSON to send to a channel or request
		case <-writeTicker.C:
			ca.writeJSON()

		case <-ca.quitCh:
			return
		}
	}
}

func (ca *ContainerAggregator) computeAggregates() {
	ca.mu.Lock()
	defer ca.mu.Unlock()
	if ca.count == 0 {
		return
	}
	cpuAgg := MetricAggregate{
		Sum: ca.sumCPU,
		Avg: ca.sumCPU / float64(ca.count),
		Min: ca.minCPU, Max: ca.maxCPU,
		Count: ca.count,
	}
	memAgg := MetricAggregate{
		Sum: ca.sumMem,
		Avg: ca.sumMem / float64(ca.count),
		Min: ca.minMem, Max: ca.maxMem,
		Count: ca.count,
	}
	netAgg := MetricAggregate{
		Sum: ca.sumNet,
		Avg: ca.sumNet / float64(ca.count),
		Min: ca.minNet, Max: ca.maxNet,
		Count: ca.count,
	}
	diskAgg := MetricAggregate{
		Sum: ca.sumDisk,
		Avg: ca.sumDisk / float64(ca.count),
		Min: ca.minDisk, Max: ca.maxDisk,
		Count: ca.count,
	}
	ca.lastResult = AggregateResult{
		Timestamp:   time.Now(),
		ContainerID: ca.containerID,
		CPU:         cpuAgg,
		Memory:      memAgg,
		Network:     netAgg,
		Disk:        diskAgg,
	}

	ca.count = 0
	ca.sumCPU, ca.sumMem, ca.sumNet, ca.sumDisk = 0, 0, 0, 0
	ca.minCPU, ca.minMem, ca.minNet, ca.minDisk = 1e308, 1e308, 1e308, 1e308
	ca.maxCPU, ca.maxMem, ca.maxNet, ca.maxDisk = 0, 0, 0, 0
}

func (ca *ContainerAggregator) writeJSON() {
	result := ca.lastResult
	data, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		fmt.Printf("Error marshaling JSON for container %s: %v\n", ca.containerID, err)
		return
	}
	filename := fmt.Sprintf("%s_%d.json", ca.containerID, time.Now().Unix()) // Change the filename startegy
	if err := os.WriteFile(filename, data, 0644); err != nil {
		fmt.Printf("Error writing JSON file for container %s: %v\n", ca.containerID, err)
		return
	}
}

func (ca *ContainerAggregator) Stop() {
	close(ca.quitCh)
}

type AggregatorManager struct {
	windowSec   int
	writeSec    int
	mu          sync.Mutex
	aggregators map[string]*ContainerAggregator
	Input       chan container.Container
}

func NewAggregatorManager(windowSec, writeSec int) *AggregatorManager {
	return &AggregatorManager{
		windowSec:   windowSec,
		writeSec:    writeSec,
		aggregators: make(map[string]*ContainerAggregator),
		Input:       make(chan container.Container),
	}
}

func (am *AggregatorManager) Run() {
	for c := range am.Input {
		am.mu.Lock()
		agg, exists := am.aggregators[c.ID]
		if !exists {
			agg = NewContainerAggregator(c.ID, am.windowSec, am.writeSec)
			am.aggregators[c.ID] = agg
		}
		am.mu.Unlock()
		agg.inputCh <- c.Stats
	}
}

func (am *AggregatorManager) Stop() {
	am.mu.Lock()
	defer am.mu.Unlock()
	close(am.Input)
	for _, agg := range am.aggregators {
		agg.Stop()
	}
}
