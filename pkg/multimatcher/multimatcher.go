package multimatcher

import (
	"regexp"
	"sync"

	"github.com/fluxninja/aperture/pkg/log"
)

// Expr represents an expression that returns a bool decision based on labels.
//
// Note that even though this interface is public, it's recommended to use
// Exprs from this package, as multimatcher can optimize them better.
//
// (we're not doing it right now, but in future we might want to perform some
// optimizations like common expression elimination or running all regexes in
// one go (e.g. with engine like hyperscan)).
type Expr interface {
	Evaluate(Labels) bool
}

// All returns an expression that checks if all the given expressions are true.
func All(exprs []Expr) Expr {
	switch len(exprs) {
	case 0:
		return constNode(true)
	case 1:
		return exprs[0]
	default:
		return andNode{nodes: exprs}
	}
}

// Any returns an expression that checks if any of the given expressions are true.
func Any(exprs []Expr) Expr {
	switch len(exprs) {
	case 0:
		return constNode(false)
	case 1:
		return exprs[0]
	default:
		return orNode{nodes: exprs}
	}
}

// Not returns an expression that negates the given expression.
func Not(expr Expr) Expr { return notNode{node: expr} }

// LabelExists returns an expression that checks if the given label exists.
func LabelExists(label string) Expr { return existsNode{label: label} }

// LabelEquals returns an expression that checks if the given label is exactly equal to the value.
func LabelEquals(label string, value string) Expr {
	return exactMatchNode{label: label, value: value}
}

// LabelMatchesRegex compiles a regex and returns an expression that checks if the given label contains any match of the regex.
func LabelMatchesRegex(key string, pattern string) (Expr, error) {
	regex, err := regexp.Compile(pattern)
	if err != nil {
		log.Warn().Err(err).Msg("Failed to compile regex")
		return nil, err
	}
	return &regexMatchNode{
		label: key,
		regex: regex,
	}, nil
}

// Labels is a map from label to value.
type Labels map[string]string

// MultiMatcher is a database of items (entries) and their corresponding LabelSelector.
// Provides APIs to Add/Remove entries, Match against labels.
//
// Key is the type that identifies entries – used in adding / deleting, e.g. string.
// ResultCollection is a collection type that will be used in collecting
// results from Match, e.g. slice. Zero value of the ResultCollection should
// represent an empty one.
type MultiMatcher[Key comparable, ResultCollection any] struct {
	// Read/Write mutex to protect the maps and id
	mutex sync.RWMutex
	// Map from arbitrary key to entry
	entries map[Key]*matchEntry[ResultCollection]
}

// New returns a new multi matcher object.
func New[Key comparable, ResultCollection any]() *MultiMatcher[Key, ResultCollection] {
	return &MultiMatcher[Key, ResultCollection]{
		entries: make(map[Key]*matchEntry[ResultCollection]),
	}
}

// MatchCallback is called on entry in case of successful match. It should
// modify result collection and return changed one (e.g. by appending some
// object).
type MatchCallback[ResultCollection any] func(ResultCollection) ResultCollection

// Appender can be used to create a simple MatchCallback in case the
// ResultCollection is []T.
func Appender[T any](value T) MatchCallback[[]T] {
	return func(slice []T) []T { return append(slice, value) }
}

// Length returns number of entries in the MultiMatcher.
func (mm *MultiMatcher[_, _]) Length() int {
	mm.mutex.RLock()
	defer mm.mutex.RUnlock()
	return len(mm.entries)
}

// AddEntry adds an entry to the MultiMatcher.
//
// Key is an arbitrary key to map this entry – should be unique, can be used to remove the entry.
// MatchCallback is a callback used by Match in case of successful match.
// The entry is considered matching if matchExpr.Evaluate(labels) returns true.
func (mm *MultiMatcher[Key, ResultCollection]) AddEntry(key Key, matchExpr Expr, mc MatchCallback[ResultCollection]) error {
	mm.mutex.Lock()
	defer mm.mutex.Unlock()

	// Prepare entry
	entry := &matchEntry[ResultCollection]{
		mc:   mc,
		expr: matchExpr,
	}

	// Add entry
	_, ok := mm.entries[key]
	if ok {
		// Delete the previous entry
		_ = mm.removeEntryUnsafe(key)
	}
	mm.entries[key] = entry
	return nil
}

// RemoveEntry removes an entry in MultiMatcher at key.
func (mm *MultiMatcher[Key, _]) RemoveEntry(key Key) error {
	mm.mutex.Lock()
	defer mm.mutex.Unlock()

	return mm.removeEntryUnsafe(key)
}

func (mm *MultiMatcher[Key, _]) removeEntryUnsafe(key Key) error {
	delete(mm.entries, key)
	return nil
}

// Match runs matching against labels and returns a collection of results.
//
// Starting with empty ResultCollection, calls MatchCallback for every entry
// which matches given labels, and returns resulting ResultCollection.
func (mm *MultiMatcher[_, ResultCollection]) Match(labels Labels) ResultCollection {
	var resultCollection ResultCollection

	mm.mutex.RLock()
	defer mm.mutex.RUnlock()

	// Range through entries, invoke callback on match
	for _, entry := range mm.entries {
		matched := entry.expr.Evaluate(labels)
		if matched {
			resultCollection = entry.mc(resultCollection)
		}
	}

	return resultCollection
}

type matchEntry[ResultCollection any] struct {
	mc MatchCallback[ResultCollection]
	// Note: we're using Expr directly here, but nothing prevents us from
	// translating it into different representation
	expr Expr
}

type andNode struct {
	nodes []Expr
}

// Evaluate iterates through all andNodes and returns true if all of them return true.
func (an andNode) Evaluate(l Labels) bool {
	for _, node := range an.nodes {
		if !node.Evaluate(l) {
			return false
		}
	}
	return true
}

type orNode struct {
	nodes []Expr
}

// Evaluate iterates through all orNodes and returns true if any of them return true.
func (on orNode) Evaluate(l Labels) bool {
	for _, node := range on.nodes {
		if node.Evaluate(l) {
			return true
		}
	}
	return false
}

type notNode struct {
	node Expr
}

// Evaluate checks the negation of node expression.
func (nn notNode) Evaluate(l Labels) bool {
	return !nn.node.Evaluate(l)
}

type existsNode struct {
	label string
}

// Evaluate checks if the label exists in the labels.
func (en existsNode) Evaluate(l Labels) bool {
	_, ok := l[en.label]
	return ok
}

type exactMatchNode struct {
	label string
	value string
}

// Evaluate checks if the label exists in the labels and if so, checks whether the label is exactly equal to the value.
func (em exactMatchNode) Evaluate(l Labels) bool {
	value, ok := l[em.label]
	if ok {
		if value == em.value {
			return true
		}
	}
	return false
}

type regexMatchNode struct {
	regex *regexp.Regexp
	label string
}

// Evaluate checks if the label exists in the labels and if so, checks whether the label contains any match of the regex.
func (rm regexMatchNode) Evaluate(l Labels) bool {
	value, ok := l[rm.label]
	if ok {
		return rm.regex.MatchString(value)
	}
	return false
}

type constNode bool

// Evaluate returns boolean value of constNode.
func (cn constNode) Evaluate(l Labels) bool { return bool(cn) }
