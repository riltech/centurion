package bus

import "fmt"

// Enums
const EventTypeRegistration = "registration"

// Describes a message sent to the bus
type BusEvent struct {
	Type        string
	Information interface{}
}

// Decodes a registration event
func (be BusEvent) DecodeRegistration() (*RegistrationEvent, error) {
	if be.Type != EventTypeRegistration {
		return nil, fmt.Errorf("Event is not registration")
	}
	if conv, ok := be.Information.(RegistrationEvent); ok {
		return &conv, nil
	}
	return nil, fmt.Errorf("Event is not registration")
}

// Describes a registration event
type RegistrationEvent struct {
	Name string
	Team string
	ID   string
}
