package main

import (
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
		fmt.Printf("Usage of %s:\n", os.Args[0])
		fmt.Println(modeUsage)
		flag.PrintDefaults()
	}

	flag.Parse()

	switch *mode {
	case "server":
		server.Run()
	case "logger":
		fmt.Println("Work In Progress please come back soon")
	case "cli":
		fmt.Println("Work In Progress please come back soon")
	case "mock":
		mock.Run()
	default:
		fmt.Println(modeError)
		os.Exit(1)
	}
}
