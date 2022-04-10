package core

import (
	"fmt"
	"log"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"github.com/riltech/centurion/core/bus"
	"github.com/riltech/centurion/core/dashboard"
	"github.com/sirupsen/logrus"
)

// Describes a dashboard interface
type IDashboard interface {
	// Starts the dashboard process [This is a blocking call]
	Start()
}

// Dashboard implementation
type Dashboard struct {
	// dependencies
	bus bus.IBus

	// channels
	playerRegisteredCh <-chan *bus.BusEvent
	playerJoinedCh     <-chan *bus.BusEvent
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

	welcome := widgets.NewParagraph()
	welcome.Text = "Placeholder"
	base := 1.0 / 10

	eventLog := dashboard.GetEventLog()
	grid.Set(
		ui.NewRow(base,
			dashboard.GetHeader()...,
		),
		ui.NewRow(base*5.0,
			welcome,
		),
		ui.NewRow(base*4.0,
			eventLog.List,
		),
	)

	logrus.Info("Dashboard is rendering for the first time")
	ui.Render(grid)

	termUIEvents := ui.PollEvents()
	for {
		select {
		case value := <-d.playerRegisteredCh:
			event, err := value.DecodeRegistrationEvent()
			if err != nil {
				logrus.Error(err)
				continue
			}
			eventLog.Push(fmt.Sprintf("[Registration] %s registered to be a %s", event.Name, event.Team))
			ui.Render(grid)
			continue
		case e := <-termUIEvents:
			if e.Type == ui.KeyboardEvent {
				logrus.Infoln("Dashboard quits")
				return
			}
			continue
		}
	}
}

// Constructor for dashboard
func NewDashboard(eventBus bus.IBus) IDashboard {
	playerRegisteredCh := eventBus.Listen(bus.EventTypeRegistration)
	playerJoinedCh := eventBus.Listen(bus.EventTypePlayerJoined)
	return Dashboard{eventBus, playerRegisteredCh, playerJoinedCh}
}
