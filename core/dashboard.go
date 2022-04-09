package core

import (
	"log"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"github.com/riltech/centurion/core/bus"
)

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

func (d Dashboard) Start() {
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	grid := ui.NewGrid()
	termWidth, termHeight := ui.TerminalDimensions()
	grid.SetRect(0, 0, termWidth, termHeight)

	p := widgets.NewParagraph()
	p.Text = "Welcome to Riltech's Centurion!"
	p.Title = "Game is on!"

	grid.Set(
		ui.NewRow(1,
			ui.NewCol(1, p),
		),
	)

	ui.Render(grid)

	for e := range ui.PollEvents() {
		if e.Type == ui.KeyboardEvent {
			break
		}
	}
}
func (d Dashboard) Stop() {}

// Constructor for dashboard
func NewDashboard(bus bus.IBus) IDashboard {
	return Dashboard{bus}
}
