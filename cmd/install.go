package cmd

import (
	"archive/tar"
	"errors"
	"flag"
	"fmt"
	"os"
	"path"

	"github.com/zhk-kk/raftpm/pkg"
)

type installCommand struct {
	fs   *flag.FlagSet
	name string
}

func newInstallCommand() *installCommand {
	ic := &installCommand{
		fs: flag.NewFlagSet("install", flag.ContinueOnError),
	}

	// ic.fs.StringVar(&ic.name, "name", "World", "name of the person to be greeted")

	return ic
}

func (i *installCommand) Name() string {
	return i.fs.Name()
}

func (i *installCommand) Init(args []string) error {
	return i.fs.Parse(args)
}

// [TODO]: Implement caching
func (i *installCommand) Run() error {
	fmt.Println("`install` was called.")

	pkgnames := make([]pkg.PackageName, 0)
	for _, pkgname := range i.fs.Args() {
		parsedPkgname, err := pkg.NewPackageName(pkgname)
		if err != nil {
			return fmt.Errorf("failed to parse the package name: %s", err)
		}
		pkgnames = append(pkgnames, parsedPkgname)
	}

	fmt.Println(pkgnames)

	// dirEntries, err := os.ReadDir("./raftpm-local-repo")
	// if err != nil {
	// 	return fmt.Errorf("unable to read the local repository")
	// }

	// pkgFileEntries := make([]fs.DirEntry, 0)
	// for _, entry := range dirEntries {
	// 	// Ignore directories and non-`.raftpm` files.
	// 	if entry.IsDir() || path.Ext(entry.Name()) != ".raftpm" {
	// 		continue
	// 	}

	// 	pkgFileEntries = append(pkgFileEntries, entry)

	// 	file, err :=
	// }

	// fmt.Println(pkgFileEntries)

	// [TODO]: Move everything related to working with the repositories to a separate package.

	localRepoPath := path.Join(".", "raftpm-local-repo")

	localRepoDir, err := os.Open(localRepoPath)
	if err != nil {
		return errors.New("failed to open the local repository")
	}
	defer localRepoDir.Close()
	// [TODO]: Create the directory if it doesn't exist

	localRepoContents, err := localRepoDir.ReadDir(0)
	if err != nil {
		return errors.New("failed to read the local repository")
	}

	localRepoPkgs := make([]string, 0)
	for _, entry := range localRepoContents {
		if entry.IsDir() || path.Ext(entry.Name()) != ".raftpm" {
			continue
		}
		localRepoPkgs = append(localRepoPkgs, entry.Name())
	}

	// Open each package.
	for _, pkgFileName := range localRepoPkgs {
		file, err := os.Open(path.Join(localRepoPath, pkgFileName))
		if err != nil {
			return fmt.Errorf("failed to read `%s` package file in the repository", pkgFileName)
		}
		defer file.Close()

		tarReader := tar.NewReader(nil)
		// [TODO]: finish the thing.
	}

	fmt.Println("ahstht", localRepoPkgs)

	return nil
}
