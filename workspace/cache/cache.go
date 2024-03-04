package cache

import (
	"os"

	"github.com/zhk-kk/raftpm/workspace/config"
)

type Cache struct {
	path   string
	config *config.Config
}

func NewCache(cachePath string, config *config.Config) *Cache {
	l := Cache{path: cachePath, config: config}
	return &l
}

func (c *Cache) Init() error {
	if err := os.MkdirAll(c.path, os.ModePerm); err != nil {
		return err
	}

	return nil
}

func (c *Cache) Load() error {
	return nil
}
