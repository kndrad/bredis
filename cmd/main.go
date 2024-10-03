package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
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
		msg := []byte("+OK (pong)\r\n")
		if _, err := conn.Write(msg); err != nil {
			fmt.Println(err)
		}
	}

	// Create input
	input := "$5\r\nMagda\r\n"
	reader := bufio.NewReader(strings.NewReader(input))

	// Read the RESP string to determine number of characters we need to read.
	// Above is 5, plus additional 2 bytes.
	b, err := reader.ReadByte()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
	if b != '$' {
		fmt.Printf("Invalid first byte type, expecting bulk strings only.\n")
	}

	// Determine number of characters in a string.
	size, err := reader.ReadByte()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
	strSize, err := strconv.ParseInt(string(size), 10, 64)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
	// Consume /r/n to get rid of 2 bytes '\r\n' that follows the number.
	if _, err := reader.ReadByte(); err != nil {
		fmt.Printf("Error: %v\n", err)
	}
	if _, err := reader.ReadByte(); err != nil {
		fmt.Printf("Error: %v\n", err)
	}
	name := make([]byte, strSize)
	if _, err := reader.Read(name); err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	fmt.Println(string(name))
}

func gracefulShutdown(funcs ...func() error) func() {
	return func() {
		for _, f := range funcs {
			if err := f(); err != nil {
				fmt.Println(err)
			}
		}
		os.Exit(1)
	}
}
