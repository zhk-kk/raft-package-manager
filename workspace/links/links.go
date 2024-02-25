package links

import "os"

type Links struct {
	path string
}

func NewLinks(linksPath string) *Links {
	l := Links{path: linksPath}
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
