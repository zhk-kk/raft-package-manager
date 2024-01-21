package main

import (
	"fmt"
	"os"

	"github.com/zhk-kk/raftpm/cmd"
)

func main() {
	if err := cmd.HandleRootArgs(os.Args[1:]); err != nil {
		fmt.Println("[ERROR]:", err)
		os.Exit(1)
	}
}
