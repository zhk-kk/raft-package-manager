package cmd

import (
	"flag"
	"fmt"
)

type runCommand struct {
	fs   *flag.FlagSet
	name string
}

func newRunCommand() *runCommand {
	rc := &runCommand{
		fs: flag.NewFlagSet("run", flag.ContinueOnError),
	}

	// rc.fs.StringVar(&rc.name, "name", "World", "name of the person to be greeted")

	return rc
}

func (r *runCommand) Name() string {
	return r.fs.Name()
}

func (r *runCommand) Init(args []string) error {
	return r.fs.Parse(args)
}

func (r *runCommand) Run() error {
	fmt.Println("`run` was called.")
	return nil
}
