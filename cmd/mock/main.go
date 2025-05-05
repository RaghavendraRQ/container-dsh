package main

import (
	"container-dsh/internal/container"
	"fmt"
	"time"
)

func main() {
	cli := container.GetClient()
	containers, err := container.GetContainerList(cli)
	if err != nil {
		panic(err)
	}
	for _, containerId := range containers {
		fmt.Println("Container ID: ", containerId)
		go func() {
			containerData, err := container.GetContainerData(cli, containerId)
			if err != nil {
				panic(err)
			}
			fmt.Println(containerData.String())
		}()
	}
	// Wait for goroutines to finish
	time.Sleep(2 * time.Second)

}
