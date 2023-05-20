package notifier

import (
	filesystem "github.com/fluxninja/aperture/v2/pkg/file-system"
	"github.com/fluxninja/aperture/v2/pkg/notifiers"
)

// KeyToFSNotifier holds the state of a notifier that writes raw/transformed contents of a watched file to another file.
type KeyToFSNotifier struct {
	fileInfo *filesystem.FileInfo
	notifiers.KeyBase
	path string
}

// Make sure KeyToFSNotifier implements KeyNotifier.
var _ notifiers.KeyNotifier = (*KeyToFSNotifier)(nil)

// NewKeyToFSNotifier returns a new notifier that writes raw/transformed contents to another file.
func NewKeyToFSNotifier(key notifiers.Key, dir string, ext string) *KeyToFSNotifier {
	fi := filesystem.NewFileInfo(dir, key.String(), ext)

	n := &KeyToFSNotifier{
		KeyBase:  notifiers.NewKeyBase(key),
		path:     fi.GetFilePath(),
		fileInfo: fi,
	}

	return n
}

// Start starts the key notifier.
func (n *KeyToFSNotifier) Start() error {
	// Should we remove the file or write nil to it?
	_ = n.fileInfo.RemoveFile()
	return nil
}

// Stop stops the key notifier.
func (n *KeyToFSNotifier) Stop() error {
	return nil
}

// Notify writes/removes to filesystem based on received event.
func (n *KeyToFSNotifier) Notify(event notifiers.Event) {
	// Should the new fi be created on every event?
	// Should it be the same directory as n.FileInfo.Directory?
	fi := filesystem.NewFileInfo(n.fileInfo.GetDirectory(), event.Key.String(), n.fileInfo.GetFileExt())

	switch event.Type {
	case notifiers.Write:
		_ = fi.WriteByteBufferToFile(event.Value)
	case notifiers.Remove:
		_ = fi.RemoveFile()
	}
}
