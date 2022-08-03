package sentry

import (
	"io"
	"os"
	"sync"

	"github.com/eapache/queue"
	"github.com/fluxninja/lumberjack"
)

const logCountLimit = 100

var globalCrashWriter = getCrashWriter()

// GetCrashWriter returns a global crash writer.
func GetCrashWriter() *CrashWriter {
	return globalCrashWriter
}

func getCrashWriter() *CrashWriter {
	return NewCrashWriter(logCountLimit)
}

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

// Write writes the crash logs to the buffer and updates CrashWriter's buffer status.
func (w *CrashWriter) Write(data []byte) (n int, err error) {
	w.crashLock.Lock()
	defer w.crashLock.Unlock()

	// Puts data on the end of the queue buffer.
	w.buffer.Add(data)

	// Removes the element from the front of the queue when number of elements stored in the queue exceeds the logCountLimit.
	for {
		if w.buffer.Length() > w.logCountLimit {
			_ = w.buffer.Remove()
		} else {
			break
		}
	}

	return len(data), nil
}

// Flush writes last 100 lines of logs up until crash to the disk.
func (w *CrashWriter) Flush(lg io.Writer) {
	w.crashLock.Lock()
	defer w.crashLock.Unlock()

	for {
		if w.buffer.Length() > 0 {
			log := w.buffer.Remove()
			_, _ = lg.Write(log.([]byte))
		} else {
			break
		}
	}
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
