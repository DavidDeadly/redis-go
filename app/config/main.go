package config

import (
	"flag"
	"fmt"
)

var DIR string
var DB_FILENAME string

func Config() {
	flag.StringVar(&DIR, "dir", "/home/daviddeadly/redis-go", "Provide a directory where RDB files will be stored")
	flag.StringVar(&DB_FILENAME, "dbfilename", "redis-go.rdb", "Provide the name of the RDB file")

	flag.Parse()

	fmt.Println("dir: ", DIR)
	fmt.Println("dbFilename: ", DB_FILENAME)
}
