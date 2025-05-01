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
	ticker := time.NewTicker(a.interval)
	defer ticker.Stop()
	file, err := os.OpenFile("fck.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error opening file: ", err)
		return
	}

	go func() {
		for {
			select {
			case metric := <-a.MetricsChannel:
				a.updateStats(metric)
			case <-ticker.C:
				fmt.Println("Aggregated Stats:")
				for _, met := range Metrics {
					if stat, ok := a.Stats[met]; ok {
						file.WriteString(fmt.Sprintf("Metric: %s, Max: %.2f, Min: %.2f, Avg: %.2f, Total: %d\n",
							stat.Metric, stat.Max, stat.Min, stat.Avg, stat.Total))
					} else {
						fmt.Printf("Metric: %s, No data available\n", met)
					}
				}
				a.Stats = make(map[string]AggregatorStats)

			}
		}
	}()
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
