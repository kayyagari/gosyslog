package server

import (
	"bytes"
	"fmt"
	"net"
	"os"
)

var maxLen int = 1024

func Start(host string) {
	listener, err := net.Listen("udp", host)
	if err != nil {
		fmt.Println("Error while listening at host", err.Error())
		os.Exit(2)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		go handleClient(conn)
	}
}

func handleClient(conn net.Conn) {
	defer close(conn)

	buf := make([]byte, maxLen)
	for {
		n, rerr := conn.Read(buf)
		if rerr != nil {
			if rerr.Error() != "EOF" {
				fmt.Println("Failed to read the data", rerr.Error())
			}
			return
		}

		go Parse(bytes.NewBuffer(buf[0:n]))
	}
}

func close(conn net.Conn) {
	fmt.Println("closing the connection")
	conn.Close()
}
