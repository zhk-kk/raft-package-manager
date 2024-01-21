package cmd

import (
	"flag"
	"fmt"
)

type infoCommand struct {
	fs   *flag.FlagSet
	name string
}

func newInfoCommand() *infoCommand {
	ic := &infoCommand{
		fs: flag.NewFlagSet("info", flag.ContinueOnError),
	}

	// ic.fs.StringVar(&ic.name, "name", "World", "name of the person to be greeted")

	return ic
}

func (i *infoCommand) Name() string {
	return i.fs.Name()
}

func (i *infoCommand) Init(args []string) error {
	return i.fs.Parse(args)
}

func (i *infoCommand) Run() error {
	fmt.Println("`info` was called.")
	return nil
}
