package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
)

var (
	workdir         = flag.String("workdir", "", "working directory")
	command         = flag.String("command", "", "command to run")
	tfTarget        = flag.String("tf-target", "", "TF target directory")
	additionalFlags = flag.String("additional-flags", "", "additional flags to pass to the command")
)

func main() {
	flag.Parse()

	ctx, cancelFn := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancelFn()

	args := []string{
		"run", ".", *command,
		"--tf-target", *tfTarget,
	}

	af := strings.Split(*additionalFlags, " ")
	args = append(args, af...)

	fmt.Println("executing go @", *workdir, "with", args)

	cmd := exec.CommandContext(ctx, "go", args...)
	cmd.Dir = *workdir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		fmt.Println("Error running command:", err)
		os.Exit(1)
	}
}
