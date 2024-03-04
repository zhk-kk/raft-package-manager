package cmd

import (
	"flag"
	"fmt"

	"github.com/zhk-kk/raftpm/workspace"
)

type workspaceInit struct {
	fs       *flag.FlagSet
	path     string
	portable bool
}

func NewWorkspaceInit() *workspaceInit {
	fs := flag.NewFlagSet("workspace-init", flag.ContinueOnError)
	wi := workspaceInit{fs: fs}
	fs.StringVar(&wi.path, "path", "", "path to the workspace")
	fs.BoolVar(&wi.portable, "portable", false, "makes the workspace portable")
	return &wi
}

func (wi *workspaceInit) Parse(args []string) error {
	wi.fs.Parse(args)

	if wi.path == "" {
		return fmt.Errorf("workspace-init: %w: `-path`", ErrArgumentMustBeSpecified)
	}

	w := workspace.NewWorkspace(wi.path)

	if err := w.Init(); err != nil {
		// os.Remove(wi.path)
		return fmt.Errorf("couldn't initialize the workspace: %w", err)
	}

	err := w.Editor().
		Portable(wi.portable).
		ApplyChanges()

	if err != nil {
		return err
	}

	if err := w.Load(); err != nil {
		return err
	}

	return nil
}

func (workspaceInit) Name() string { return "workspace-init" }
