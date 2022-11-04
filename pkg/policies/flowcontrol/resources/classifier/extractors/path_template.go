package extractors

import (
	"fmt"
	"strings"
)

// PathTemplate is an OpenAPI-inspired path template
//
// See github.com/fluxninja/aperture/api/gen/proto/go/aperture/classification/v1.PathTemplateMatcher.
type PathTemplate struct {
	Prefix              []string
	NumParams           int
	HasTrailingWildcard bool
}

// String returns string representation of the path template.
func (pt PathTemplate) String() string {
	tail := ""
	if pt.HasTrailingWildcard {
		tail = "/*"
	}
	_ = tail
	return "/" + strings.Join(pt.Prefix, "/") + strings.Repeat("/{}", pt.NumParams) + tail
}

// ParsePathTemplate parses a path template string into a PathTemplate.
func ParsePathTemplate(pathTemplate string) (PathTemplate, error) {
	pt := PathTemplate{}

	segments := strings.Split(pathTemplate, "/")
	for _, segment := range segments {
		if segment == "" {
			continue
		}
		if pt.HasTrailingWildcard {
			return PathTemplate{}, fmt.Errorf("* must be last segment in %q", pathTemplate)
		}

		if segment == "*" {
			pt.HasTrailingWildcard = true
		} else if strings.ContainsAny(segment, "{}") {
			if !strings.HasPrefix(segment, "{") || !strings.HasSuffix(segment, "}") {
				return PathTemplate{}, fmt.Errorf(
					"invalid parameter syntax in segment %q", segment,
				)
			}
			pt.NumParams += 1
		} else {
			if pt.NumParams != 0 {
				return PathTemplate{}, fmt.Errorf(
					"parametrized segments must come after stagic segments in %q",
					pathTemplate,
				)
			}
			pt.Prefix = append(pt.Prefix, segment)
		}
	}
	return pt, nil
}

// NumSegments returns the number of segments in the path template.
func (pt PathTemplate) NumSegments() int {
	return len(pt.Prefix) + pt.NumParams
}

// IsMoreSpecificThan returns whether path template might specific than other
// (eg. /foo is more specific than /*).
//
// Note: This function might have false positives (eg. /foo/bar is treated as
// more specific than /xyz, even though they're not related), but it shouldn't
// have any false negatives.
//
// This function can be used as sorting comparator.
func (pt PathTemplate) IsMoreSpecificThan(other PathTemplate) bool {
	if len(pt.Prefix) != len(other.Prefix) {
		return len(pt.Prefix) > len(other.Prefix)
	}
	if pt.NumParams != other.NumParams {
		return pt.NumParams > other.NumParams
	}
	return !pt.HasTrailingWildcard && other.HasTrailingWildcard
}
