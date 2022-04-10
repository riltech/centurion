package challenge

// Describes a challenge
type Model struct {
	ID          string
	Name        string
	Description string
	Example     Example
}

// Example of a challenge
type Example struct {
	Hints    []interface{}
	Solution interface{}
}
