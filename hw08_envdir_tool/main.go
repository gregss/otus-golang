package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	args := os.Args[1:]
	env, err := ReadDir(args[0])
	if err != nil {
		log.Fatal(err)
	}

	returnCode := RunCmd(os.Args[2:], env)
	if returnCode > 0 {
		fmt.Printf("код ответа: %v", returnCode)
	}
}
