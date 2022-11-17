package extractors

import (
	"fmt"
	"strings"
)

// AttributePath is a dot-separated path to attribute to extract data from
//
// Should be either:
// * one of the fields of [Attribute Context][ctx], or
// * a special "request.http.bearer" pseudo-attribute.
//
// Eg. Request.http.method
//
// [ctx]: https://www.envoyproxy.io/docs/envoy/latest/api-v3/service/auth/v3/attribute_context.proto
type AttributePath []string

// String returns the attribute path as a dot-separated string.
func (p AttributePath) String() string {
	return strings.Join(p, ".")
}

// ParseAttributePath parses a string into an attribute path.
func ParseAttributePath(path string) AttributePath {
	return strings.Split(path, ".")
}

// this function is used only internally on the `inputTemplate` document.
func (p AttributePath) traverse(document interface{}) (interface{}, error) {
	if len(p) == 0 {
		return document, nil
	}
	switch document := document.(type) {
	case map[string]interface{}:
		if subdoc, exists := document[p[0]]; exists {
			return AttributePath(p[1:]).traverse(subdoc)
		} else {
			return nil, fmt.Errorf("no such field %s", p[0])
		}
	default:
		return nil, fmt.Errorf("cannot descend into %T", document)
	}
}

func (p AttributePath) isBearer() bool {
	return len(p) == 3 && p[0] == "request" && p[1] == "http" && p[2] == "bearer"
}

// validate checks if attribute path points to a valid scalar field.
func (p AttributePath) validate() error {
	if p.isBearer() {
		return nil
	}

	// any header is ok
	if len(p) == 4 && p[0] == "request" && p[1] == "http" && p[2] == "headers" {
		return nil
	}

	// any filter metadata item (or sub-item) is ok
	if len(p) >= 3 && p[0] == "metadataContext" && p[1] == "filterMetadata" {
		return nil
	}

	value, err := p.traverse(inputTemplate)
	if err != nil {
		return fmt.Errorf("unknown attribute %v: %v", p, err)
	}
	if _, isMap := value.(map[string]interface{}); isMap {
		return fmt.Errorf("%v exists but points to an object; scalar value expected", p)
	}

	return nil
}

// validateAddress checks if attribute path points to an address object.
func (p AttributePath) validateAddress() error {
	if len(p) == 2 &&
		(p[0] == "source" || p[0] == "destination") &&
		p[1] == "address" {
		return nil
	}
	return fmt.Errorf("%v cannot be used with type address", p)
}
