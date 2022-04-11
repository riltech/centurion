package challenge

import (
	"fmt"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/google/uuid"
)

// First of the default modules
var defaultModuleReverseSorter = Model{
	ID:          uuid.NewString(),
	Type:        ChallengeTypeDefault,
	CreatorID:   "",
	Name:        "Reverse sorter",
	Description: "You receive a random length string array in the first parameter of the hints. Your aim is to change the order of the array and send it back as the first parameter of the solution array",
	Example: Example{
		Hints:     []interface{}{"123456"},
		Solutions: []interface{}{"654321"},
	},
}

// Returns the built in challenges
func getDefaultChallenges() []Model {
	return []Model{
		defaultModuleReverseSorter,
	}
}

// Returns hints for a default module
func getHintsForDefaultModule(m Model) ([]interface{}, error) {
	if m.Type != ChallengeTypeDefault {
		return nil, fmt.Errorf("Default module error: %s (%s) is %s", m.ID, m.Name, m.Type)
	}
	gofakeit.Seed(0)
	if m.Name == defaultModuleReverseSorter.Name {
		hint := fmt.Sprintf("%s%s%s", gofakeit.Word(), gofakeit.HipsterWord(), gofakeit.BuzzWord())
		return []interface{}{hint}, nil
	}
	return nil, fmt.Errorf("Default module not found: %s (%s) is %s", m.ID, m.Name, m.Type)
}

// Validates a given solution to a given default module
func isValidDefaultModuleSolution(m Model, hints []interface{}, solutions []interface{}) (bool, error) {
	if m.Type != ChallengeTypeDefault {
		return false, fmt.Errorf("Default module error: %s (%s) is %s", m.ID, m.Name, m.Type)
	}
	if m.Name == defaultModuleReverseSorter.Name {
		if len(solutions) == 0 {
			return false, fmt.Errorf("Solutions needs to be at least 1 long")
		}
		if len(hints) == 0 {
			return false, fmt.Errorf("Hints needs to be at least 1 long")
		}
		switch solutions[0].(type) {
		case string:
			solution := solutions[0].(string)
			switch hints[0].(type) {
			case string:
				hint := hints[0].(string)
				return isValidReverseSorterSolution(hint, solution), nil
			default:
				return false, fmt.Errorf("Hint has to be string")
			}
		default:
			return false, fmt.Errorf("Solution has to be string")
		}
	}
	return false, fmt.Errorf("Module is not a known module")
}

// Validates the solution directly
func isValidReverseSorterSolution(hint string, solution string) bool {
	if len(hint) != len(solution) {
		return false
	}
	for i := range hint {
		if hint[i] != solution[len(solution)-(1+i)] {
			return false
		}
	}
	return true
}
