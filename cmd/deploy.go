package cmd

import (
	"flag"
	"fmt"

	"github.com/zhk-kk/raftpm/global"
	"github.com/zhk-kk/raftpm/workspace"
)

type deploy struct {
	fs              *flag.FlagSet
	destinationPath string
	portable        bool
}

func NewDeploy() *deploy {
	fs := flag.NewFlagSet("deploy", flag.ContinueOnError)
	d := deploy{fs: fs}
	fs.StringVar(&d.destinationPath, "dest", "", "path to the destination of deployment")
	fs.BoolVar(&d.portable, "portable", false, "makes the installation portable (intended to be used on removable drives)")
	return &d
}

func (d *deploy) Parse(args []string) error {
	d.fs.Parse(args)

	if d.destinationPath == "" {
		return fmt.Errorf("deploy: %w: `-dest`", ErrArgumentMustBeSpecified)
	}

	// [TODO]:
	//			1. (Done) Load the current workspace.
	//			2. (DONE) Initialize the workspace in the destination directory.
	//			3. Clone the store over to the destination workspace.
	//			4. Self-package to the destination's store.
	//			5. Run all the detection scripts, caching the result.
	//			6. Install all the required packages in the destination workspace.
	//			7. Install raftpm in the destination workspace.

	// Load the current workspace.
	curWork := workspace.NewWorkspace(global.Global.RunningExecutableDir())
	if err := curWork.Init(); err != nil {
		return err
	}

	// Initialize the workspace in the destination directory.
	destWork := workspace.NewWorkspace(d.destinationPath)
	if err := destWork.Init(); err != nil {
		return err
	}

	err := destWork.Editor().
		Portable(d.portable).
		ApplyChanges()
	if err != nil {
		return err
	}

	return nil
}

func (*deploy) Name() string { return "deploy" }
