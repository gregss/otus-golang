package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"sync"
	"time"
)

func main() {
	timeoutf := flag.String("timeout", "1s", "timeout in second")
	flag.Parse()
	timeout, _ := time.ParseDuration(*timeoutf)
	client := NewTelnetClient(net.JoinHostPort(flag.Arg(0), flag.Arg(1)), timeout, os.Stdin, os.Stdout)

	defer client.Close()
	if err := client.Connect(); err != nil {
		fmt.Println(err.Error())
		return
	}

	wg := &sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		client.Receive()
	}()
	go func() {
		defer wg.Done()
		client.Send()
	}()

	wg.Wait()
}
