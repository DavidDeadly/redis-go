package main

import (
	"fmt"
	"io"
	"net"
	"os"
)

func main() {
	fmt.Println("Logs from your program will appear here!")

	listener, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		printError(err, "Failed to bind to port 6379")
		os.Exit(1)
	}

	defer listener.Close()

	for {
		var connection net.Conn
		connection, err = listener.Accept()
		if err != nil {
			printError(err, "Error accepting connection")
			os.Exit(1)
		}

		fmt.Println("\nNew connection")
		go handleConnection(connection)
	}
}

func handleConnection(conn net.Conn) {
  fmt.Println("Handling...")
	defer conn.Close()

	for {
		message := make([]byte, 1024)
		messageBytes, err := conn.Read(message)
		if err == io.EOF {
			break
		}
		if err != nil {
			printError(err, "Error reading data from the conn")
			return
		}

		fmt.Printf("Received '%s'\n", string(message[:messageBytes]))

		response := []byte("+PONG\r\n")
		bytes, err := conn.Write(response)
		if err != nil {
			printError(err, "Error sending data to the connection")
			return
		}

		fmt.Printf("Send %v bytes\n", bytes)
	}
}

func printError(err error, msg string) {
	fmt.Printf("%s: %v\n\n", msg, err.Error())
}
