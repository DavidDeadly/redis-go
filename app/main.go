package main

import (
	"github.com/codecrafters-io/redis-starter-go/app/config"
	"github.com/codecrafters-io/redis-starter-go/app/server"
)

func main() {
	config.Config()

	listener := server.StartServerOn("6379")

	defer listener.Close()

	server.ListenForConnections(listener)
}
