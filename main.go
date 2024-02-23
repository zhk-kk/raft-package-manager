package main

import (
	"fmt"
	"os"

	"github.com/zhk-kk/raftpm/cmd"
)

func main() {
	// Register the subcommands.
	resolver := cmd.NewResolver([]cmd.Subcommand{
		cmd.NewNested("develop", []cmd.Subcommand{cmd.NewPkgCompile()}),
	})

	// Parse the arguments, running requested modules.
	if err := resolver.Parse(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
