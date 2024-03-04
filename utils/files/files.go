package files

import "io/fs"

func IsUnixExecutableFile(fileInfo fs.FileInfo) bool {
	return fileInfo.Mode()&0100 != 0 || fileInfo.Mode()&0010 != 0 || fileInfo.Mode()&0001 != 0
}
