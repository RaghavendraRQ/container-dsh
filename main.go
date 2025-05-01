package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

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

    containers, err := cli.ContainerList(ctx, cont.ListOptions{All:true})
    if err != nil {
        log.Fatalln("Error: ", err)
    }

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
                float64(statsData.MemoryStats.Usage)/(1024*1024),
            )

        } (container.ID)

    }

    time.Sleep(3 * time.Second)
}

func calculateCPUPercent(stat *cont.StatsResponse) float64 {
    cpuDelta := float64(stat.CPUStats.CPUUsage.TotalUsage - stat.PreCPUStats.CPUUsage.TotalUsage)
    systemDelta := float64(stat.CPUStats.SystemUsage - stat.PreCPUStats.SystemUsage)
    if systemDelta > 0.0 && cpuDelta > 0.0 {
        return (cpuDelta / systemDelta) * float64(len(stat.CPUStats.CPUUsage.PercpuUsage)) * 100.0
    }
    return 0.0
}
