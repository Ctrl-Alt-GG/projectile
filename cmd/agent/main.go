package main

import (
	"fmt"
	"os"

	"github.com/Ctrl-Alt-GG/projectile/pkg/utils"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: agent [daemon|test|version] (-debug)")
		return
	}
	mode := os.Args[1]
	switch mode {
	case "version":
		fmt.Println(utils.GetLongVersion())
		fmt.Println("Built on", utils.GetBuildTimestamp())
	case "daemon":
		daemon()
	case "test":
		test()
	default:
		fmt.Println("Invalid mode. Valid modes: daemon, test")
	}
}
