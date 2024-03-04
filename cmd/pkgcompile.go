package cmd

import (
	"flag"
	"fmt"
	"os"
	"path"

	"github.com/zhk-kk/raftpm/pkg"
)

type pkgCompile struct {
	fs      *flag.FlagSet
	srcPath string
	outPath string
}

func NewPkgCompile() *pkgCompile {
	fs := flag.NewFlagSet("pkg-compile", flag.ContinueOnError)
	pc := pkgCompile{fs: fs}
	fs.StringVar(&pc.srcPath, "src", "", "path to the package source")
	fs.StringVar(&pc.outPath, "out", path.Join(".", "pkg.raftpm"), "output path of the new package")
	return &pc
}

func (pc *pkgCompile) Parse(args []string) error {
	pc.fs.Parse(args)

	if pc.srcPath == "" {
		return fmt.Errorf("pkg-compile: %w: `-src`", ErrArgumentMustBeSpecified)
	}

	// Create the resulting file.
	out, err := os.Create(pc.outPath)
	if err != nil {
		return err
	}
	defer out.Close()

	if err := pkg.CompileTemplate(pc.srcPath, out); err != nil {
		// Delete the newly-created file and return an error.
		os.Remove(pc.outPath)
		return fmt.Errorf("couldn't compile the template: %w", err)
	}

	return nil
}

func (*pkgCompile) Name() string { return "pkg-compile" }
