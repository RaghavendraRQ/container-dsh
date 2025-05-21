package main

import (
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
		fmt.Println("You have choosen server")
	case "logger":
		fmt.Println("You have choosen logger")
	case "cli":
		fmt.Println("You have choosen cli")
	case "mock":
		fmt.Println("You have choosen mock")
	default:
		fmt.Println(modeError)
	}
}
