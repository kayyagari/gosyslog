package server

import (
	"fmt"
	"net"
	"os"
)

func Start(name string) {
	addrs, err := net.LookupHost(name)
	if err != nil {
		fmt.Println("Error while resolving host", err.Error())
		os.Exit(2)
	}

	fmt.Println(addrs)
}
