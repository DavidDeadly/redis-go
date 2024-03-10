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

  _,  err = listener.Accept()
  exitOn(err, "Error accepting connection")
}

func exitOn(err error, msg string) {
  if err != nil {
    fmt.Printf("%s: %v", msg, err.Error())
    os.Exit(1)
  }
}
