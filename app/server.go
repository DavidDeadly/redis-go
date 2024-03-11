package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type CommandHandler func (params []string) []byte

var commandHandlers = map[string]CommandHandler{

  "PING": func(_ []string) []byte {
		return SimpleString("PONG")
  },

  "ECHO": func(params []string) []byte {
    message := strings.Join(params, " ")

    return SimpleString(message)
  },
}

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
		request := make([]byte, 1024)
		reqBytes, err := conn.Read(request)
		if err == io.EOF {
			break
		}
		if err != nil {
			printError(err, "Error reading data from the conn")
			return
		}

    command, params := parseRedisProtocolRequest(string(request[:reqBytes]))

    fmt.Printf("Command: '%s', Params: '%s'\n", command, params)

		response := commandHandlers[command](params)
    bytes, err := conn.Write(response)

		if err != nil {
			printError(err, "Error sending data to the connection")
			return
		}

		fmt.Printf("Send %v bytes\n", bytes)
	}
}

func SimpleString(message string) []byte {
  return []byte(fmt.Sprintf("+%s\r\n", message))
}

func parseRedisProtocolRequest(request string) (string, []string) {
	places := strings.Fields(request)
	regex := regexp.MustCompile(`[\$\*]\d+`)
  numElements, err := strconv.Atoi(places[0][1:])

  if err != nil {
    fmt.Println("error parsing redis-cli message")

    return "DEFAULT", []string{}
  }

	parsedMessage := make([]string, 0, numElements)

	for _, data := range places {
		matches := regex.MatchString(data)

		if !matches {
      parsedMessage = append(parsedMessage, data)
		}
	}

  if len(parsedMessage) == 0 {
    fmt.Println("error getting the command and params from the request, none of them")
    return "DEFAULT", []string{}
  }

  command := strings.ToUpper(parsedMessage[0])

  return command, parsedMessage[1:]
}

func printError(err error, msg string) {
	fmt.Printf("%s: %v\n\n", msg, err.Error())
}
