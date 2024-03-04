package store

import (
	"os"
	"path"

	"github.com/zhk-kk/raftpm/workspace/config"
)

type Store struct {
	path   string
	config *config.Config
}

func (s Store) appsPath() string     { return path.Join(s.path, "apps") }
func (s Store) iscriptsPath() string { return path.Join(s.path, "iscripts") }

func NewStore(storePath string, config *config.Config) *Store {
	l := Store{path: storePath, config: config}
	return &l
}

func (s *Store) Init() error {
	if err := os.MkdirAll(s.appsPath(), os.ModePerm); err != nil {
		return err
	}

	if err := os.MkdirAll(s.iscriptsPath(), os.ModePerm); err != nil {
		return err
	}

	return nil
}

func (s *Store) Load() error {
	return nil
}
