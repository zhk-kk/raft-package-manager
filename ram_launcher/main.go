package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
)

func main() {
	if err := root(os.Args[1:]); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

type Runner interface {
	Init([]string) error
	Run() error
	Name() string
}

type InfoCommand struct {
	fs   *flag.FlagSet
	name string
}

func NewInfoCommand() *InfoCommand {
	ic := &InfoCommand{
		fs: flag.NewFlagSet("greet", flag.ContinueOnError),
	}

	ic.fs.StringVar(&ic.name, "name", "World", "name of the person to be greeted")

	return ic
}

func (i *InfoCommand) Name() string {
	return i.fs.Name()
}

func (i *InfoCommand) Init(args []string) error {
	return i.fs.Parse(args)
}

func (i *InfoCommand) Run() error {
	fmt.Println("Info was called.")
	return nil
}

func root(args []string) error {
	if len(args) < 1 {
		return errors.New("You must provide a sub-command when running `raftpm`")
	}

	cmds := []Runner{
		NewInfoCommand(),
	}

	subcommand := os.Args[1]

	for _, cmd := range cmds {
		if cmd.Name() == subcommand {
			cmd.Init(os.Args[2:])
			return cmd.Run()
		}
	}

	return fmt.Errorf("Unknown subcommand: %s", subcommand)
}
