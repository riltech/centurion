package config

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

// Describes system config setup
type Specification struct {
	// Describes if the example clients are enabled
	ExampleEnabled bool `envconfig:"example_enabled"`
	// Describes the port number to use
	Port int `envconfing:"port" default:"8080"`
}

// Inits configuration
func Init() (*Specification, error) {
	var s Specification
	err := envconfig.Process("centurion", &s)
	if err != nil {
		return nil, err
	}
	fmt.Println(s.Port)
	return &s, nil
}
