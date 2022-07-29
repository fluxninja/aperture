package uuid

import "github.com/google/uuid"

// Provider is an interface used to generate UUIDs.
type Provider interface {
	New() string
}

// NewDefaultProvider returns instance of DefaultProvider.
func NewDefaultProvider() *DefaultProvider {
	return &DefaultProvider{}
}

// DefaultProvider is using github.com/google/uuid to provide UUIDs.
type DefaultProvider struct{}

// New returns new random UUID.
func (up *DefaultProvider) New() string { return uuid.New().String() }

// NewTestProvider returns instance of TestProvider.
func NewTestProvider(uuid string) *TestProvider {
	return &TestProvider{UUID: uuid}
}

// TestProvider always returns the same UUID.
type TestProvider struct {
	UUID string
}

// New returns the defined UUID.
func (up *TestProvider) New() string { return up.UUID }
