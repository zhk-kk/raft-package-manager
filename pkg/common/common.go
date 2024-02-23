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
