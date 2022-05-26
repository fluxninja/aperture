package notifier

import (
	"aperture.tech/aperture/pkg/filesystem"
	"aperture.tech/aperture/pkg/notifiers"
)

// PrefixToFSNotifier holds the state of a notifier that writes raw/transformed contents of a watched prefix to a directory.
type PrefixToFSNotifier struct {
	notifiers.PrefixNotifierBase
	path string
	ext  string
}

// Make	sure PrefixToFSNotifier implements PrefixNotifier.
var _ notifiers.PrefixNotifier = (*PrefixToFSNotifier)(nil)

// NewPrefixToFSNotifier returns a new prefix notifier that writes raw/transformed contents to a directory.
func NewPrefixToFSNotifier(path string, ext string) *PrefixToFSNotifier {
	n := &PrefixToFSNotifier{
		path: path,
		ext:  ext,
	}

	return n
}

// Start starts the prefix notifier.
func (n *PrefixToFSNotifier) Start() error {
	_ = filesystem.PurgeDirectory(n.path)
	return nil
}

// Stop stops the prefix notifier.
func (n *PrefixToFSNotifier) Stop() error {
	return nil
}

// GetKeyNotifier gets the underlying key notifier from prefix notifier.
func (n *PrefixToFSNotifier) GetKeyNotifier(key notifiers.Key) notifiers.KeyNotifier {
	fi := filesystem.NewFileInfo(n.path, key.String(), n.ext)
	kn := &KeyToFSNotifier{
		path:     fi.GetFilePath(),
		fileInfo: fi,
	}
	return kn
}
