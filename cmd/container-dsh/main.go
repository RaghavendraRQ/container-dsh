package main

import (
	"container-dsh/cmd/mock"
	"container-dsh/cmd/server"
	"flag"
	"fmt"
)

const (
	modeUsage = "Any one of the server, logger, cli, mock\nDefalut is \"Server\""
	modeError = "Usage: " + modeUsage
)

func main() {
	mode := flag.String("mode", "server", modeUsage)
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
	}
}
