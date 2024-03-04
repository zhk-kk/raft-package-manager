package archiver

import (
	"archive/zip"
	"io"
	"io/fs"
	"strings"
)

// archiver struct works a wrapper on top of zip.Writer. May use custom compression method.
type archiver struct {
	w *zip.Writer
}

// NewArchiver() returns a new archiver struct, which is a wrapper for zip.Writer.
func NewArchiver(w io.Writer) *archiver {
	a := archiver{
		w: zip.NewWriter(w),
	}
	return &a
}

// CreateDir() creates a directory using the provided path.
func (a *archiver) CreateDir(path string) error {
	_, err := a.w.Create(strings.TrimRight(path, "/") + "/")
	return err
}

// Comment() lets you add a comment to the resulting archive.
func (a *archiver) Comment(comment string) {
	a.w.SetComment(comment)
}

// Writer() returns a reference to the underlying writer. Usually this is **NOT** what you need.
func (a *archiver) Writer() *zip.Writer { return a.w }

// Close() closes the underlying zip.Writer
func (a *archiver) Close() error { return a.w.Close() }

type fileBuilder struct {
	a       *archiver
	path    string
	comment *string
	mode    *fs.FileMode
}

// FileBuilder() returns a fileBuilder struct, used to set all the options for the file.
// After that please call Build() to get the writer,
// to which the file contents should be written.
func (a *archiver) FileBuilder(path string) *fileBuilder {
	fb := fileBuilder{a: a, path: path}
	return &fb
}

// Build() finishes the building of the file, returning the writer,
// to which the file contents should be written.
func (fb *fileBuilder) Build() (io.Writer, error) {
	fh := zip.FileHeader{Name: fb.path}
	if fb.comment != nil {
		fh.Comment = *fb.comment
	}
	if fb.mode != nil {
		fh.SetMode(*fb.mode)
	}
	return fb.a.w.CreateHeader(&fh)
}

func (fb *fileBuilder) Comment(comment string) *fileBuilder {
	fb.comment = &comment
	return fb
}

func (fb *fileBuilder) Mode(mode fs.FileMode) *fileBuilder {
	fb.mode = &mode
	return fb
}
