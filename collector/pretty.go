package collector

import "fmt"

func PrintPretty(id string, metric Metrics) {
	fmt.Printf("[Container %s] CPU: %.2f%% | Mem: %.2fMB | NetIO: %.2fMB | DiskIO: %.2fMB\n",
		id[:10],
		metric.CPUUsage,
		metric.MemUsage,
		metric.NetIO,
		metric.DiskIO,
	)
}
