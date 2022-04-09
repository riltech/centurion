package dto

import "github.com/google/uuid"

// Describes a registration request
type RegisterRequest struct {
	Name string `json:"name"`
	Team string `json:"team"`
}

// Describes a registration response
type RegisterResponse struct {
	CenturionResponse
	ID uuid.UUID `json:"id"`
}
