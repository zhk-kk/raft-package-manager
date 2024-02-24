package global

import (
	"os"
	"path"
	"path/filepath"
)

var Global *global = nil

type global struct {
	runningExecutablePath string
	runningExecutableDir  string
	portableMode          bool
}

func (g global) PortableMode() bool            { return g.portableMode }
func (g global) RunningExecutablePath() string { return g.runningExecutablePath }
func (g global) RunningExecutableDir() string  { return g.runningExecutableDir }

func Init() error {
	if Global != nil {
		return nil
	}

	g := global{}

	exeSymlink, err := os.Executable()
	if err != nil {
		return err
	}
	g.runningExecutablePath, err = filepath.EvalSymlinks(exeSymlink)
	if err != nil {
		return err
	}
	g.runningExecutableDir = path.Dir(g.runningExecutablePath)

	_, err = os.Stat(path.Join(g.runningExecutableDir, ".portable"))
	if err != nil {
		g.portableMode = true
	} else if os.IsNotExist(err) {
		g.portableMode = false
	} else {
		return err
	}

	Global = &g

	return nil
}
