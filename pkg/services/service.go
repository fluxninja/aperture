package services

// ServiceID uniquely identifies a service.
type ServiceID struct {
	Service string
}

func (s ServiceID) String() string {
	return s.Service
}
