package cmd

import (
	"flag"
	"fmt"
	"os"

	"github.com/zhk-kk/raftpm/pkg"
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
		return fmt.Errorf("self-package: %w: `-out`", ErrArgumentMustBeSpecified)
	}

	// Create the resulting file.
	out, err := os.Create(sp.pathOut)
	if err != nil {
		return err
	}
	defer out.Close()

	if err := pkg.GenerateSelfPackage(out); err != nil {
		os.Remove(sp.pathOut)
		return err
	}

	return nil
}

func (*selfPackage) Name() string { return "self-package" }
