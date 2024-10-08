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

	// Read RESP input
	input := "$5\r\nMagda\r\n"

	name, err := ReadRESP(input)
	if err != nil {
		fmt.Printf("Error: %v", err)
	}

	fmt.Println(name)
}

var ErrEmptyRESP = errors.New("ReadRESP: empty input")

// ReadRESP first reads the string to determine number of characters we need to read.
// For example "Magda" is 5, plus additional 2 bytes.
//
// Then, consumes /r/n to get rid of 2 bytes '\r\n' that follows the number.
//
// Finally, returns the string.
func ReadRESP(input string) (string, error) {
	if input == "" {
		return "", ErrEmptyRESP
	}

	reader := bufio.NewReader(strings.NewReader(input))

	b, err := reader.ReadByte()
	if err != nil {
		return "", fmt.Errorf("ReadRESP: %w", err)
	}
	if b != '$' {
		return "", fmt.Errorf("ReadRESP: Invalid first byte type, expecting bulk strings only: %w", err)
	}

	size, err := reader.ReadByte()
	if err != nil {
		return "", fmt.Errorf("ReadRESP: %w", err)
	}
	strSize, err := strconv.ParseInt(string(size), 10, 64)
	if err != nil {
		return "", fmt.Errorf("ReadRESP: %w", err)
	}
	if _, err := reader.ReadByte(); err != nil {
		return "", fmt.Errorf("ReadRESP: %w", err)
	}
	if _, err := reader.ReadByte(); err != nil {
		return "", fmt.Errorf("ReadRESP: %w", err)
	}
	name := make([]byte, strSize)
	if _, err := reader.Read(name); err != nil {
		return "", fmt.Errorf("ReadRESP: %w", err)
	}

	return string(name), nil
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
