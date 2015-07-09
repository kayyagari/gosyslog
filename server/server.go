package server

import (
	"fmt"
	"net"
	"os"
)

func Start(host string) {
	listener, err := net.Listen("tcp", host)
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

	var buf [128]byte
	for {
		n, rerr := conn.Read(buf[0:])
		if rerr != nil {
			if rerr.Error() != "EOF" {
				fmt.Println("Failed to read the data", rerr.Error())
			}
			return
		}

		_, werr := conn.Write(buf[0:n])
		if werr != nil {
			fmt.Println("Failed to write the data", werr.Error())
			return
		}
	}
}

func close(conn net.Conn) {
	fmt.Println("closing the connection")
	conn.Close()
}
