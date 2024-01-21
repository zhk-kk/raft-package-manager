package cmd

import (
	"flag"
	"fmt"
)

type verifyPkgCommand struct {
	fs   *flag.FlagSet
	name string
}

func newVerifyPkgCommand() *verifyPkgCommand {
	vpc := &verifyPkgCommand{
		fs: flag.NewFlagSet("info", flag.ContinueOnError),
	}

	// vpc.fs.StringVar(&vpc.name, "name", "World", "name of the person to be greeted")

	return vpc
}

func (v *verifyPkgCommand) Name() string {
	return v.fs.Name()
}

func (v *verifyPkgCommand) Init(args []string) error {
	return v.fs.Parse(args)
}

func (v *verifyPkgCommand) Run() error {
	fmt.Println("`verify-pkg` was called.")
	return nil
}
