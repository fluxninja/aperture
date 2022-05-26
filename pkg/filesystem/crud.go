package filesystem

import (
	"bytes"
	"errors"
	"os"
	"path"

	"github.com/natefinch/atomic"
	"google.golang.org/protobuf/proto"

	"aperture.tech/aperture/pkg/log"
)

// WriteByteBufferToFile atomically writes byte buffer to file.
func (fileInfo *FileInfo) WriteByteBufferToFile(data []byte) error {
	r := bytes.NewReader(data)
	file := fileInfo.GetFilePath()
	err := atomic.WriteFile(file, r)
	if err != nil {
		log.Error().Err(err).Str("file", file).Msg("Failed to write file on local filesystem")
		return err
	}
	return nil
}

// WriteMessageAsProtobufToFile serializes message with protobuf wire format and atomically writes to file.
func (fileInfo *FileInfo) WriteMessageAsProtobufToFile(msg proto.Message) error {
	dat, marshalErr := proto.Marshal(msg)
	if marshalErr != nil {
		log.Error().Err(marshalErr).Interface("message", msg).Msg("Failed to marshal message")
		return marshalErr
	}
	return fileInfo.WriteByteBufferToFile(dat)
}

// RemoveFile removes a file.
func (fileInfo *FileInfo) RemoveFile() error {
	file := fileInfo.GetFilePath()
	rmvErr := os.Remove(file)
	if rmvErr != nil {
		return rmvErr
	}
	return nil
}

// ReadAsByteBufferFromFile reads a file as a byte buffer.
func (fileInfo *FileInfo) ReadAsByteBufferFromFile() ([]byte, error) {
	file := fileInfo.GetFilePath()
	content, err := os.ReadFile(file)
	if err != nil {
		log.Error().Err(err).Str("file", file).Msg("Failed to read from file")
		return nil, err
	}
	return content, nil
}

// ReadAsStringFromFile reads a file as a string.
func (fileInfo *FileInfo) ReadAsStringFromFile() (string, error) {
	b, err := fileInfo.ReadAsByteBufferFromFile()
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// PurgeDirectory recursively removes all directories and files in a given directory.
func (fileInfo *FileInfo) PurgeDirectory() error {
	return PurgeDirectory(fileInfo.directory)
}

// PurgeDirectory recursively removes all directories and files in a given directory.
func PurgeDirectory(directory string) error {
	dirs, err := os.ReadDir(directory)
	if err != nil {
		log.Error().Err(err).Str("directory", directory).Msg("Failed to read directory on local filesystem")
		return err
	}
	for _, dir := range dirs {
		err := os.RemoveAll(path.Join([]string{directory, dir.Name()}...))
		if err != nil {
			log.Error().Err(err).Str("path", dir.Name()).Msg("Failed to remove file/directory on local filesystem")
			return err
		}
	}
	return nil
}

// ExistsFile checks if a file exists.
func (fileInfo *FileInfo) ExistsFile() (bool, error) {
	file := fileInfo.GetFilePath()
	_, err := os.Stat(file)
	if err != nil {
		switch {
		case errors.Is(err, os.ErrNotExist):
			// File does not exist
			return false, nil
		default:
			// File may or may not exist. See err for details.
			log.Warn().Err(err).Str("file", file).Msg("Unable to check if file exists")
			return false, err
		}
	}
	return true, nil
}
