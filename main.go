package main

import (
	"fmt"
	"gosyslog/server"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s hostname\n", os.Args[0])
		os.Exit(0)
	}

	host := os.Args[1]
	server.Start(host)
}
