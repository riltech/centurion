package config

import (
	"github.com/kelseyhightower/envconfig"
)

// Describes system config setup
type Specification struct {
	// Describes if the example clients are enabled
	ExampleEnabled bool `envconfig:"example_enabled"`
}

// Inits configuration
func Init() (*Specification, error) {
	var s Specification
	err := envconfig.Process("centurion", &s)
	if err != nil {
		return nil, err
	}
	return &s, nil
}
