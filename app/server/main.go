package server

import (
	"fmt"
	"io"
	"net"
	"os"

	"github.com/codecrafters-io/redis-starter-go/app/handlers"
	"github.com/codecrafters-io/redis-starter-go/app/utils"
)


func StartServerOn(port string) net.Listener {
	listener, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", port))

	if err != nil {
		utils.PrintError(err, "Failed to bind to port 6379")
		os.Exit(1)
	}

  fmt.Printf("Server running on PORT: %s", port)

  return listener
}

func ListenForConnections(listener net.Listener) net.Conn {
	go handlers.ListenExpirations()

	for {
    connection, err := listener.Accept()
		if err != nil {
			utils.PrintError(err, "Error accepting connection")
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
		request := make([]byte, 1024)
		reqBytes, err := conn.Read(request)
		if err == io.EOF {
			break
		}
		if err != nil {
			utils.PrintError(err, "Error reading data from the conn")
			return
		}

		command, params := utils.ParseRedisProtocolRequest(string(request[:reqBytes]))

		fmt.Printf("Command: '%s', Params: '%s'\n", command, params)

		response := handlers.CommandHandlers[command](params)
		bytes, err := conn.Write(response)
		if err != nil {
			utils.PrintError(err, "Error sending data to the connection")
			return
		}

		fmt.Printf("Send %v bytes\n", bytes)
	}
}


