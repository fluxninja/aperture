package panic

import (
	"io"
	"os"
	"sync"

	"github.com/FluxNinja/lumberjack"
)

const logCountLimit = 100

// CrashWriter defines a crash writer with buffer to store the logs when the app crashes.
type CrashWriter struct {
	// buffer stores captured log in circular manner
	buffer [][]byte
	// current tracks the current log index in Buffer
	current int
	// logCountLimit limits the number of lines of last logs to capture
	logCountLimit int
	crashLock     sync.Mutex
	// full tracks if the buffer has reached the limit
	full bool
}

// Write writes the crash logs to the buffer and updates CrashWriter's buffer status.
func (w *CrashWriter) Write(data []byte) (n int, err error) {
	w.crashLock.Lock()
	defer w.crashLock.Unlock()

	length := len(data)
	if len(w.buffer[w.current]) != length {
		w.buffer[w.current] = make([]byte, length)
	}
	copy(w.buffer[w.current], data)

	if w.current == w.logCountLimit-1 {
		w.full = true
		w.current = 0
	} else {
		w.current++
	}

	return length, nil
}

// Flush writes last 100 lines of logs up until crash to the disk.
func (w *CrashWriter) Flush(lg io.Writer) {
	w.crashLock.Lock()
	defer w.crashLock.Unlock()

	if w.full {
		for _, log := range w.buffer[w.current:] {
			_, _ = lg.Write(log)
		}
	}
	for _, log := range w.buffer[:w.current] {
		_, _ = lg.Write(log)
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

var globalCrashWriter = getCrashWriter()

// GetCrashWriter returns a global crash writer.
func GetCrashWriter() *CrashWriter {
	return globalCrashWriter
}

func getCrashWriter() *CrashWriter {
	return NewCrashWriter(logCountLimit)
}

// NewCrashWriter returns a new crash writer with new log buffer.
func NewCrashWriter(limit int) *CrashWriter {
	crashWriter := &CrashWriter{
		buffer:        make([][]byte, logCountLimit),
		current:       0,
		full:          false,
		logCountLimit: limit,
	}

	return crashWriter
}
