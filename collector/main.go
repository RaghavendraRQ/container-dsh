package collector

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

type Metrics struct {
	CPUUsage float64 `json:"cpu_usage"`
	MemUsage float64 `json:"mem_usage"`
	NetIO    float64 `json:"net_io"`
	DiskIO   float64 `json:"disk_io"`
}

var (
	ctx       context.Context
	statsPool sync.Pool
)

func init() {
	// Initialize the collector
	ctx = context.Background()
	statsPool = sync.Pool{
		New: func() interface{} {
			return &container.StatsResponse{}
		},
	}
}

func Start() error {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return fmt.Errorf("error in creating client: %v", err)
	}
	defer cli.Close()
	containers, err := cli.ContainerList(ctx, container.ListOptions{All: true})
	if err != nil {
		return fmt.Errorf("error in getting containers: %v", err)
	}
	GetContainerData(cli, containers)
	return nil
}

func GetContainerData(cli *client.Client, containers []container.Summary) {
	for _, cont := range containers {

		go func(id string) {
			stats, err := cli.ContainerStatsOneShot(ctx, id)
			if err != nil {
				fmt.Println("Error in ContainerStatsOneShot: ", err)
				return
			}
			defer stats.Body.Close()

			statsData := statsPool.Get().(*container.StatsResponse)
			if err := json.NewDecoder(stats.Body).Decode(statsData); err != nil {
				fmt.Println("Decode err: ", err)
				return
			}

			metrics := NewMetrics(statsData)
			PrintPretty(id, metrics)

			statsPool.Put(statsData)
		}(cont.ID)
	}
}

func NewMetrics(statsData *container.StatsResponse) Metrics {
	return Metrics{
		CPUUsage: calculateCPUPercent(statsData),
		MemUsage: calculateMemUsage(statsData),
		NetIO:    calculateNetIO(statsData),
		DiskIO:   calculateDiskIO(statsData),
	}

}

func calculateCPUPercent(stat *container.StatsResponse) float64 {
	cpuDelta := float64(stat.CPUStats.CPUUsage.TotalUsage - stat.PreCPUStats.CPUUsage.TotalUsage)
	systemDelta := float64(stat.CPUStats.SystemUsage - stat.PreCPUStats.SystemUsage)
	if systemDelta > 0.0 && cpuDelta > 0.0 {
		return (cpuDelta / systemDelta) * float64(len(stat.CPUStats.CPUUsage.PercpuUsage)) * 100.0
	}
	return 0.0
}

func calculateMemUsage(stat *container.StatsResponse) float64 {
	return float64(stat.MemoryStats.Usage) / (1024 * 1024)
}

func calculateNetIO(stat *container.StatsResponse) float64 {
	if len(stat.Networks) > 0 {
		for _, net := range stat.Networks {
			return float64(net.RxBytes+net.TxBytes) / (1024 * 1024)
		}
	}
	return 0.0
}
func calculateDiskIO(stat *container.StatsResponse) float64 {
	if len(stat.BlkioStats.IoServiceBytesRecursive) > 0 {
		return float64(stat.BlkioStats.IoServiceBytesRecursive[0].Value) / (1024 * 1024)
	}
	return 0.0
}
