package core

import "github.com/riltech/centurion/core/bus"

// Describes a dashboard interface
type IDashboard interface {
	// Starts the dashboard process [This is a blocking call]
	Start()
	// Stops the dashboard process
	Stop()
}

// Dashboard implementation
type Dashboard struct {
	bus bus.IBus
}

// Interface check
var _ IDashboard = (*Dashboard)(nil)

func (d Dashboard) Start() {}
func (d Dashboard) Stop()  {}

// Constructor for dashboard
func NewDashboard(bus bus.IBus) IDashboard {
	return Dashboard{bus}
}
