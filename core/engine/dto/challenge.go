package dto

// Describes a registration response
type FetchChallengesResponse struct {
	CenturionResponse
	Challenges []*ChallengeResponseDTO `json:"challenges"`
}

// Describes a challenge in fetch challenge response
type ChallengeResponseDTO struct {
	ID          string              `json:"id"`
	Name        string              `json:"name"`
	Description string              `json:"description"`
	Example     ChallengeExampleDTO `json:"example"`
}

// Describes an example in a challenge DTO
type ChallengeExampleDTO struct {
	Hints     []interface{} `json:"hints"`
	Solutions []interface{} `json:"solutions"`
}

// Describes a request body received in install challenge endpoint
type InstallChallengeRequest struct {
	// ID of the defender
	DefenderID string `json:"defenderId"`
	// Name of the challenge
	// NOTE: This has to be unique throughut the game
	Name string `json:"name"`
	// Description of the challenge
	Description string `json:"description"`
	// Example for the challenge
	Example ChallengeExampleDTO
}

// Success response of challenge installation endpoint
type InstallChallengeResponse struct {
	CenturionResponse
	// ID of the challenge you created
	// NOTE: you should persist this ID as you will be required
	// to provide hints and validations based on this ID
	ID string `json:"id"`
}
