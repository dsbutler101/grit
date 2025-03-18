package main

import (
	"fmt"
	"os"
	"strconv"
)

func main() {
	if len(os.Args) > 2 && os.Args[1] == "fail" {
		ec, err := strconv.Atoi(os.Args[2])
		if err != nil {
			panic(fmt.Sprintf("fail to convert %q to int: %v", os.Args[2], err))
		}

		_, _ = fmt.Fprintf(os.Stderr, "Exiting with %d\n", ec)

		os.Exit(ec)
	}

	_, _ = fmt.Fprintln(os.Stdout, "out")
	_, _ = fmt.Fprintln(os.Stderr, "err")
}
