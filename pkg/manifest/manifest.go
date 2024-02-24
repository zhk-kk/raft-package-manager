package manifest

import (
	"encoding/json"
	"errors"
	"fmt"

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

	// Retrieve the package type.
	pkgType, err := common.ParseMapField[string](m, "type")
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
	case "isPkg":
		result := IntegrationScriptsPkg{}
		if err := json.Unmarshal(raw, &result); err != nil {
			return nil, err
		}
		return result, nil
	default:
		return nil, fmt.Errorf("%w `%s`", ErrUnknownPkgType, pkgType)
	}
}

type PkgCommonInfo struct {
	PkgVersion    semver.Version
	RaftpmVersion semver.Version
}

type BinaryPkg struct {
	Name        string                    `json:"name"`
	Arch        map[string][]string       `json:"arch"`
	About       map[string]string         `json:"about"`
	BinRegistry map[string]common.PkgPath `json:"binRegistry"`
	BinShellExe map[string]string         `json:"binShellExe"`
}

type IntegrationScriptsPkg struct {
	TargetName          string         `json:"targetName"`
	TargetType          string         `json:"targetType"`
	DetectionScriptPath common.PkgPath `json:"detectionScript"`
	CapabilityScripts   []CapabilityScriptDesc
}

type CapabilityScriptDesc struct {
	Capability string         `json:"capability"`
	Path       common.PkgPath `json:"path"`
}
