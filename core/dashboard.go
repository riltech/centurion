package core

import (
	"fmt"
	"log"
	"time"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"github.com/riltech/centurion/core/bus"
	"github.com/riltech/centurion/core/combat"
	"github.com/riltech/centurion/core/dashboard"
	"github.com/riltech/centurion/core/player"
	"github.com/riltech/centurion/core/scoreboard"
	"github.com/sirupsen/logrus"
)

// Describes a dashboard interface
type IDashboard interface {
	// Starts the dashboard process [This is a blocking call]
	Start()
}

// Dashboard implementation
type Dashboard struct {
	createdAt     time.Time
	bus           bus.IBus
	scoreService  scoreboard.IService
	combatService combat.IService
	playerService player.IService

	// channels

	playerRegisteredCh       <-chan *bus.BusEvent
	playerJoinedCh           <-chan *bus.BusEvent
	attackInitiatedCh        <-chan *bus.BusEvent
	attackFinishedCh         <-chan *bus.BusEvent
	defenseModuleInstalledCh <-chan *bus.BusEvent
	defenseFailedCh          <-chan *bus.BusEvent
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

	clockWindow := dashboard.NewClockWindow(d.createdAt)
	eventLog := dashboard.GetEventLog(d.createdAt)
	uptimeWindow := dashboard.NewUptimeTrackerWindow(d.combatService)
	attackerSuccessWindow := dashboard.NewAttackerSuccessWindow(d.combatService)
	bestDefendersWindow := dashboard.NewBestDefendersWindow(d.playerService)
	bestAttackersWindow := dashboard.NewBestAttackersWindow(d.playerService)
	refresh := func() {
		bestDefendersWindow.Refresh()
		attackerSuccessWindow.Refresh()
		bestAttackersWindow.Refresh()
		bestDefendersWindow.Refresh()
	}
	grid.Set(
		ui.NewRow(base,
			dashboard.GetHeader(clockWindow.GetWidget())...,
		),
		ui.NewRow(base*4.0,
			ui.NewCol(0.5, bestDefendersWindow.GetWidget()),
			ui.NewCol(0.5, bestAttackersWindow.GetWidget()),
		),
		ui.NewRow(base*1.0,
			ui.NewCol(0.5, uptimeWindow.GetWidget()),
			ui.NewCol(0.5, attackerSuccessWindow.GetWidget()),
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
			clockWindow.Refresh()
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
			refresh()
			ui.Render(grid)
			continue
		case value := <-d.defenseFailedCh:
			event, err := value.DecodeDefenseFailedEvent()
			if err != nil {
				logrus.Error(err)
				continue
			}
			eventLog.Push(fmt.Sprintf("[Defense] %s failed a defense against %s", event.DefenderName, event.AttackerName))
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
			eventLog.Push(fmt.Sprintf("[Combat] %s %s '%s' challenge", event.AttackerName, result, event.ChallengeName))
			refresh()
			ui.Render(grid)
			continue
		case value := <-d.attackInitiatedCh:
			event, err := value.DecodeAttackInitiatedEvent()
			if err != nil {
				logrus.Error(err)
				continue
			}
			eventLog.Push(fmt.Sprintf("[Combat] %s initiated attack on '%s' challenge", event.AttackerName, event.ChallengeName))
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
func NewDashboard(
	eventBus bus.IBus,
	scoreService scoreboard.IService,
	combatService combat.IService,
	playerService player.IService,
) IDashboard {
	playerRegisteredCh := eventBus.Listen(bus.EventTypeRegistration)
	playerJoinedCh := eventBus.Listen(bus.EventTypePlayerJoined)
	attackFinishedCh := eventBus.Listen(bus.EventTypeAttackFinished)
	attackInitiatedCh := eventBus.Listen(bus.EventTypeAttackInitiated)
	defenseModuleInstalledCh := eventBus.Listen(bus.EventTypeDefenseModuleInstalled)
	defenseFailedCh := eventBus.Listen(bus.EventTypeDefenseFailed)
	return Dashboard{
		createdAt:                time.Now(),
		bus:                      eventBus,
		scoreService:             scoreService,
		combatService:            combatService,
		playerService:            playerService,
		playerRegisteredCh:       playerRegisteredCh,
		playerJoinedCh:           playerJoinedCh,
		attackInitiatedCh:        attackInitiatedCh,
		attackFinishedCh:         attackFinishedCh,
		defenseModuleInstalledCh: defenseModuleInstalledCh,
		defenseFailedCh:          defenseFailedCh,
	}
}
