package pkg

import (
	"archive/zip"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/zhk-kk/raftpm/pkg/common"
	"github.com/zhk-kk/raftpm/pkg/manifest"
	pkgpath "github.com/zhk-kk/raftpm/pkg/path"
)

var (
	ErrUnregisteredBinaryReferenced = errors.New("unregistered binary was referenced")
	ErrCouldNotValidateTemplate     = errors.New("couldn't validate the template")
	ErrRequiredFileIsDir            = errors.New("required file is a directory")
	ErrRequiredDirIsFile            = errors.New("required directory is a file")
)

var (
	AllowedArchCpu = []string{"x86_64", "x86", "aarch64", "aarch32"}
	AllowedArchOs  = []string{"bsd", "linux", "macos"}
)

// CompileTemplate validates and compiles the template.
func CompileTemplate(templatePath string, w io.Writer) error {
	// Parse the manifest file.
	manifestFile, err := os.Open(path.Join(templatePath, pkgpath.ManifestFile))
	if err != nil {
		return err
	}
	defer manifestFile.Close()

	rawManifest, err := io.ReadAll(manifestFile)
	if err != nil {
		return err
	}

	pkgCommonInfo := manifest.PkgCommonInfo{}
	pkgManifest, err := manifest.ParseManifest(rawManifest, &pkgCommonInfo)
	if err != nil {
		return fmt.Errorf("couldn't parse manifest: %w", err)
	}

	// Validate the template according to it's type.
	switch pkgManifest := pkgManifest.(type) {
	case manifest.BinaryPkg:
		if err := validateBinaryPkgTemplate(templatePath, pkgManifest); err != nil {
			return fmt.Errorf("%w: %w", ErrCouldNotValidateTemplate, err)
		}
	default:
		return fmt.Errorf("[BUG]: CompileTemplate() got a package type that it couldn't process")
	}

	// Create the archive.
	zipW := zip.NewWriter(w)
	defer zipW.Close()

	executableFiles := ""

	if err := filepath.WalkDir(templatePath, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Ignore the template directory.
		if p == templatePath {
			return nil
		}

		// [BUG]: Possible one? Could break if the template path is ".", but idk tbh.
		relativePath := strings.TrimPrefix(path.Clean(p), path.Clean(templatePath)+"/")
		relativeRootDir := strings.Split(relativePath, string(os.PathSeparator))[0]

		// Ignore everything in the `.ignore` directory.
		if relativeRootDir == ".ignore" {
			return nil
		}

		// Get the stats.
		fileInfo, err := os.Stat(p)
		if err != nil {
			return err
		}

		// All directories are made equal.
		if fileInfo.IsDir() {
			if _, err := zipW.Create(relativePath + "/"); err != nil {
				return err
			}
			return nil
		}

		// Check if the file is executable.
		if (fileInfo.Mode()&0100 != 0 || fileInfo.Mode()&0010 != 0 || fileInfo.Mode()&0001 != 0) && relativeRootDir != "metadata" {
			executableFiles += relativePath + "\n"
		}

		// Read the file.
		file, err := os.Open(p)
		if err != nil {
			return err
		}
		defer file.Close()

		fileBuf, err := io.ReadAll(file)
		if err != nil {
			return err
		}

		zipFilePath := relativePath // To be modified if metadata file relativeRootDir == "metadata"is encountered.

		// Special treatment for the manifest file.
		if relativeRootDir == "metadata" {
			isJson := filepath.Ext(relativePath) == "json"
			fileBuf, err = encodeMetadataFile(fileBuf, isJson)
			if err != nil {
				return err
			}
			zipFilePath = relativePath[:len(relativePath)-len(filepath.Ext(relativePath))]
		}

		// Add the file to the zip.
		w, err := zipW.Create(zipFilePath)
		if err != nil {
			return err
		}
		if _, err := w.Write(fileBuf); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return err
	}

	// Add the metadata/raftpmGen directory.
	if _, err := zipW.Create("metadata/raftpmGen/"); err != nil {
		return err
	}

	// Add the executableFiles list.
	if w, err := zipW.Create("metadata/raftpmGen/executableFiles"); err != nil {
		return err
	} else {
		executableFiles = strings.Trim(executableFiles, "\n")
		e, err := encodeMetadataFile([]byte(executableFiles), false)
		if err != nil {
			return err
		}
		if _, err := w.Write(e); err != nil {
			return err
		}
	}

	return nil
}

// validateBinaryPkgTemplate validates the binary package template, or returns an error.
func validateBinaryPkgTemplate(templatePath string, binPkgManifest manifest.BinaryPkg) error {
	const (
		pathRequiredFile = iota
		pathRequiredDir  = iota
		pathMasked       = iota
	)

	pathRegistry := map[string]int{
		templatePath: pathRequiredDir,
		path.Join(templatePath, pkgpath.MetadataDir):  pathRequiredDir,
		path.Join(templatePath, pkgpath.ManifestFile): pathRequiredFile,
		path.Join(templatePath, pkgpath.IgnoreDir):    pathMasked,
		path.Join(templatePath, pkgpath.CopyDataDir):  pathRequiredDir,
	}

	// Require all the BinRegistry paths.
	for _, p := range binPkgManifest.BinRegistry {
		if p.Type != common.PkgPathTypeLocal {
			continue
		}
		pathRegistry[path.Join(templatePath, pkgpath.CopyDataDir, p.Path)] = pathRequiredFile
	}

	// Verify that only registered binaries are referenced.
	for _, bin := range binPkgManifest.BinShellExe {
		if _, ok := binPkgManifest.BinRegistry[bin]; !ok {
			return fmt.Errorf("%w: `%s`", ErrUnregisteredBinaryReferenced, bin)
		}
	}

	// Go through the path registry.
	for p, mode := range pathRegistry {
		if mode == pathMasked {
			continue
		}

		stat, err := os.Stat(p)
		if err != nil {
			return err
		}

		switch mode {
		case pathRequiredFile:
			if stat.IsDir() {
				return fmt.Errorf("%w: `%s`", ErrRequiredFileIsDir, p)
			}
		case pathRequiredDir:
			if !stat.IsDir() {
				return fmt.Errorf("%w: `%s`", ErrRequiredDirIsFile, p)
			}
		}
	}

	return nil
}

// encodeMetadataFile strips the metadata file of all the unnecessary characters, and applies a base64 encoding.
func encodeMetadataFile(metadataFile []byte, isJson bool) ([]byte, error) {
	file := metadataFile

	if isJson {
		// Strip the file of all the unnecessary characters.
		compacted := bytes.NewBuffer([]byte{})
		if err := json.Compact(compacted, metadataFile); err != nil {
			return nil, err
		}
		file = compacted.Bytes()
	}

	// Encode the metadata file.
	buf := bytes.NewBuffer([]byte{})
	encoder := base64.NewEncoder(base64.StdEncoding, buf)
	if _, err := encoder.Write(file); err != nil {
		return nil, err
	}
	encoder.Close()
	return buf.Bytes(), nil
}

// decodeMetadataFile decodes the provided file, and returns a reader.
func decodeMetadataFile(r io.Reader) io.Reader { return base64.NewDecoder(base64.StdEncoding, r) }
