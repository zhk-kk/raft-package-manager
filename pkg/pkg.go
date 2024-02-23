package pkg

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"path"

	"github.com/zhk-kk/raftpm/pkg/manifest"
	pkgpath "github.com/zhk-kk/raftpm/pkg/path"
)

// CompileTemplate validates and compiles the template.
func CompileTemplate(templatePath string) (*bufio.Reader, error) {
	manifestFile, err := os.Open(path.Join(templatePath, pkgpath.ManifestFile))
	if err != nil {
		return nil, err
	}
	defer manifestFile.Close()

	rawManifest, err := io.ReadAll(manifestFile)
	if err != nil {
		return nil, err
	}

	pkgCommonInfo := manifest.PkgCommonInfo{}
	pkgManifest, err := manifest.ParseManifest(rawManifest, &pkgCommonInfo)
	if err != nil {
		return nil, fmt.Errorf("couldn't parse manifest: %w", err)
	}
	fmt.Println(pkgCommonInfo, pkgManifest)

	return nil, errors.ErrUnsupported
}
