package cmd

import (
	"errors"
	"flag"
	"fmt"
	"path"

	"github.com/zhk-kk/raftpm/pkg"
)

var (
	ErrPkgCompileNoSourcePathProvided = errors.New("pkg-compile: no source path provided")
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

	reader, err := pkg.CompileTemplate(pc.srcPath)
	if err != nil {
		return err
	}

	var p []byte
	if _, err := reader.Read(p); err != nil {
		return err
	}
	fmt.Println(p)

	return nil
}

func (pc *pkgCompile) Name() string { return "pkg-compile" }
