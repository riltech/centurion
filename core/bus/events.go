package bus

import "fmt"

// Enums
const EventTypeRegistration = "registration"
const EventTypePlayerJoined = "player_joined"
const EventTypePanic = "panic"
const EventTypeAttackFinished = "attack_finished"
const EventTypeAttackInitiated = "attack_initiated"
const EventTypeDefenseModuleInstalled = "defense_module_installed"
const EventTypeDefenseFailed = "defense_failed"

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

// Decodes an attack finished event
func (be BusEvent) DecodeAttackFinishedEvent() (*AttackFinishedEvent, error) {
	if be.Type != EventTypeAttackFinished {
		return nil, fmt.Errorf("Event is not attack finished")
	}
	if conv, ok := be.Information.(AttackFinishedEvent); ok {
		return &conv, nil
	}
	return nil, fmt.Errorf("Event is not attack finished")
}

// Decodes an attack initiated event
func (be BusEvent) DecodeAttackInitiatedEvent() (*AttackInitiatedEvent, error) {
	if be.Type != EventTypeAttackInitiated {
		return nil, fmt.Errorf("Event is not attack initiated")
	}
	if conv, ok := be.Information.(AttackInitiatedEvent); ok {
		return &conv, nil
	}
	return nil, fmt.Errorf("Event is not attack initiated")
}

// Decodes a defense module installed event
func (be BusEvent) DecodeDefenseModuleInstalledEvent() (*DefenseModuleInstalledEvent, error) {
	if be.Type != EventTypeDefenseModuleInstalled {
		return nil, fmt.Errorf("Event is not defense module installed")
	}
	if conv, ok := be.Information.(DefenseModuleInstalledEvent); ok {
		return &conv, nil
	}
	return nil, fmt.Errorf("Event is not defense module installed")
}

// Decodes a defense failed event
func (be BusEvent) DecodeDefenseFailedEvent() (*DefenseFailedEvent, error) {
	if be.Type != EventTypeDefenseFailed {
		return nil, fmt.Errorf("Event is not defense failed")
	}
	if conv, ok := be.Information.(DefenseFailedEvent); ok {
		return &conv, nil
	}
	return nil, fmt.Errorf("Event is failed")
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
type AttackFinishedEvent struct {
	AttackerName  string
	ChallengeName string
	Success       bool
}

// Happens when a new attack is started
type AttackInitiatedEvent struct {
	AttackerName  string
	ChallengeName string
}

// Happens when a new defense module is added to the system
type DefenseModuleInstalledEvent struct {
	Name        string
	CreatorName string
}

// Happens when a defender fails to defend their module
type DefenseFailedEvent struct {
	DefenderName string
	AttackerName string
}
