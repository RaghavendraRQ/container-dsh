package main

import (
	"container-dsh/internal/container"
	"fmt"
	"sync"
	"time"

	"github.com/docker/docker/client"
)

type SnapShotData struct {
	sync.Mutex
	TimeStamp  time.Time
	Contianers []container.Container
}

func SnapShot(cli *client.Client) (*SnapShotData, error) {
	containers, err := container.GetContainerList(cli)
	if err != nil {
		return nil, fmt.Errorf("failed to get container list: %v", err)
	}
	snapShotData := SnapShotData{
		TimeStamp:  time.Now(),
		Contianers: make([]container.Container, 0),
	}

	for _, containerId := range containers {
		go func(id string) {
			containerData, err := container.GetContainerData(cli, id)
			if err != nil {
				fmt.Printf("failed to get container data for %s: %v\n", id, err)
				return
			}
			snapShotData.Lock()
			defer snapShotData.Unlock()
			snapShotData.Contianers = append(snapShotData.Contianers, containerData)
		}(containerId)
	}
	return &snapShotData, nil
}
