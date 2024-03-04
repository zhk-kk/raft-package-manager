package workspace

import (
	"fmt"
	"io/fs"
	"os"
	"path"

	"github.com/zhk-kk/raftpm/workspace/cache"
	"github.com/zhk-kk/raftpm/workspace/common"
	"github.com/zhk-kk/raftpm/workspace/config"
	"github.com/zhk-kk/raftpm/workspace/links"
	"github.com/zhk-kk/raftpm/workspace/store"
)

type Workspace struct {
	path string

	config *config.Config

	cache *cache.Cache
	links *links.Links
	store *store.Store

	portable bool
}

func (w Workspace) storeDir() string        { return path.Join(w.path, "store") }
func (w Workspace) linksDir() string        { return path.Join(w.path, "links") }
func (w Workspace) cacheDir() string        { return path.Join(w.path, "cache") }
func (w Workspace) configDir() string       { return path.Join(w.path, "config") }
func (w Workspace) workConfigPath() string  { return path.Join(w.configDir(), "workspace") }
func (w Workspace) storeConfigPath() string { return path.Join(w.configDir(), "store") }
func (w Workspace) linksConfigPath() string { return path.Join(w.configDir(), "links") }
func (w Workspace) cacheConfigPath() string { return path.Join(w.configDir(), "cache") }

func (w Workspace) portableFlagFilePath() string { return path.Join(w.path, ".portable") }

func NewWorkspace(workspacePath string) *Workspace {
	w := Workspace{
		path:     workspacePath,
		portable: false,
	}

	// Create the config.
	w.config = config.NewConfig(w.workConfigPath())
	w.config.AddBool("isPortable", true, false)

	return &w
}

// Init initializes workspace, creating the workspace structure.
// Initializes all the contents of the workspace.
// Does nothing if the workspace is already initialized, apart from creating missing structures.
func (w *Workspace) Init() error {
	// Create the workspace directory.
	if err := os.MkdirAll(w.path, fs.ModePerm); err != nil {
		return fmt.Errorf("%s: %w", "unable to initialize the workspace directory", err)
	}

	// Create all the core elements.
	w.cache = cache.NewCache(w.cacheDir(), config.NewConfig(w.cacheConfigPath()))
	w.links = links.NewLinks(w.linksDir(), config.NewConfig(w.linksConfigPath()))
	w.store = store.NewStore(w.storeDir(), config.NewConfig(w.storeConfigPath()))

	// Initialize all the core workspace elements.
	workElemsInitOrder := []common.WorkspaceElement{
		w.cache,
		w.links,
		w.store,
	}

	for _, e := range workElemsInitOrder {
		if err := e.Init(); err != nil {
			return err
		}
	}

	return nil
}

// Load loads the workspace, along with all the core elements.
func (w *Workspace) Load() error {
	// Read the config.

	// Check if portable.
	// if _, err := os.Stat(w.portableFlagFilePath()); err == nil {
	// 	w.portable = true
	// } else if os.IsNotExist(err) {
	// 	w.portable = false
	// } else {
	// 	return err
	// }

	// Load all the core elements.
	workElemsLoadOrder := []common.WorkspaceElement{
		w.cache,
		w.links,
		w.store,
	}

	for _, e := range workElemsLoadOrder {
		if err := e.Load(); err != nil {
			return err
		}
	}

	return nil
}

func (w *Workspace) Portable() bool { return w.portable }

func (w *Workspace) Editor() *workspaceEditor {
	return &workspaceEditor{
		w: w,
	}
}

type workspaceEditor struct {
	w        *Workspace
	portable *bool
}

func (we *workspaceEditor) Portable(p bool) *workspaceEditor { we.portable = &p; return we }

func (we *workspaceEditor) ApplyChanges() error {
	if we.portable != nil {
		if *we.portable {
			if _, err := os.Create(we.w.portableFlagFilePath()); err != nil {
				return err
			}
		} else {
			if err := os.Remove(we.w.portableFlagFilePath()); err != nil && !os.IsNotExist(err) {
				return err
			}
		}
	}

	return nil
}
