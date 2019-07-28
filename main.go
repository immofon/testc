package main

import (
	"fmt"
	"os"
)

func main() {
	cmd := ""
	if len(os.Args) > 1 {
		cmd = os.Args[1]
	}
	if cmd == "" {
		cmd = os.Getenv("cmd")
	}
	switch cmd {
	case "server":
		server()
	case "client":
		client()
	default:
		fmt.Println("cmd: server|client")
	}
}
