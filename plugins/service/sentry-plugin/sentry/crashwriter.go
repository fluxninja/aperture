package sentry

import (
	"encoding/json"
	"os"
	"sync"

	"github.com/eapache/queue"
	"github.com/fluxninja/lumberjack"
)

const logCountLimit = 20

// CrashWriter defines a crash writer with buffer to store the logs when the app crashes.
type CrashWriter struct {
	crashLock sync.Mutex
	// buffer stores captured logs in ring-buffer queue
	buffer *queue.Queue
	// logCountLimit limits the number of lines of last logs to capture
	logCountLimit int
}

// NewCrashWriter returns a new crash writer with new log buffer.
func NewCrashWriter(limit int) *CrashWriter {
	crashWriter := &CrashWriter{
		buffer:        queue.New(),
		logCountLimit: limit,
	}
	return crashWriter
}

// Write writes the crash logs to the buffer in map[string]interface{} format for Sentry Breadcrumb data field.
func (w *CrashWriter) Write(data []byte) (n int, err error) {
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
func (w *CrashWriter) Flush() {
	w.crashLock.Lock()
	defer w.crashLock.Unlock()

	for {
		if w.buffer.Length() > 0 {
			_ = w.buffer.Remove()
		} else {
			break
		}
	}
}

// GetCrashLogs returns the logs buffered in the crash writer until the crash.
func (w *CrashWriter) GetCrashLogs() []map[string]interface{} {
	w.crashLock.Lock()
	defer w.crashLock.Unlock()

	logs := []map[string]interface{}{}
	for w.buffer.Length() > 0 {
		log := w.buffer.Remove().(map[string]interface{})
		logs = append(logs, log)
	}

	return logs
}

// NewCrashFileWriter returns a lumberjack rolling logger which is used to write crash logs to the output file.
func NewCrashFileWriter(filename string) *lumberjack.Logger {
	writer := &lumberjack.Logger{
		Filename:   filename,
		MaxBackups: 10,
		MaxAge:     7,
	}
	return writer
}

// CloseCrashFileWriter closes the crash file writer.
func CloseCrashFileWriter(lg *lumberjack.Logger) {
	filename := lg.Filename
	_ = lg.Rotate()
	_ = lg.Close()
	_ = os.Remove(filename)
}
