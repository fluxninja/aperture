package extractors

import (
	"regexp"
	"sort"
	"strings"
)

// isValidPackageName checks package name according to strict rules.
func isValidPackageName(pkg string) bool {
	// This function is a little bit more strict than what rego allows. We
	// require path segments to be dot-separated list of valid identifiers while
	// rego also allows arbitrary strings, when `["xxx"]` notation is used.
	// See https://www.openpolicyagent.org/docs/latest/policy-language/#packages
	if pkg == "" {
		return false
	}
	for _, elem := range strings.Split(pkg, ".") {
		if !isRegoIdent(elem) {
			return false
		}
	}
	return true
}

// isRegoIdent checks if string is a valid rego identifier – not a rego
// keyword and composed of allowed characters - alphanumeric and underscore.
func isRegoIdent(ident string) bool {
	return ident != "" && !isRegoKeyword(ident) && regoIdentRegex.MatchString(ident)
}

var regoIdentRegex = regexp.MustCompile(`^[a-zA-Z_][0-9a-zA-Z_]*$`)

// isRegoKeyword checks if the string is one of rego keywords.
func isRegoKeyword(s string) bool {
	i := sort.SearchStrings(regoKeywords, s)
	return i < len(regoKeywords) && regoKeywords[i] == s
}

// https://www.openpolicyagent.org/docs/latest/policy-reference/#reserved-names
// rego keywords already sorted.
var regoKeywords = []string{
	"as",
	"default",
	"else",
	"false",
	"import",
	"not",
	"null",
	"package",
	"some",
	"true",
	"with",
}
