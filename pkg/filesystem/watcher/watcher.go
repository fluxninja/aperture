package watcher

import (
	"os"
	"path/filepath"
	"strings"
	"sync"

	"gopkg.in/fsnotify.v1"

	"github.com/fluxninja/aperture/v2/pkg/filesystem"
	"github.com/fluxninja/aperture/v2/pkg/log"
	"github.com/fluxninja/aperture/v2/pkg/notifiers"
	"github.com/fluxninja/aperture/v2/pkg/panichandler"
)

// watcher holds the state of the watcher.
// We create separate watcher at directory level as opposed to adding directories to single fsnotify instance.
// That way these named singletons exist for the lifetime of the fx app they are running in.
type watcher struct {
	waitGroup sync.WaitGroup
	fswatcher *fsnotify.Watcher
	notifiers.Trackers
	fileToSymlink map[string]string
	directory     string
	fileExt       string
}

// Make sure Watcher implements notifiers.Watcher interface.
var _ notifiers.Watcher = &watcher{}

// NewWatcher creates a new watcher instance that starts watching a directory via fsnotify.
func NewWatcher(directory, fileExt string) (*watcher, error) {
	fInfo, err := os.Stat(directory)
	if err != nil {
		log.Warn().Err(err).Str("directory", directory).Msg("Unable to stat directory")
		return nil, err
	} else if !fInfo.IsDir() {
		log.Error().Err(err).Str("directory", directory).Msg("Watcher being created on a non-directory")
		return nil, err
	}

	watcher := &watcher{
		directory:     filepath.Clean(directory),
		fileExt:       fileExt,
		fileToSymlink: make(map[string]string),
		Trackers:      notifiers.NewDefaultTrackers(),
	}

	return watcher, nil
}

// Start starts the watcher go routines and handles events from fsnotify.
func (w *watcher) Start() error {
	err := w.Trackers.Start()
	if err != nil {
		return err
	}
	// use fsnotify
	w.fswatcher, err = fsnotify.NewWatcher()
	if err != nil {
		log.Error().Err(err).Msg("Unable to create fsnotify watcher! Check system limits.")
		return err
	}

	// start tracking and accumulating events
	err = w.fswatcher.Add(w.directory)
	if err != nil {
		log.Warn().Err(err).Msg("Unable to add directory to fsnotify watcher!")
		// return err
	}

	w.waitGroup.Add(1)

	// bootstrap existing files -- we do not have to add notifiers here as they get added when watch
	// events go routine starts -- this code is for caching the file contents in trackers
	files, err := os.ReadDir(w.directory)
	if err != nil {
		log.Warn().Err(err).Str("directory", w.directory).Msg("Unable to list files of directory")
	}
	for _, file := range files {
		if !file.IsDir() {
			fileExt := filepath.Ext(file.Name())
			if fileExt == w.fileExt {
				filename := strings.TrimSuffix(file.Name(), fileExt)
				finfo := w.getFileInfo(filename)
				b, err := finfo.ReadAsByteBufferFromFile()
				if err != nil {
					log.Warn().Err(err).Str("file", w.getFileInfo(filename).String()).Msg("Unable to read file")
				}
				filePath := filepath.Clean(filepath.Join(w.directory, filename+w.fileExt))
				symLinkFilePath, _ := filepath.EvalSymlinks(filePath)
				w.fileToSymlink[filePath] = symLinkFilePath
				if b != nil {
					w.WriteEvent(notifiers.Key(filename), b)
				}
			}
		}
	}

	// watch events
	panichandler.Go(func() {
		defer w.waitGroup.Done()
	OUTER:
		for {
			select {
			case event, ok := <-w.fswatcher.Events:
				log.Trace().Interface("event", event).Bool("ok", ok).Str("dir", w.directory).Msg("got events")
				if !ok {
					break OUTER
				}
				op := event.Op
				filePath := filepath.Clean(event.Name)

				processEvent := func() {
					_, fileWithExt := filepath.Split(filePath)
					fileExt := filepath.Ext(fileWithExt)
					filename := strings.TrimSuffix(filepath.Base(fileWithExt), fileExt)
					symLinkFilePath, _ := filepath.EvalSymlinks(filePath)
					if symLinkFilePath != "" {
						symLinkFilePath = filepath.Clean(symLinkFilePath)
					}
					finfo := w.getFileInfo(filename)

					log.Trace().
						Str("event", event.String()).
						Str("filePath", filePath).
						Str("filename", filename).
						Str("fileExt", fileExt).
						Str("symLinkFilePath", symLinkFilePath).
						Msg("fsnotify")

					// only track specific extensions
					if fileExt == w.fileExt {
						// check whether file was modified or whether symlink changed
						if op&(fsnotify.Create|fsnotify.Write) != 0 ||
							(symLinkFilePath != "" && w.fileToSymlink[filePath] != symLinkFilePath) {
							w.fileToSymlink[filePath] = symLinkFilePath

							b, err := finfo.ReadAsByteBufferFromFile()
							if err != nil {
								log.Warn().Err(err).Str("file", finfo.String()).Msg("Unable to read file")
							}
							if b != nil {
								w.WriteEvent(notifiers.Key(filename), b)
							}
						} else if op&(fsnotify.Remove|fsnotify.Rename) != 0 {
							delete(w.fileToSymlink, filePath)
							w.RemoveEvent(notifiers.Key(filename))
						}
					}
				}
				// check whether filePath is a directory ending with "..data" - Kubernetes keeps configMap data in that directory
				if strings.HasSuffix(filePath, "..data") {
					// check whether filePath is a directory
					fstat, err := os.Stat(filePath)
					if err == nil {
						if fstat.IsDir() {
							log.Info().Str("filePath", filePath).Interface("event", event).Msg("fsnotify event on Kubernetes configmap data directory")
							// loop through w.fileToSymlink and for each key send a write event
							for f := range w.fileToSymlink {
								// rewrite event
								op = fsnotify.Write
								filePath = f
								processEvent()
							}
							continue
						}
					}
				}

				processEvent()

			case err, ok := <-w.fswatcher.Errors:
				log.Debug().Interface("err", err).Bool("ok", ok).Str("dir", w.directory).Msg("got errors")
				if !ok {
					break OUTER
				}
				log.Error().Err(err).Str("directory", w.directory).Msg("fsnotify error")
			}
		}
		log.Debug().Msg("exited fs watcher loop")
	})
	return nil
}

// Stop stops the watcher go routines.
func (w *watcher) Stop() error {
	w.fswatcher.Close()
	w.waitGroup.Wait()
	return w.Trackers.Stop()
}

func (w *watcher) getFileInfo(filename string) *filesystem.FileInfo {
	return filesystem.NewFileInfo(w.directory, filename, w.fileExt)
}
