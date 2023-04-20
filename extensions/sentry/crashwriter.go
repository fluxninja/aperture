package sentry

import (
	"encoding/json"
	"sync"

	"github.com/eapache/queue"
)

const logCountLimit = 20

// crashWriter defines a crash writer with buffer to store the logs when the app crashes.
type crashWriter struct {
	crashLock sync.Mutex
	// buffer stores captured logs in ring-buffer queue
	buffer *queue.Queue
	// logCountLimit limits the number of lines of last logs to capture
	logCountLimit int
}

// newCrashWriter returns a new crash writer with new log buffer.
func newCrashWriter(limit int) *crashWriter {
	crashWriter := &crashWriter{
		buffer:        queue.New(),
		logCountLimit: limit,
	}
	return crashWriter
}

// Write writes the crash logs to the buffer in map[string]interface{} format for Sentry Breadcrumb data field.
func (w *crashWriter) Write(data []byte) (n int, err error) {
	w.crashLock.Lock()
	defer w.crashLock.Unlock()

	// Removes the element from the front of the queue when number of elements stored in the queue exceeds the logCountLimit.
	if w.buffer.Length() > w.logCountLimit-1 {
		_ = w.buffer.Remove()
	}

	var log map[string]interface{}
	err = json.Unmarshal(data, &log)
	if err != nil {
		return 0, err
	}

	// Puts data on the end of the queue buffer.
	w.buffer.Add(log)

	return len(data), nil
}

// Flush drains the crash writer buffer.
func (w *crashWriter) Flush() {
	w.crashLock.Lock()
	defer w.crashLock.Unlock()

	for w.buffer.Length() > 0 {
		_ = w.buffer.Remove()
	}
}

// GetCrashLogs returns the logs buffered in the crash writer until the crash.
func (w *crashWriter) GetCrashLogs() []map[string]interface{} {
	w.crashLock.Lock()
	defer w.crashLock.Unlock()

	logs := []map[string]interface{}{}
	for w.buffer.Length() > 0 {
		log := w.buffer.Remove().(map[string]interface{})
		logs = append(logs, log)
	}

	return logs
}
