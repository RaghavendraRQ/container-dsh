package aggregator

import (
	"container-dsh/logger"
	"fmt"
	"math"
	"os"
	"time"
)

var Metrics = []string{"cpu_usage", "mem_usage", "net_io", "disk_io"}

type AggregatorStats struct {
	Metric        string
	Max, Min, Avg float64
	Total         int
}

type Aggregator struct {
	TimeStamp      string
	interval       time.Duration
	MetricsChannel chan logger.ContainerMetrics
	Stats          map[string]AggregatorStats
}

func NewAggregator(interval time.Duration) *Aggregator {
	return &Aggregator{
		interval:       interval,
		MetricsChannel: make(chan logger.ContainerMetrics),
		Stats:          make(map[string]AggregatorStats),
	}
}

func (a *Aggregator) Start() {
	file, err := os.OpenFile("aggregator.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic("Error opening file: " + err.Error())
	}
	defer file.Close()
	for {
		select {
		case metric := <-a.MetricsChannel:
			a.updateStats(metric)
			a.TimeStamp = time.Now().Format(time.RFC3339)

		case <-time.After(a.interval):
			for _, met := range Metrics {
				if stat, ok := a.Stats[met]; ok {
					file.WriteString(fmt.Sprintf(
						"[%s] %s: Max: %.2f, Min: %.2f, Avg: %.2f, Total: %d\n",
						a.TimeStamp,
						stat.Metric,
						stat.Max,
						stat.Min,
						stat.Avg,
						stat.Total,
					))
				}
			}
			a.Stats = make(map[string]AggregatorStats)
			a.TimeStamp = time.Now().Format(time.RFC3339)
		}
	}
}
func (a *Aggregator) updateStats(metric logger.ContainerMetrics) {
	for _, met := range Metrics {
		if _, ok := a.Stats[met]; !ok {
			a.Stats[met] = AggregatorStats{
				Metric: met,
				Max:    metric.CpuUsage,
				Min:    metric.CpuUsage,
				Avg:    metric.CpuUsage,
				Total:  1,
			}
		} else {
			a.Stats[met] = AggregatorStats{
				Metric: a.Stats[met].Metric,
				Max:    math.Max(a.Stats[met].Max, metric.CpuUsage),
				Min:    math.Min(a.Stats[met].Min, metric.CpuUsage),
				Avg:    (a.Stats[met].Avg*float64(a.Stats[met].Total) + metric.CpuUsage) / float64(a.Stats[met].Total+1),
				Total:  a.Stats[met].Total + 1,
			}
		}
	}

}
