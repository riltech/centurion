package dto

// Describes a registration request
type RegisterRequest struct {
	Name string `json:"name"`
	Team string `json:"team"`
}

// Describes a registration response
type RegisterResponse struct {
	CenturionResponse
	ID string `json:"id"`
}
