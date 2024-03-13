package handlers

import (
	"errors"
	"fmt"
	"strings"

	"github.com/codecrafters-io/redis-starter-go/app/config"
	"github.com/codecrafters-io/redis-starter-go/app/utils"
)

func CONFIG(params []string) []byte {
		var message string
		numParams := len(params)
		if numParams == 0 {
			message = "The CONFIG command needs a sub-action."
			return utils.SimpleError(errors.New(message))
		}

		action := strings.ToUpper(params[0])

		switch action {
		case "GET":
			if numParams <= 1 {
				message = "The CONFIG GET command needs a key to search for"
				return utils.SimpleError(errors.New(message))
			}

			key := strings.ToLower(params[1])

			value, ok := config.CONFIG[key]

			if ok {
				return utils.RespArray([]string{key, value})
			}
		default:
			message = fmt.Sprintf("The '%s' sub-action is not supported for the CONFIG command.", action)

			return utils.SimpleError(errors.New(message))
		}

		return utils.RespArray([]string{})
	}
