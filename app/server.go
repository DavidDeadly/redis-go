package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	fmt.Println("Logs from your program will appear here!")

  listener, err := net.Listen("tcp", "0.0.0.0:6379")
  exitOn(err, "Failed to bind to port 6379")

  var connection net.Conn
  connection, err = listener.Accept()
  fmt.Println("HEEERE")

  message := make([]byte, 1024)
  messageBytes, err := connection.Read(message)
  exitOn(err, "Error reading data from the connection")
  fmt.Printf("Received '%s'\n", string(message[:messageBytes]))

  response := []byte("+PONG\r\n")
  bytes, err := connection.Write(response)
  exitOn(err, "Error sending data to the connection")
  fmt.Printf("Send %v bytes\n", bytes)

  exitOn(err, "Error accepting connection")
}

func exitOn(err error, msg string) {
  if err != nil {
    fmt.Printf("%s: %v", msg, err.Error())
    os.Exit(1)
  }
}
