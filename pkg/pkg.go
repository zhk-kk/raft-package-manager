package pkg

import (
	"errors"
	"strings"

	"github.com/blang/semver/v4"
)

func IsNameAllowed(name string) bool {
	for _, c := range name {
		if !(('A' <= c && c <= 'Z') || ('a' <= c && c <= 'z') || ('0' <= c && c <= '9') || c == '_' || c == '-') {
			return false
		}
	}
	return true
}

type PackageName struct {
	Name    string
	Version *semver.Version
}

// `pkgname` is written in the following format: `name:0.0.1-beta`. Version is optional.
func NewPackageName(pkgname string) (PackageName, error) {
	name, version, isContainsVersion := strings.Cut(pkgname, ":")

	new := PackageName{
		Name:    name,
		Version: nil,
	}

	if isContainsVersion {
		newVersion, err := semver.Make(version)
		if err != nil {
			return new, err
		}
		new.Version = &newVersion
	}

	if err := new.IsValid(); err != nil {
		return new, err
	}

	return new, nil
}

// Validate() returns an error if the provided PackageName is invalid.
func (p *PackageName) IsValid() error {
	if !IsNameAllowed(p.Name) {
		return errors.New("package name contains prohibited special characters")
	}

	return nil
}

func (p *PackageName) IsEqual(other *PackageName) bool {
	if p.Version == nil && other.Version == nil {
		return p.Name == other.Name
	} else if p.Version != nil && other.Version != nil {
		return p.Name == other.Name && p.Version.Equals(*other.Version)
	}
	return false
}
