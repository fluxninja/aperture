package classifier

import (
	"github.com/open-policy-agent/opa/ast"
)

// Input is an input to Classify method.
//
// Classifier sometimes needs the input in ast.Value format and sometimes in
// interface{} format. Using this interface allows to minimize number of
// conversions.  Note that values returned from both methods should be
// equivalent.
//
// Each method can be called multiple times (including zero).
type Input interface {
	Value() ast.Value
	Interface() interface{}
}
