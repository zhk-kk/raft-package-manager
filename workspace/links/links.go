package links

import (
	"os"

	"github.com/zhk-kk/raftpm/workspace/config"
)

type Links struct {
	path   string
	config *config.Config
}

func NewLinks(linksPath string, config *config.Config) *Links {
	l := Links{path: linksPath, config: config}
	return &l
}

func (l *Links) Init() error {
	if err := os.MkdirAll(l.path, os.ModePerm); err != nil {
		return err
	}

	return nil
}

func (l *Links) Load() error {
	return nil
}
