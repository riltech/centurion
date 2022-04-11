package bus

import "fmt"

// Enums
const EventTypeRegistration = "registration"
const EventTypePlayerJoined = "player_joined"
const EventTypePanic = "panic"
const EventTypeAttackStateUpdate = "attack_state_update"

// Describes a message sent to the bus
type BusEvent struct {
	Type        string
	Information interface{}
}

// Decodes a registration event
func (be BusEvent) DecodeRegistrationEvent() (*RegistrationEvent, error) {
	if be.Type != EventTypeRegistration {
		return nil, fmt.Errorf("Event is not registration")
	}
	if conv, ok := be.Information.(RegistrationEvent); ok {
		return &conv, nil
	}
	return nil, fmt.Errorf("Event is not registration")
}

// Decodes a player joined event
func (be BusEvent) DecodePlayerJoinedEvent() (*PlayerJoinedEvent, error) {
	if be.Type != EventTypePlayerJoined {
		return nil, fmt.Errorf("Event is not player joined")
	}
	if conv, ok := be.Information.(PlayerJoinedEvent); ok {
		return &conv, nil
	}
	return nil, fmt.Errorf("Event is not player joined")
}

// Decodes an attack status update event
func (be BusEvent) DecodeAttackStateUpdateEvent() (*AttackStateUpdateEvent, error) {
	if be.Type != EventTypeAttackStateUpdate {
		return nil, fmt.Errorf("Event is not attack state update")
	}
	if conv, ok := be.Information.(AttackStateUpdateEvent); ok {
		return &conv, nil
	}
	return nil, fmt.Errorf("Event is not attack state update")
}

// Describes a registration event
type RegistrationEvent struct {
	Name string
	Team string
	ID   string
}

// Player joined the live game event
type PlayerJoinedEvent struct {
	Name string
	Team string
}

// Describes a success or a fail for an attack
type AttackStateUpdateEvent struct {
	AttackerName  string
	ChallengeName string
	Success       bool
}
