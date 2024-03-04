package cmd

import (
	"flag"
	"fmt"
)

type selfPackage struct {
	fs      *flag.FlagSet
	pathOut string
}

func NewSelfPackage() *selfPackage {
	fs := flag.NewFlagSet("self-package", flag.ContinueOnError)
	sp := selfPackage{fs: fs}
	fs.StringVar(&sp.pathOut, "out", "", "path to the package")
	return &sp
}

func (sp *selfPackage) Parse(args []string) error {
	sp.fs.Parse(args)

	if sp.pathOut == "" {
		return fmt.Errorf("self-package: %w: `-path`", ErrArgumentMustBeSpecified)
	}

	return nil
}

func (*selfPackage) Name() string { return "self-package" }
