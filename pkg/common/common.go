package common

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

var (
	ErrExpectedFieldNotFound = errors.New("expected field not found")
	ErrWrongFieldType        = errors.New("wrong field type")
	ErrWrongPathFormat       = errors.New("wrong path format")
	ErrUnknownPathType       = errors.New("unknown path type")
)

func ParseMapField[T any](m map[string]interface{}, field string) (T, error) {
	var result T
	raw := m[field]
	if raw == nil {
		return result, fmt.Errorf("%w: `%s`", ErrExpectedFieldNotFound, field)
	}
	result, ok := raw.(T)
	if !ok {
		return result, fmt.Errorf("%w: field `%s`, expected type `%s`, got `%s`",
			ErrWrongFieldType, field, reflect.TypeOf(result).String(), reflect.TypeOf(raw).String())
	}
	return result, nil
}

const (
	PkgPathTypeLocal = iota
)

type PkgPath struct {
	Type int
	Path string
}

// Parses the provided string, producing a PkgPath struct.
func PkgPathFromString(path string) (PkgPath, error) {
	pkgPath := PkgPath{}
	split := strings.SplitN(path, ":", 2)
	if len(split) != 2 {
		return pkgPath, ErrWrongPathFormat
	}
	pkgPath.Path = split[1]
	switch split[0] {
	case "local":
		pkgPath.Type = PkgPathTypeLocal
	default:
		return pkgPath, fmt.Errorf("%w: got `%s`", ErrUnknownPathType, split[0])
	}
	return pkgPath, nil
}

// TypeString returns the string representation of a the path's type.
// If the type couldn't be recognized, an empty string is returned.
func (p PkgPath) TypeString() string {
	switch p.Type {
	case PkgPathTypeLocal:
		return "local"
	default:
		return ""
	}
}

// Returns true if the type is valid, otherwise returns false.
func (p PkgPath) ValidateType() bool { return p.TypeString() != "" }

// String returns the string representation of the path.
// If the type of a given path is invalid, it will not be added.
// In that case the returned string would look something like this:
// `:the/rest/of/the/path`. That way, if the path starts with a column (`:`),
// the type of the path is invalid.
func (p PkgPath) String() string { return p.TypeString() + ":" + p.Path }

// [TODO]: Properly parse the string.
// MarshalJSON implements Marshaler interface from the json package.
func (p PkgPath) MarshalJSON() ([]byte, error) {
	s := p.String()
	if s[0] == ':' {
		return nil, ErrUnknownPathType
	}
	return []byte("\"" + s + "\""), nil
}

// [TODO]: Properly parse the string.
// UnmarshalJSON implements Unmarshaler interface from the json package.
func (p *PkgPath) UnmarshalJSON(json []byte) error {
	path, err := PkgPathFromString(strings.Trim(string(json), "\""))
	p.Path = path.Path
	p.Type = path.Type
	return err
}
