package handlers

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/codecrafters-io/redis-starter-go/app/config"
	"github.com/codecrafters-io/redis-starter-go/app/utils"
)

var (
	Expiry = make(chan string)
	Data   = map[string]string{}
)

type CommandHandler func(params []string) []byte

var CommandHandlers = map[string]CommandHandler{
	"PING": func(_ []string) []byte {
		return utils.SimpleString(utils.PONG)
	},

	"ECHO": func(params []string) []byte {
		message := strings.Join(params, " ")

		return utils.SimpleString(message)
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

		return utils.SimpleString(utils.OK)
	},

	"GET": func(params []string) []byte {
		key := params[0]

		value, ok := Data[key]

		if !ok {
			return utils.BulkString(nil)
		}

		return utils.SimpleString(value)
	},

	"CONFIG": func(params []string) []byte {
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
	},
}

func Expire(key string, timeMs int, channel chan string) {
	time.Sleep(time.Millisecond * time.Duration(timeMs))
	fmt.Printf("%d Milliseconds has passed %s expired...", timeMs, key)
	channel <- key
}

func ListenExpirations() {
	for key := range Expiry {
		fmt.Println("Deleted")
		delete(Data, key)
	}
}
