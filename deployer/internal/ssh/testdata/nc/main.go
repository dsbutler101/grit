package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"os/signal"
	"syscall"
)

var (
	tcpSocket  = flag.String("tcp-socket", "", "TCP socket to connect to in form of <address>:<port>")
	unixSocket = flag.String("unix-socket", "", "Unix socket to connect to in form of </path/to/socket>")
)

func main() {
	flag.Parse()
	if *tcpSocket == "" && *unixSocket == "" {
		fmt.Println("Need at least one of tcp-socket, unix-socket")

		flag.Usage()
		os.Exit(1)
	}

	if *tcpSocket != "" && *unixSocket != "" {
		fmt.Println("Need at most one of tcp-socket, unix-socket")

		flag.Usage()
		os.Exit(1)
	}

	ctx, cancelFn := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancelFn()

	var conn net.Conn
	var err error

	if *tcpSocket != "" {
		conn, err = net.Dial("tcp", *tcpSocket)
	} else {
		conn, err = net.Dial("unix", *unixSocket)
	}

	if err != nil {
		panic(fmt.Sprintf("failed to dial target: %v", err))
	}

	go func() {
		_, err := io.Copy(os.Stdout, conn)
		if err != nil {
			fmt.Printf("failed to copy data to stdout: %v\n", err)
		}
	}()

	go func() {
		_, err := io.Copy(conn, os.Stdin)
		if err != nil {
			fmt.Printf("failed to copy data from stdin: %v\n", err)
		}
	}()

	<-ctx.Done()
	err = conn.Close()
	if err != nil {
		panic(fmt.Sprintf("failed to close connection: %v", err))
	}
}
