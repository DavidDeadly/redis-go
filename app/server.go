package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var Expiry = make(chan string)

const (
	OK        = "OK"
	PONG      = "PONG"
	NOT_FOUND = "$-1\r\n"
)

var Data = map[string]string{}

type CommandHandler func(params []string) []byte

var commandHandlers = map[string]CommandHandler{
	"PING": func(_ []string) []byte {
		return SimpleString(PONG)
	},

	"ECHO": func(params []string) []byte {
		message := strings.Join(params, " ")

		return SimpleString(message)
	},

	"SET": func(receivedParams []string) []byte {
		params := make([]string, 4)
		copy(params, receivedParams)

		key, value := params[0], params[1]

		arg1 := strings.ToUpper(params[2])
    param1 := params[3]

    if arg1 == "PX" {
      expireTime, _ := strconv.Atoi(param1)

      if expireTime != 0 {
        go Expire(key, expireTime, Expiry)
      }
    }

		Data[key] = value

		return SimpleString(OK)
	},

	"GET": func(params []string) []byte {
		key := params[0]

		value, ok := Data[key]

		if !ok {
			var msg *string
			return BulkString(msg)
		}

		return SimpleString(value)
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

  go ListenExpirations()

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

func ListenExpirations() {
  for key := range Expiry {
    fmt.Println("Deleted")
    delete(Data, key)
  }
}

func Expire(key string, timeMs int, channel chan string) {
  time.Sleep(time.Millisecond * time.Duration(timeMs))
  fmt.Printf("%d Milliseconds has passed %s expired...", timeMs, key)
  channel <- key
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

func SimpleError(err error) []byte {
	return []byte(fmt.Sprintf("-%s\r\n", err.Error()))
}

func BulkString(message *string) []byte {
	if message == nil {
		return []byte(NOT_FOUND)
	}

	return []byte(fmt.Sprintf("$%d\r\n%s\r\n", len(*message), *message))
}

func parseRedisProtocolRequest(request string) (string, []string) {
	places := strings.Split(request, "\r\n")
	places = places[:len(places)-1]

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
