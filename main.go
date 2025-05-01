package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"container-dsh/aggregator"
	"container-dsh/logger"

	cont "github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

func main() {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Fatalln("Error: ", err)
	}

	defer cli.Close()

	containers, err := cli.ContainerList(ctx, cont.ListOptions{All: true})
	if err != nil {
		log.Fatalln("Error: ", err)
	}
	timeSeriesLogger := logger.TimeSeries{}
	aggregator := aggregator.NewAggregator(3 * time.Second)
	go aggregator.Start()
	for _, container := range containers {

		go func(id string) {
			stats, err := cli.ContainerStatsOneShot(ctx, id)
			if err != nil {
				log.Fatalln("Error in ContainerStatsOneShot: ", err)
			}
			defer stats.Body.Close()

			var statsData cont.StatsResponse

			if err := json.NewDecoder(stats.Body).Decode(&statsData); err != nil {
				fmt.Println("Decode err: ", err)
				return
			}
			fmt.Printf("[Container %s] CPU: %.2f%% | Mem: %.2fMB\n",
				id[:10],
				calculateCPUPercent(&statsData),
				float64(statsData.MemoryStats.Stats["Usage"])/(1024*1024),
			)
			timeSeriesLogger.Log(id[:10], "CPU", calculateCPUPercent(&statsData))
			timeSeriesLogger.Log(id[:10], "Mem", float64(statsData.MemoryStats.Stats["Usage"])/(1024*1024))

			if id[:10] == "eeb37b74c1" {
				log.Println("Container ID: ", id)
				aggregator.MetricsChannel <- logger.ContainerMetrics{
					CpuUsage: 12,
					MemUsage: float64(statsData.MemoryStats.Stats["Usage"]) / (1024 * 1024),
				}
			}

		}(container.ID)

	}

	time.Sleep(7 * time.Second)
	timeSeriesLogger.Dump("test.json")
}

func calculateCPUPercent(stat *cont.StatsResponse) float64 {
	cpuDelta := float64(stat.CPUStats.CPUUsage.TotalUsage - stat.PreCPUStats.CPUUsage.TotalUsage)
	systemDelta := float64(stat.CPUStats.SystemUsage - stat.PreCPUStats.SystemUsage)
	if systemDelta > 0.0 && cpuDelta > 0.0 {
		return (cpuDelta / systemDelta) * float64(len(stat.CPUStats.CPUUsage.PercpuUsage)) * 100.0
	}
	return 0.0
}
