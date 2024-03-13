package config

import (
	"flag"
)

const (
	DIR         = "dir"
	DB_FILENAME = "dbfilename"
)

var CONFIG = map[string]string{}

func Config() {
	var dir, dbFilename string

	flag.StringVar(&dir, "dir", "/home/daviddeadly/redis-go", "Provide a directory where RDB files will be stored")
	flag.StringVar(&dbFilename, "dbfilename", "redis-go.rdb", "Provide the name of the RDB file")

	flag.Parse()

	CONFIG[DIR] = dir
	CONFIG[DB_FILENAME] = dbFilename
}
