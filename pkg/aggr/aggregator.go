package aggr

import (
	"container-dsh/internal/container"
	"fmt"
	"log"
	"os"
	"sync"
	"time"
)

type AggregatorMetrics struct {
	Count int
	Min   float64
	Max   float64
	Avg   float64
}

type Aggregator struct {
	sync.WaitGroup
	mu             sync.Mutex
	AggregatedData map[string]*AggregatorMetrics
	filePath       string
	Completed      chan bool
	MetricsChannel chan container.Container
	Interval       time.Duration
}

func NewAggregator(filePath string, inteval time.Duration) *Aggregator {
	return &Aggregator{
		AggregatedData: make(map[string]*AggregatorMetrics),
		filePath:       filePath,
		Completed:      make(chan bool),
		MetricsChannel: make(chan container.Container),
		Interval:       inteval,
	}
}

func (a *Aggregator) Start() {
	a.Add(1)
	defer a.Done()

	ticker := time.NewTicker(a.Interval)
	defer ticker.Stop()

	for {
		select {
		case <-a.Completed:
			a.Dump()
			log.Println("Aggregator completed")
			return
		case metric, ok := <-a.MetricsChannel:
			if !ok {
				a.Dump()
				log.Println("Aggregator channel closed")
				return
			}
			a.mu.Lock()
			a.aggregate(metric)
			a.mu.Unlock()
		case <-ticker.C:
			a.mu.Lock()
			a.Dump()
			a.mu.Unlock()
			log.Println("Dumping aggregated data")
		}
	}
}

func (a *Aggregator) aggregate(metric container.Container) {
	if _, exsists := a.AggregatedData[metric.ID]; !exsists {
		a.AggregatedData[metric.ID] = &AggregatorMetrics{
			Count: 0,
			Min:   metric.CpuUsage,
			Max:   metric.CpuUsage,
			Avg:   metric.CpuUsage,
		}
	}
	a.AggregatedData[metric.ID].Count++
	if metric.CpuUsage < a.AggregatedData[metric.ID].Min {
		a.AggregatedData[metric.ID].Min = metric.CpuUsage
	}
	if metric.CpuUsage > a.AggregatedData[metric.ID].Max {
		a.AggregatedData[metric.ID].Max = metric.CpuUsage
	}
	a.AggregatedData[metric.ID].Avg = (a.AggregatedData[metric.ID].Avg + metric.CpuUsage) / float64(a.AggregatedData[metric.ID].Count)
}

func (a *Aggregator) Dump() {
	file, err := os.OpenFile(a.filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println("Error opening file: ", err)
		return
	}
	defer file.Close()

	for id, data := range a.AggregatedData {
		file.WriteString(
			fmt.Sprintf("Container ID: %s, Count: %d, Min: %.2f, Max: %.2f, Avg: %.2f\n",
				id,
				data.Count,
				data.Min,
				data.Max,
				data.Avg,
			),
		)
	}
}
