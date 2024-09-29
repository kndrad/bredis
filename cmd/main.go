package main

import (
	"flag"
	"fmt"
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

	defer conn.Close() // Close once finished

	// Infinite loop for receiving commands from client and responding to them
}
