package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
)

var port = flag.String("port", "6379", "tcp port")

func main() {
	flag.Parse()

	// Start a TCP Listener
	print("TCP Listening on: ", ":"+*port)
	lis, err := net.Listen("tcp", ":"+*port)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Start receiving requests
	conn, err := lis.Accept()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	shutdown := gracefulShutdown(conn.Close)
	defer shutdown() // Shutdown gracefully once finished

	// Infinite loop for receiving commands from client and responding to them
	for {
		size := 1024
		buf := make([]byte, size)

		// Read message from client
		if _, err := conn.Read(buf); err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			fmt.Println(err)
		}

		// Send back 'PONG'
		msg := []byte("+OK\r\n")
		if _, err := conn.Write(msg); err != nil {
			fmt.Println(err)
		}
	}
}

func gracefulShutdown(closeConn func() error) func() {
	return func() {
		if err := closeConn(); err != nil {
			fmt.Println(err)
		}
		os.Exit(1)
	}
}
