package utils

import "fmt"

func PrintError(err error, msg string) {
	fmt.Printf("%s: %v\n\n", msg, err.Error())
}

