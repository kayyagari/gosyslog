package main

import (
	"fmt"
	_ "gosyslog/message"
	"gosyslog/server"
	"os"
)

func main() {
	host := "localhost:6514"

	if len(os.Args) >= 2 {
		//fmt.Fprintf(os.Stderr, "Usage: %s hostname\n", os.Args[0])
		//os.Exit(0)
		host = os.Args[1]
	}

	fmt.Println("Using host and port ", host)
	server.Start(host)
}
