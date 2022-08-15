package panichandler

import (
	"fmt"
	"path"
	"runtime"
	"strconv"
	"strings"
)

// Callstack is a full stacktrace.
type Callstack []uintptr

const stackLimit = 50

// Capture returns a full stacktrace.
func Capture() Callstack {
	callers := make([]uintptr, stackLimit)
	count := runtime.Callers(2, callers)
	stack := callers[:count]
	return Callstack(stack)
}

// Location holds the physical location of a stack entry.
type Location struct {
	// Directory is the directory the source file is from.
	Directory string
	// File is the filename of the source file.
	File string
	// Line is the line index in the file.
	Line int
}

// Function holds the logical location of a stack entry.
type Function struct {
	// Package is the go package the stack entry is from.
	Package string
	// Name is the function name the stack entry is from.
	Name string
}

// Entry holds the human understandable form of a StackTrace entry.
type Entry struct {
	// Location holds the physical location for this entry.
	Location Location
	// Location holds the logical location for this entry.
	Function Function
	// PC is the program counter for this entry.
	PC uintptr
}

// Entries returns all the entries for the stack trace.
func (c Callstack) Entries() []Entry {
	frames := runtime.CallersFrames([]uintptr(c))
	out := []Entry{}
	for {
		frame, more := frames.Next()
		dir, file := path.Split(frame.File)
		// fmt.Println("dir:", dir, "file:", file)
		fullname := frame.Function
		// fmt.Println("fullname: ",fullname)
		var pkg, name string
		if i := strings.LastIndex(fullname, "/"); i > 0 {
			i += strings.IndexRune(fullname[i+1:], '.')
			// we find the last /, then find the next . to split the function name from the package name
			pkg, name = fullname[:i+1], fullname[i+2:]
		} else {
			fullnameSplit := strings.Split(fullname, ".")
			pkg, name = fullnameSplit[0], fullnameSplit[1]
		}

		out = append(out, Entry{
			Location: Location{
				Directory: dir,
				File:      file,
				Line:      frame.Line,
			},
			Function: Function{
				Package: pkg,
				Name:    name,
			},
			PC: frame.PC,
		})
		if !more {
			break
		}
	}

	return out
}

// GetEntries returns stacktrace of Callstack in map[string]interface{} format.
func (c Callstack) GetEntries() map[string]interface{} {
	entries := c.Entries()
	lines := make(map[string]interface{})
	for i, e := range entries {
		lines["#"+strconv.Itoa(i+1)] = e.String()
	}

	return lines
}

// String returns stacktrace of Entry.
func (e Entry) String() string {
	return fmt.Sprint(e.Location, ":", e.Function)
}

// String returns Location of stack entry.
func (l Location) String() string {
	const strip = "fluxninja/"
	dir := l.Directory
	if i := strings.LastIndex(dir, strip); i > 0 {
		dir = dir[i+len(strip):]
	}
	return fmt.Sprint(dir, l.File, "@", l.Line)
}

// String returns name of Function.
func (f Function) String() string {
	return f.Name
}
