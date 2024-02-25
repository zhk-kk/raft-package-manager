package workspace

import (
	"io/fs"
	"os"
	"path"

	"github.com/zhk-kk/raftpm/workspace/cache"
	"github.com/zhk-kk/raftpm/workspace/common"
	"github.com/zhk-kk/raftpm/workspace/links"
	"github.com/zhk-kk/raftpm/workspace/store"
)

type Workspace struct {
	path  string
	cache *cache.Cache
	links *links.Links
	store *store.Store

	portable bool
}

func NewWorkspace(workspacePath string) *Workspace {
	w := Workspace{
		path:     workspacePath,
		portable: false,
	}
	return &w
}

func (w Workspace) storePath() string { return path.Join(w.path, "store") }
func (w Workspace) cachePath() string { return path.Join(w.path, "cache") }
func (w Workspace) linksPath() string { return path.Join(w.path, "links") }

func (w Workspace) portableFlagFilePath() string { return path.Join(w.path, ".portable") }

func (w *Workspace) coreElements() []common.WorkspaceElement {
	return []common.WorkspaceElement{
		w.cache,
		w.links,
		w.store,
	}
}

// Init initializes workspace, creating the workspace structure.
// Initializes all the contents of the workspace.
// Does nothing if the workspace is already initialized, apart from creating missing structures.
func (w *Workspace) Init() error {
	// Create the workspace directory.
	if err := os.MkdirAll(w.path, fs.ModePerm); err != nil {
		return err
	}

	// Create all the core elements.
	w.cache = cache.NewCache(w.cachePath())
	w.links = links.NewLinks(w.linksPath())
	w.store = store.NewStore(w.storePath())

	// Initialize all the core workspace elements.
	for _, e := range w.coreElements() {
		if err := e.Init(); err != nil {
			return err
		}
	}

	return nil
}

// Load loads the workspace, along with all the core elements.
func (w *Workspace) Load() error {
	// Check if portable.
	if _, err := os.Stat(w.portableFlagFilePath()); err == nil {
		w.portable = true
	} else if os.IsNotExist(err) {
		w.portable = false
	} else {
		return err
	}

	// Load all the core elements.
	for _, e := range w.coreElements() {
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
