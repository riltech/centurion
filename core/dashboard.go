package core

import (
	"fmt"
	"log"
	"time"

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
	createdAt time.Time
	bus       bus.IBus

	// channels

	playerRegisteredCh       <-chan *bus.BusEvent
	playerJoinedCh           <-chan *bus.BusEvent
	attackInitiatedCh        <-chan *bus.BusEvent
	attackFinishedCh         <-chan *bus.BusEvent
	defenseModuleInstalledCh <-chan *bus.BusEvent
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

	clockModule := dashboard.ClockWindow{
		CreatedAt: d.createdAt,
	}
	eventLog := dashboard.GetEventLog(d.createdAt)
	grid.Set(
		ui.NewRow(base,
			dashboard.GetHeader(clockModule.GetWidget())...,
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
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			clockModule.Refresh()
			ui.Render(grid)
			continue
		case value := <-d.playerRegisteredCh:
			event, err := value.DecodeRegistrationEvent()
			if err != nil {
				logrus.Error(err)
				continue
			}
			eventLog.Push(fmt.Sprintf("[Registration] %s registered to be a %s", event.Name, event.Team))
			ui.Render(grid)
			continue
		case value := <-d.playerJoinedCh:
			event, err := value.DecodePlayerJoinedEvent()
			if err != nil {
				logrus.Error(err)
				continue
			}
			eventLog.Push(fmt.Sprintf("[Join] %s joined %s team", event.Name, event.Team))
			ui.Render(grid)
			continue
		case value := <-d.attackFinishedCh:
			event, err := value.DecodeAttackFinishedEvent()
			if err != nil {
				logrus.Error(err)
				continue
			}
			result := "resolved"
			if !event.Success {
				result = "failed"
			}
			eventLog.Push(fmt.Sprintf("[Combat] %s %s %s challenge", event.AttackerName, result, event.ChallengeName))
			ui.Render(grid)
			continue
		case value := <-d.attackInitiatedCh:
			event, err := value.DecodeAttackInitiatedEvent()
			if err != nil {
				logrus.Error(err)
				continue
			}
			eventLog.Push(fmt.Sprintf("[Combat] %s initiated attack on %s challenge", event.AttackerName, event.ChallengeName))
			ui.Render(grid)
			continue
		case value := <-d.defenseModuleInstalledCh:
			event, err := value.DecodeDefenseModuleInstalledEvent()
			if err != nil {
				logrus.Error(err)
				continue
			}
			eventLog.Push(fmt.Sprintf("[Defense] %s installed new module '%s'", event.CreatorName, event.Name))
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
	attackFinishedCh := eventBus.Listen(bus.EventTypeAttackFinished)
	attackInitiatedCh := eventBus.Listen(bus.EventTypeAttackInitiated)
	defenseModuleInstalledCh := eventBus.Listen(bus.EventTypeDefenseModuleInstalled)
	return Dashboard{
		createdAt:                time.Now(),
		bus:                      eventBus,
		playerRegisteredCh:       playerRegisteredCh,
		playerJoinedCh:           playerJoinedCh,
		attackInitiatedCh:        attackInitiatedCh,
		attackFinishedCh:         attackFinishedCh,
		defenseModuleInstalledCh: defenseModuleInstalledCh,
	}
}
