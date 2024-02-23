package manifest

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"

	"github.com/blang/semver/v4"
)

var (
	ErrExpectedFieldNotFound = errors.New("expected field not found")
	ErrWrongFieldType        = errors.New("wrong field type")
	ErrUnknownPkgType        = errors.New("unknown package type")
)

func ParseManifest(raw []byte, outPkgCommonInfo *PkgCommonInfo) (interface{}, error) {
	var m map[string]interface{}
	if err := json.Unmarshal(raw, &m); err != nil {
		return nil, err
	}

	// Retrieve the required raftpm version.
	raftpmVersionString, err := parseField[string](m, "raftpmVersion")
	if err != nil {
		return nil, err
	}

	raftpmVersion, err := semver.Make(raftpmVersionString)
	if err != nil {
		return nil, err
	}
	(*outPkgCommonInfo).RaftpmVersion = raftpmVersion

	// Retrieve the package version.
	pkgVersionString, err := parseField[string](m, "version")
	if err != nil {
		return nil, err
	}

	pkgVersion, err := semver.Make(pkgVersionString)
	if err != nil {
		return nil, err
	}
	(*outPkgCommonInfo).PkgVersion = pkgVersion

	// Retrieve the package name.
	pkgName, err := parseField[string](m, "name")
	if err != nil {
		return nil, err
	}
	(*outPkgCommonInfo).PkgName = pkgName

	// Retrieve the package type.
	pkgType, err := parseField[string](m, "type")
	if err != nil {
		return nil, err
	}

	// Parse the rest of the manifest according to the package type.
	switch pkgType {
	case "binPkg":
		result := BinaryPkg{}
		if err := json.Unmarshal(raw, &result); err != nil {
			return nil, err
		}
		return result, nil
	default:
		return nil, ErrUnknownPkgType
	}
}

func parseField[T any](m map[string]interface{}, field string) (T, error) {
	var result T
	raw := m[field]
	if raw == nil {
		return result, fmt.Errorf("%w: `%s`", ErrExpectedFieldNotFound, field)
	}
	result, ok := raw.(T)
	if !ok {
		return result, fmt.Errorf("%w: field `%s`, expected type `%s`", ErrWrongFieldType, field, reflect.TypeOf(result).String())
	}
	return result, nil
}

type PkgCommonInfo struct {
	PkgVersion    semver.Version
	RaftpmVersion semver.Version
	PkgName       string
}

type BinaryPkg struct {
	Arch        map[string][]string `json:"arch"`
	About       map[string]string   `json:"about"`
	BinRegistry map[string]string   `json:"binRegistry"`
	BinShellExe map[string]string   `json:"binShellExe"`
}

// {
//     "raftpmVersion": "0.0.0",
//     "type": "binPkg",
//     "arch": {
//         "cpu": ["x86_64"],
//         "os": ["linux"]
//     },
//     "about": {
//         "description": "A sample hello world program"
//     },
//     "binRegistry": {
//         "appExe": "path:helloworld"
//     },
//     "binShellExe": {
//         "helloworld": "appExe"
//     }
// }
