package manifest

import (
	"encoding/json"
	"errors"

	"github.com/blang/semver/v4"
	"github.com/zhk-kk/raftpm/pkg/common"
)

var (
	ErrUnknownPkgType = errors.New("unknown package type")
)

func ParseManifest(raw []byte, outPkgCommonInfo *PkgCommonInfo) (interface{}, error) {
	var m map[string]interface{}
	if err := json.Unmarshal(raw, &m); err != nil {
		return nil, err
	}

	// Retrieve the required raftpm version.
	raftpmVersionString, err := common.ParseMapField[string](m, "raftpmVersion")
	if err != nil {
		return nil, err
	}

	raftpmVersion, err := semver.Make(raftpmVersionString)
	if err != nil {
		return nil, err
	}
	(*outPkgCommonInfo).RaftpmVersion = raftpmVersion

	// Retrieve the package version.
	pkgVersionString, err := common.ParseMapField[string](m, "version")
	if err != nil {
		return nil, err
	}

	pkgVersion, err := semver.Make(pkgVersionString)
	if err != nil {
		return nil, err
	}
	(*outPkgCommonInfo).PkgVersion = pkgVersion

	// Retrieve the package name.
	pkgName, err := common.ParseMapField[string](m, "name")
	if err != nil {
		return nil, err
	}
	(*outPkgCommonInfo).PkgName = pkgName

	// Retrieve the package type.
	pkgType, err := common.ParseMapField[string](m, "type")
	if err != nil {
		return nil, err
	}

	// Parse the rest of the manifest according to the package type.
	switch pkgType {
	case "binPkg":
		result := BinaryPkg{BinRegistry: make(map[string]common.PkgPath)}
		if err := json.Unmarshal(raw, &result); err != nil {
			return nil, err
		}

		// Retrieve the binary registry.
		binRegistry, err := common.ParseMapField[map[string]interface{}](m, "binRegistry")
		if err != nil {
			return nil, err
		}

		// Parse the binary registry.
		for entry, p := range binRegistry {
			pString, ok := p.(string)
			if !ok {
				return result, common.ErrWrongFieldType
			}
			pkgPath, err := common.PkgPathFromString(pString)
			if err != nil {
				return result, err
			}
			result.BinRegistry[entry] = pkgPath
		}
		return result, nil
	default:
		return nil, ErrUnknownPkgType
	}
}

type PkgCommonInfo struct {
	PkgVersion    semver.Version
	RaftpmVersion semver.Version
	PkgName       string
}

type BinaryPkg struct {
	Arch        map[string][]string       `json:"arch"`
	About       map[string]string         `json:"about"`
	BinRegistry map[string]common.PkgPath `json:"-"`
	BinShellExe map[string]string         `json:"binShellExe"`
}
