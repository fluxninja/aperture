package filesystem

import (
	"fmt"
	"path/filepath"
	"strings"
)

// FileInfo holds fields used for internal tracking of files, events in notifications and so on.
type FileInfo struct {
	directory string
	name      string
	ext       string
}

// NewFileInfo creates a new instance of FileInfo.
func NewFileInfo(directory, name, ext string) *FileInfo {
	fileInfo := &FileInfo{directory: filepath.Clean(directory), name: name, ext: ext}
	return fileInfo
}

// String returns a string representation of the FileInfo.
func (fileInfo *FileInfo) String() string {
	return fmt.Sprintf("FileInfo<"+
		"Name: %s "+
		"| Ext: %s "+
		"| Directory: %s"+
		">",
		fileInfo.name,
		fileInfo.ext,
		fileInfo.directory)
}

// ParseFilePath parses a filesystem path into a FileInfo.
func ParseFilePath(fpath string) *FileInfo {
	path := filepath.Clean(fpath)
	ext := filepath.Ext(path)
	name := strings.TrimSuffix(filepath.Base(path), ext)
	directory := filepath.Dir(path)
	return NewFileInfo(directory, name, ext)
}

// GetDirectory returns the directory from FileInfo.
func (fileInfo *FileInfo) GetDirectory() string {
	return fileInfo.directory
}

// GetFileName returns the file name from FileInfo.
func (fileInfo *FileInfo) GetFileName() string {
	return fileInfo.name
}

// GetFileExt returns the file extension from FileInfo.
func (fileInfo *FileInfo) GetFileExt() string {
	return fileInfo.ext
}

// GetFilePath returns the full file path from FileInfo, i.e., "directory/filename+extension".
func (fileInfo *FileInfo) GetFilePath() string {
	return filepath.Join(fileInfo.directory, fileInfo.name+fileInfo.ext)
}

// GetFilePathWithoutExt returns the full file path from FileInfo, i.e., "directory/filename".
func (fileInfo *FileInfo) GetFilePathWithoutExt() string {
	return filepath.Join(fileInfo.directory, fileInfo.name)
}
