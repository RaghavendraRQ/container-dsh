package aggr

import (
	"container-dsh/internal/container"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"
	"time"
)

type Status string

// MetricAggregate holds aggregated values for one metric.
type MetricAggregate struct {
	Sum   float64 `json:"sum"`
	Avg   float64 `json:"avg"`
	Min   float64 `json:"min"`
	Max   float64 `json:"max"`
	Count int64   `json:"count"`
}

// AggregateResult is the JSON structure for aggregated stats output.
type AggregateResult struct {
	Timestamp   time.Time       `json:"timestamp"`
	ContainerID string          `json:"container_id"`
	CPU         MetricAggregate `json:"cpu_usage"`
	Memory      MetricAggregate `json:"memory_usage"`
	Network     MetricAggregate `json:"net_io"`
	Disk        MetricAggregate `json:"disk_io"`
}

// ContainerAggregator collects stats for a single container and outputs aggregates.
type ContainerAggregator struct {
	containerID string
	windowSec   int // N seconds aggregation window
	writeSec    int // M seconds JSON output interval
	inputCh     chan container.Stats
	quitCh      chan struct{}
	mu          sync.Mutex // protects counters below

	// Running counters for the current window:
	count                           int64
	sumCPU, sumMem, sumNet, sumDisk float64
	minCPU, minMem, minNet, minDisk float64
	maxCPU, maxMem, maxNet, maxDisk float64

	lastResult AggregateResult // most recent computed aggregate
}

// NewContainerAggregator creates and starts a new aggregator for a container.
func NewContainerAggregator(containerID string, windowSec, writeSec int) *ContainerAggregator {
	// Initialize min values to +Inf so first value sets them correctly
	ca := &ContainerAggregator{
		containerID: containerID,
		windowSec:   windowSec,
		writeSec:    writeSec,
		inputCh:     make(chan container.Stats),
		quitCh:      make(chan struct{}),
		minCPU:      1e308, minMem: 1e308,
		minNet: 1e308, minDisk: 1e308,
	}
	go ca.run() // start the internal loop
	return ca
}

// run processes incoming Stats and triggers aggregation on time intervals.
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
			// Update min/max
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
			// Compute aggregate for the last window and reset counters.
			ca.computeAggregates()

		case <-writeTicker.C:
			// Serialize the last aggregated result to JSON file.
			ca.writeJSON()

		case <-ca.quitCh:
			// Stop signal received: exit goroutine.
			log.Print("Stopping aggregator for container:", ca.containerID)
			return
		}
	}
}

// computeAggregates calculates averages/min/max and resets the counters.
func (ca *ContainerAggregator) computeAggregates() {
	ca.mu.Lock()
	defer ca.mu.Unlock()
	if ca.count == 0 {
		// Nothing to do if no data arrived in this window.
		return
	}
	// Build MetricAggregate for each resource.
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
	// Record the timestamped result.
	ca.lastResult = AggregateResult{
		Timestamp:   time.Now(),
		ContainerID: ca.containerID,
		CPU:         cpuAgg,
		Memory:      memAgg,
		Network:     netAgg,
		Disk:        diskAgg,
	}
	// Reset counters for next window.
	ca.count = 0
	ca.sumCPU, ca.sumMem, ca.sumNet, ca.sumDisk = 0, 0, 0, 0
	ca.minCPU, ca.minMem, ca.minNet, ca.minDisk = 1e308, 1e308, 1e308, 1e308
	ca.maxCPU, ca.maxMem, ca.maxNet, ca.maxDisk = 0, 0, 0, 0
}

// writeJSON serializes the lastResult to a JSON file (indented for readability).
func (ca *ContainerAggregator) writeJSON() {
	result := ca.lastResult
	// Use MarshalIndent for pretty JSON output:contentReference[oaicite:7]{index=7}.
	data, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		fmt.Printf("Error marshaling JSON for container %s: %v\n", ca.containerID, err)
		return
	}
	// Write to a file (one per container, or could append timestamp).
	filename := fmt.Sprintf("%s_%d.json", ca.containerID, time.Now().Unix())
	log.Print("Writing JSON file:", filename)
	if err := os.WriteFile(filename, data, 0644); err != nil {
		fmt.Printf("Error writing JSON file for container %s: %v\n", ca.containerID, err)
		return
	}
}

// Stop signals the aggregator to shut down (ticker.Stop() is called via defer in run).
func (ca *ContainerAggregator) Stop() {
	close(ca.quitCh)
}

// AggregatorManager routes container stats to per-container aggregators.
type AggregatorManager struct {
	windowSec   int
	writeSec    int
	mu          sync.Mutex
	aggregators map[string]*ContainerAggregator
	Input       chan container.Container // Channel for incoming container metrics.
}

// NewAggregatorManager initializes the manager with given intervals.
func NewAggregatorManager(windowSec, writeSec int) *AggregatorManager {
	return &AggregatorManager{
		windowSec:   windowSec,
		writeSec:    writeSec,
		aggregators: make(map[string]*ContainerAggregator),
		Input:       make(chan container.Container),
	}
}

// Run listens on the Input channel and dispatches each Container to its aggregator.
func (am *AggregatorManager) Run() {
	for c := range am.Input {
		am.mu.Lock()
		agg, exists := am.aggregators[c.ID]
		if !exists {
			// Create a new aggregator if this container ID is new.
			agg = NewContainerAggregator(c.ID, am.windowSec, am.writeSec)
			am.aggregators[c.ID] = agg
		}
		am.mu.Unlock()
		// Send stats to the container's aggregator goroutine.
		agg.inputCh <- c.Stats
	}
}

// Stop gracefully shuts down all aggregators and closes the Input channel.
func (am *AggregatorManager) Stop() {
	am.mu.Lock()
	defer am.mu.Unlock()
	close(am.Input) // Stop accepting new data; causes Run() loop to exit.
	for _, agg := range am.aggregators {
		agg.Stop()
	}
}
