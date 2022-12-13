package sim

import (
	"errors"
	"strings"

	rt "github.com/fluxninja/aperture/pkg/policies/controlplane/runtime"
	"github.com/lithammer/dedent"
)

// NewReadings creates a slice of readings from a slice of floats.
func NewReadings(values []float64) []rt.Reading {
	readings := make([]rt.Reading, 0, len(values))
	for _, value := range values {
		readings = append(readings, rt.NewReading(value))
	}
	return readings
}

// SanitizeYaml cleans up indentation of multiline string literal with yaml.
//
// This is needed as yaml parsing will otherwise fail on tabs.
func SanitizeYaml(source string) (string, error) {
	// Get rid of leadings tabs that are likely to happen within multiline string literals.
	source = dedent.Dedent(source)

	// Accept either space-indented yaml or tab-indented yaml. In the latter
	// case, tabs need to be converted to spaces before parsing.
	spaces := false
	tabs := false
	for _, line := range strings.Split(source, "\n") {
	loop:
		for _, c := range line {
			switch c {
			case ' ':
				spaces = true
			case '\t':
				tabs = true
			default:
				break loop
			}
		}
	}
	if spaces && tabs {
		return "", errors.New("mixed tabs and spaces in yaml")
	}
	if tabs {
		source = strings.ReplaceAll(source, "\t", "    ")
	}

	return source, nil
}
