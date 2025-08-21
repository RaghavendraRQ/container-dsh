package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/raghavendrarq/container-dsh/cmd/cli"
	"github.com/raghavendrarq/container-dsh/cmd/mock"
	"github.com/raghavendrarq/container-dsh/cmd/prom"
	"github.com/raghavendrarq/container-dsh/cmd/server"
	"github.com/raghavendrarq/container-dsh/internal/container"
)

const (
	modeUsage = "Any one of the server, logger, cli, mock\nDefault is \"Server\""
	modeError = "Usage: " + modeUsage
)

func main() {
	mode := flag.String("mode", "server", modeUsage)

	flag.Usage = func() {
		fmt.Printf("Usage: %s [--mode=server|cli|mock|logger|prometheus]\n", os.Args[0])
		fmt.Println("Description:\n Starts different components of the container dashboard.")
		flag.PrintDefaults()
	}

	flag.Parse()

	switch *mode {
	case "server":
		server.Run()
	case "logger":
		fmt.Println("Work In Progress please come back soon")
	case "cli":
		cli.Run()
	case "mock":
		mock.Run()
	case "prometheus":
		prom.Run()

	case "redis":
		TestCache()
	default:
		fmt.Println(modeError)
		os.Exit(1)
	}
}
func TestCache() {
	cli := container.GetClient()
	for i := range 10 {
		containerIds, err := container.GetContainerList(cli)
		if err != nil {
			fmt.Println("Error", err)
		}
		fmt.Printf("ContainerIds(%d): %v\n", i, containerIds)

	}
}
