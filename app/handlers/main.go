package handlers

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/codecrafters-io/redis-starter-go/app/utils"
)

var Expiry = make(chan string)
var Data = map[string]string{}

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
			var msg *string
			return utils.BulkString(msg)
		}

		return utils.SimpleString(value)
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
