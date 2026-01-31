package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: agent [daemon|test] (-debug)")
		return
	}
	mode := os.Args[1]
	switch mode {
	case "daemon":
		daemon()
	case "test":
		test()
	default:
		fmt.Println("Invalid mode. Valid modes: daemon, test")
	}
}
