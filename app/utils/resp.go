package utils

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

const (
	OK        = "OK"
	PONG      = "PONG"
	NOT_FOUND = "$-1\r\n"
)

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

func ParseRedisProtocolRequest(request string) (string, []string) {
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
