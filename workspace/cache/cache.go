package cache

import "os"

type Cache struct {
	path string
}

func NewCache(cachePath string) *Cache {
	l := Cache{path: cachePath}
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
