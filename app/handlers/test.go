package handlers

import (
	"strings"

	"github.com/codecrafters-io/redis-starter-go/app/utils"
)

func PING(_ []string) []byte {
	return utils.SimpleString(utils.PONG)
}

func ECHO(params []string) []byte {
	message := strings.Join(params, " ")

	return utils.SimpleString(message)
}
