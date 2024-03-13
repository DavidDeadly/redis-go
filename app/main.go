package main

import (
	"github.com/codecrafters-io/redis-starter-go/app/server"
)

func main() {
	listener := server.StartServerOn("6379")

	defer listener.Close()

	server.ListenForConnections(listener)
}
