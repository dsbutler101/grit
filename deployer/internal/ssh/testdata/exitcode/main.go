package main

import (
	"fmt"
	"os"
	"strconv"
)

func main() {
	if len(os.Args) < 2 {
		os.Exit(0)
	}

	errCode, err := strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Println("err:", err)
		os.Exit(-100)
	}

	os.Exit(errCode)
}
