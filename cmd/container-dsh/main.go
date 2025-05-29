package main

import (
	"container-dsh/cmd/cli"
	"container-dsh/cmd/mock"
	"container-dsh/cmd/server"
	"flag"
	"fmt"
	"os"
)

const (
	modeUsage = "Any one of the server, logger, cli, mock\nDefault is \"Server\""
	modeError = "Usage: " + modeUsage
)

func main() {
	mode := flag.String("mode", "server", modeUsage)

	flag.Usage = func() {
		fmt.Printf("Usage: %s [--mode=server|cli|mock|logger]\n", os.Args[0])
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
	default:
		fmt.Println(modeError)
		os.Exit(1)
	}
}
