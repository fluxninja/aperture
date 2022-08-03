package services

import (
	"fmt"
)

// ServiceID uniquely identifies a service.
type ServiceID struct {
	AgentGroup string
	Service    string
}

func (s ServiceID) String() string {
	return fmt.Sprintf("%s.%s", s.AgentGroup, s.Service)
}
