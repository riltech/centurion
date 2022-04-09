package core

// Describes a dashboard interface
type IDashboard interface {
	// Starts the dashboard process [This is a blocking call]
	Start()
	// Stops the dashboard process
	Stop()
}

// Dashboard implementation
type Dashboard struct {
	bus IBus
}

// Interface check
var _ IDashboard = (*Dashboard)(nil)

func (d Dashboard) Start() {}
func (d Dashboard) Stop()  {}

// Constructor for dashboard
func NewDashboard(bus IBus) IDashboard {
	return Dashboard{bus}
}
