package pkg

import (
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
	"github.com/zhk-kk/raftpm/pkg/paths"
	"github.com/zhk-kk/raftpm/utils/archiver"
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

// GenerateSelfPackage makes a package from the currently running raftpm instance itself.
// func GenerateSelfPackage(w io.Writer) error {
// 	if err := global.Init(); err != nil {
// 		return err
// 	}

// 	files := make(map[string][]byte)

// }

// SelfPackage generates a package from currently running raftpm itself.
// func SelfPackage(w io.Writer, additionalCopyData map[string][]byte) error {
// 	if err := global.Init(); err != nil {
// 		return err
// 	}

// 	files := make(map[string][]byte)

// 	// Add all the additional copy data
// 	if len(additionalCopyData) != 0 {
// 		for p, contents := range additionalCopyData {
// 			files[paths.CopyDataDir+"/"+p] = contents
// 		}
// 	}

// 	// Add the raftpm executable itself
// 	raftpmFile, err := os.Open(global.Global.RunningExecutablePath())
// 	if err != nil {
// 		return err
// 	}
// 	files[paths.CopyDataDir+"/"+"raftpm"], err = io.ReadAll(raftpmFile)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

// CompileTemplate validates and compiles the template.
func CompileTemplate(templatePath string, w io.Writer) error {
	// Parse the manifest file.
	manifestFile, err := os.Open(path.Join(templatePath, paths.ManifestFile))
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
	case manifest.IntegrationScriptsPkg:
		if err := validateIntegrationScriptsPkgTemplate(templatePath, pkgManifest); err != nil {
			return fmt.Errorf("%w: %w", ErrCouldNotValidateTemplate, err)
		}
	default:
		return fmt.Errorf("[BUG]: CompileTemplate() got a package type that it couldn't process")
	}

	// Create the archiver.
	ar := archiver.NewArchiver(w)
	defer ar.Close()

	// Add a comment.
	ar.Comment("Package generated by the raft package manager")

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
			if err := ar.CreateDir(relativePath); err != nil {
				return err
			}
			return nil
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

		// Check if metadata.
		isMetadata := relativeRootDir == "metadata"

		// Special treatment for the manifest file.
		if isMetadata {
			isJson := filepath.Ext(relativePath) == "json"
			fileBuf, err = encodeMetadataFile(fileBuf, isJson)
			if err != nil {
				return err
			}
			zipFilePath = relativePath[:len(relativePath)-len(filepath.Ext(relativePath))]
		}

		// Add the file to the archive.
		builder := ar.FileBuilder(zipFilePath).Mode(fileInfo.Mode())

		// Add a metadata comment.
		if isMetadata {
			builder.Comment("RaftPM package metadata")
		}

		w, err := builder.Build()
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

	return nil
}

const (
	templateDirValidatorRequiredFile = iota
	templateDirValidatorRequiredDir  = iota
	templateDirValidatorMasked       = iota
)

type templateDirValidator struct {
	registry     map[string]int
	templatePath string
}

func newTemplateDirValidator(templatePath string) templateDirValidator {
	t := templateDirValidator{
		registry:     map[string]int{templatePath: templateDirValidatorRequiredDir},
		templatePath: templatePath,
	}
	return t
}

func (t templateDirValidator) Mask(path string) {
	t.registry[path] = templateDirValidatorMasked
}

func (t templateDirValidator) RequireDir(path string) {
	t.registry[path] = templateDirValidatorRequiredDir
}

func (t templateDirValidator) RequireFile(path string) {
	t.registry[path] = templateDirValidatorRequiredFile
}

func (t templateDirValidator) Validate() error {
	for p, mode := range t.registry {
		if mode == templateDirValidatorMasked {
			continue
		}

		stat, err := os.Stat(p)
		if err != nil {
			return err
		}

		switch mode {
		case templateDirValidatorRequiredFile:
			if stat.IsDir() {
				return fmt.Errorf("%w: `%s`", ErrRequiredFileIsDir, p)
			}
		case templateDirValidatorRequiredDir:
			if !stat.IsDir() {
				return fmt.Errorf("%w: `%s`", ErrRequiredDirIsFile, p)
			}
		}
	}
	return nil
}

// validateBinaryPkgTemplate validates the binary package template, or returns an error.
func validateBinaryPkgTemplate(templatePath string, binPkgManifest manifest.BinaryPkg) error {
	v := newTemplateDirValidator(templatePath)

	v.Mask(path.Join(templatePath, paths.IgnoreDir))
	v.RequireFile(path.Join(templatePath, paths.ManifestFile))
	v.RequireDir(path.Join(templatePath, paths.MetadataDir))
	v.RequireDir(path.Join(templatePath, paths.CopyDataDir))

	// Require all the BinRegistry paths.
	for _, p := range binPkgManifest.BinRegistry {
		if p.Type != common.PkgPathTypeLocal {
			continue
		}
		v.RequireFile(path.Join(templatePath, paths.CopyDataDir, p.Path))
	}

	// Verify that only registered binaries are referenced.
	for _, bin := range binPkgManifest.BinShellExe {
		if _, ok := binPkgManifest.BinRegistry[bin]; !ok {
			return fmt.Errorf("%w: `%s`", ErrUnregisteredBinaryReferenced, bin)
		}
	}

	return v.Validate()
}

// validateIntegrationScriptsPkgTemplate validates the
// integration scripts package template, or returns an error.
func validateIntegrationScriptsPkgTemplate(
	templatePath string, isPkgManifest manifest.IntegrationScriptsPkg,
) error {
	v := newTemplateDirValidator(templatePath)

	v.Mask(path.Join(templatePath, paths.IgnoreDir))
	v.RequireFile(path.Join(templatePath, paths.ManifestFile))
	v.RequireDir(path.Join(templatePath, paths.MetadataDir))
	v.RequireDir(path.Join(templatePath, paths.IntegrationScriptsDir))

	if isPkgManifest.DetectionScriptPath.Type == common.PkgPathTypeLocal {
		v.RequireFile(path.Join(
			templatePath,
			paths.IntegrationScriptsDir,
			isPkgManifest.DetectionScriptPath.Path,
		))
	}

	// Require all the capability scripts.
	for _, s := range isPkgManifest.CapabilityScripts {
		if s.Path.Type == common.PkgPathTypeLocal {
			v.RequireFile(path.Join(templatePath, paths.IntegrationScriptsDir, s.Path.Path))
		}
	}

	return v.Validate()
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
