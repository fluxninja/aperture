package runtime

import "strings"

const (
	// NestedComponentDelimiter is the delimiter used to separate the parent circuit ID and the nested circuit ID.
	NestedComponentDelimiter = "."
	// RootComponentID is the ID of the root component.
	RootComponentID = "root"
)

// ComponentID is a unique identifier for a component.
type ComponentID interface {
	String() string
	ChildID(id string) ComponentID
	ParentID() (ComponentID, bool)
}

type componentID struct {
	id string
}

// componentID implements ComponentID.
var _ ComponentID = (*componentID)(nil)

// NewComponentID creates a new ComponentID.
func NewComponentID(id string) ComponentID {
	return &componentID{
		id: id,
	}
}

// ID returns the ID of the component.
func (cID componentID) String() string {
	return cID.id
}

// ChildID returns a new ComponentID that is a child of the current ComponentID.
func (cID componentID) ChildID(id string) ComponentID {
	return NewComponentID(cID.id + NestedComponentDelimiter + id)
}

// ParentID returns the parent ComponentID of the current ComponentID.
func (cID componentID) ParentID() (ComponentID, bool) {
	// Parent of root is an empty string.
	if cID.id == RootComponentID {
		return NewComponentID(""), true
	}
	// Parent Child component IDs are delimited by dot. So, we split the child component ID by dot and return the first part.
	// For example, if the child component ID is "root.1.2", then the parent component ID is "root.1".
	// Find the last delimiter in the ID
	delimiterIndex := strings.LastIndex(cID.id, NestedComponentDelimiter)
	if delimiterIndex == -1 {
		return nil, false
	}

	return NewComponentID(cID.id[:delimiterIndex]), true
}
