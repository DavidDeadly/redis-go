package handlers

import "github.com/codecrafters-io/redis-starter-go/app/utils"

func GET(params []string) []byte {
	key := params[0]

	value, ok := Data[key]

	if !ok {
		return utils.BulkString(nil)
	}

	return utils.SimpleString(value)
}
