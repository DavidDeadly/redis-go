package handlers

import (
	"strconv"
	"strings"

	"github.com/codecrafters-io/redis-starter-go/app/utils"
)

func SET (receivedParams []string) []byte {
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
}
