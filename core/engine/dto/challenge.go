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
