package extractors

import (
	"errors"
	"fmt"
	"strings"
)

// Decided to roll our own impl, as:
// * we do not care about json pointer _logic_ (just parsing)
// * parsing is simple enough

// JSONPointer represents a parsed JSON pointer
//
// https://datatracker.ietf.org/doc/html/rfc6901
type JSONPointer struct {
	segments        []string
	escapedSegments []string // kept for fast String() call
}

// Segments returns the segments of the JSON pointer.
func (p JSONPointer) Segments() []string {
	return p.segments
}

// String returns the string representation of the JSON pointer escapedSegments.
func (p JSONPointer) String() string {
	return "/" + strings.Join(p.escapedSegments, "/")
}

// ParseJSONPointer parses a pointer into a JSONPointer.
func ParseJSONPointer(pointer string) (JSONPointer, error) {
	if pointer == "" {
		return JSONPointer{}, nil
	}

	if !strings.HasPrefix(pointer, "/") {
		return JSONPointer{}, fmt.Errorf("invalid start of json pointer: %q", pointer)
	}

	p := JSONPointer{}

	p.escapedSegments = strings.Split(pointer, "/")[1:]

	p.segments = make([]string, 0, len(p.escapedSegments))
	for _, escapedSegment := range p.escapedSegments {
		segment, err := unescapeJSONPointerSegment(escapedSegment)
		if err != nil {
			return JSONPointer{}, fmt.Errorf("%w %q", err, escapedSegment)
		}
		p.segments = append(p.segments, segment)
	}

	return p, nil
}

func unescapeJSONPointerSegment(s string) (string, error) {
	s = strings.ReplaceAll(s, "~1", "/")
	s = strings.ReplaceAll(s, "~0", "~")
	if strings.Contains(s, "~") {
		return "", errors.New("invalid ~-escape in json pointer")
	}
	return s, nil
}
