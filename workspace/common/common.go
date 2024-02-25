package common

type WorkspaceElement interface {
	Init() error
	Load() error
}
