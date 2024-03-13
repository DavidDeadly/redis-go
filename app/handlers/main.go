package handlers

import (
	"fmt"
	"time"
)

var (
	Expiry = make(chan string)
	Data   = map[string]string{}
)

type CommandHandler func(params []string) []byte

var CommandHandlers = map[string]CommandHandler{
	"PING": PING,

	"ECHO": ECHO,

	"GET": GET,

	"SET": SET,

	"CONFIG": CONFIG,
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
