package main

import (
	"fmt"
	"os"

	"github.com/zhk-kk/raftpm/cmd"
	"github.com/zhk-kk/raftpm/global"
)

func main() {
	// Initialize global variables.
	if err := global.Init(); err != nil {
		fmt.Println("Unable to initialize the package manager:\n\t%w", err)
		os.Exit(1)
	}

	// Register the subcommands.
	resolver := cmd.NewResolver([]cmd.Subcommand{
		cmd.NewNested("develop", []cmd.Subcommand{
			cmd.NewPkgCompile(),
			cmd.NewWorkspaceInit(),
			cmd.NewSelfPackage(),
		}),
		cmd.NewDeploy(),
	})

	// Parse the arguments, running requested modules.
	if err := resolver.Parse(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
