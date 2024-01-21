package cmd

import (
	"errors"
	"fmt"
	"os"
)

type runner interface {
	Init([]string) error
	Run() error
	Name() string
}

func HandleRootArgs(args []string) error {
	if len(args) < 1 {
		return errors.New("you must provide a sub-command when running `raftpm`")
	}

	cmds := []runner{
		newInfoCommand(),
		newInstallCommand(),
		newRunCommand(),
		newVerifyPkgCommand(),
	}

	subcommand := os.Args[1]

	for _, cmd := range cmds {
		if cmd.Name() == subcommand {
			cmd.Init(os.Args[2:])
			return cmd.Run()
		}
	}

	return fmt.Errorf("unknown subcommand: %s", subcommand)
}
