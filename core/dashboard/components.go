package dashboard

import (
	"fmt"
	"sort"
	"time"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"github.com/riltech/centurion/core/combat"
	"github.com/riltech/centurion/core/player"
)

// Header component for the dashboard
func GetHeader(clock *widgets.Paragraph) []interface{} {
	welcome := widgets.NewParagraph()
	welcome.Text = "Welcome to Riltech's Centurion!"
	welcome.BorderStyle.Fg = ui.ColorYellow
	welcome.TextStyle.Modifier = ui.ModifierBold

	info := widgets.NewParagraph()
	info.Text = "https://github.com/riltech/centurion"
	info.Title = "More information"
	info.BorderStyle.Fg = ui.ColorYellow
	base := 1.0 / 10
	return []interface{}{
		ui.NewCol(base*4, welcome),
		ui.NewCol(base*4, info),
		ui.NewCol(base*2, clock),
	}
}

// Describes a component that can be refreshed
type IRefreshable interface {
	// Function to call when refresh is needed
	Refresh()
}

// LogWindow is a wrapper class over the lists
// to provide handy high level functionality
// for rendering logs
type LogWindow struct {
	CreatedAt time.Time
	List      *widgets.List
}

// Pushes a given number of item into the stack
func (lw *LogWindow) Push(item string) *LogWindow {
	if lw == nil || lw.List == nil {
		panic("LogWindow or underlying list is nil")
	}
	toAdd := fmt.Sprintf("%s %s", GetTimePassedSince(lw.CreatedAt, true), item)
	lw.List.Rows = append([]string{toAdd}, lw.List.Rows...)
	return lw
}

// Returns a new list for event logs
func GetEventLog(createdAt time.Time) *LogWindow {
	l := widgets.NewList()
	l.Title = "Event logs"
	l.Rows = []string{}
	l.TextStyle = ui.NewStyle(ui.ColorYellow)
	l.WrapText = false
	return &LogWindow{createdAt, l}
}

// Describes a clock window widget
type ClockWindow struct {
	createdAt time.Time
	widget    *widgets.Paragraph
}

// Interface check
var _ IRefreshable = (*ClockWindow)(nil)

// Constructor for ClockWindow
func NewClockWindow(createdAt time.Time) *ClockWindow {
	return &ClockWindow{
		createdAt: createdAt,
	}
}

// returns the widget
func (cw *ClockWindow) GetWidget() *widgets.Paragraph {
	if cw == nil {
		return nil
	}
	if cw.widget == nil {
		clock := widgets.NewParagraph()
		clock.Text = GetTimePassedSince(cw.createdAt, false)
		clock.BorderStyle.Fg = ui.ColorYellow
		clock.TextStyle.Modifier = ui.ModifierBold
		cw.widget = clock
	}
	return cw.widget
}

// Refreshes the time on the clock
func (cw *ClockWindow) Refresh() {
	cw.widget.Text = GetTimePassedSince(cw.createdAt, false)
}

// Describes a gauge component
type GaugeComponent struct {
	Percentage int
	Gauge      *widgets.Gauge
}

// Tracks overall uptime for the defensive team
type UptimeTrackerWindow struct {
	GaugeComponent
	service combat.IService
}

// Interface check
var _ IRefreshable = (*UptimeTrackerWindow)(nil)

func (utw *UptimeTrackerWindow) GetWidget() *widgets.Gauge {
	if utw == nil {
		return nil
	}
	if utw.Gauge == nil {
		utw.Gauge = widgets.NewGauge()
		utw.Gauge.Percent = utw.Percentage
		utw.Gauge.Title = "Defender uptime"
	}
	return utw.Gauge
}

func (utw *UptimeTrackerWindow) Refresh() {
	if utw == nil {
		return
	}
	utw.Gauge.Percent = utw.service.GetDefenseFailPercent()
}

// Constructor for an UptimeTrackerWindow
func NewUptimeTrackerWindow(combatService combat.IService) *UptimeTrackerWindow {
	return &UptimeTrackerWindow{
		GaugeComponent: GaugeComponent{100, nil},
		service:        combatService,
	}
}

// Tracks the overall success chance of the attacker team
type AttackerSuccessWindow struct {
	GaugeComponent
	combatService combat.IService
}

// Interface check
var _ IRefreshable = (*AttackerSuccessWindow)(nil)

func (asw *AttackerSuccessWindow) GetWidget() *widgets.Gauge {
	if asw == nil {
		return nil
	}
	if asw.Gauge == nil {
		asw.Gauge = widgets.NewGauge()
		asw.Gauge.Percent = asw.Percentage
		asw.Gauge.Title = "Attack success ratio"
	}
	return asw.Gauge
}

func (asw *AttackerSuccessWindow) Refresh() {
	if asw == nil {
		return
	}
	asw.Gauge.Percent = asw.combatService.GetAttackerSuccessPercent()
}

// Constructor for an UptimeTrackerWindow
func NewAttackerSuccessWindow(combatService combat.IService) *AttackerSuccessWindow {
	return &AttackerSuccessWindow{
		GaugeComponent: GaugeComponent{100, nil},
		combatService:  combatService,
	}
}

// Generic type for the best player tracking
type BestPlayersWindow struct {
	List *widgets.List
}

// Tracks the top 5 attackers
type BestAttackersWindow struct {
	BestPlayersWindow
	playerService player.IService
}

// Interface check
var _ IRefreshable = (*BestAttackersWindow)(nil)

func (baw *BestAttackersWindow) GetWidget() *widgets.List {
	if baw == nil {
		return nil
	}
	if baw.List == nil {
		baw.List = widgets.NewList()
		baw.List.Title = "Top attackers"
		baw.List.Rows = []string{"No attackers yet"}
		baw.List.TextStyle = ui.NewStyle(ui.ColorRed)
		baw.List.WrapText = false
		baw.List.SelectedRowStyle = baw.List.TextStyle
	}
	return baw.List
}

func (baw *BestAttackersWindow) Refresh() {
	if baw == nil {
		return
	}
	players := baw.playerService.GetTeam(player.TeamTypeAttacker)
	sort.Slice(players, func(i, j int) bool {
		return players[i].Score > players[j].Score
	})
	selectedNames := []string{}
	for _, p := range players {
		selectedNames = append(selectedNames, fmt.Sprintf("%s - %d", p.Name, p.Score))
	}
	if len(selectedNames) == 0 {
		selectedNames = []string{"No attackers yet"}
	}
	baw.List.Rows = selectedNames
}

// Constructor for a new best attackers window
func NewBestAttackersWindow(playerService player.IService) *BestAttackersWindow {
	return &BestAttackersWindow{playerService: playerService}
}

// Tracks the top 5 defenders
type BestDefendersWindow struct {
	BestPlayersWindow
	playerService player.IService
}

// Interface check
var _ IRefreshable = (*BestDefendersWindow)(nil)

func (bdw *BestDefendersWindow) GetWidget() *widgets.List {
	if bdw == nil {
		return nil
	}
	if bdw.List == nil {
		bdw.List = widgets.NewList()
		bdw.List.Title = "Top defenders"
		bdw.List.Rows = []string{"No defenders yet"}
		bdw.List.TextStyle = ui.NewStyle(ui.ColorBlue)
		bdw.List.WrapText = false
		bdw.List.SelectedRowStyle = bdw.List.TextStyle
	}
	return bdw.List
}

func (bdw *BestDefendersWindow) Refresh() {
	if bdw == nil {
		return
	}
	players := bdw.playerService.GetTeam(player.TeamTypeDefender)
	sort.Slice(players, func(i, j int) bool {
		return players[i].Score > players[j].Score
	})
	selectedNames := []string{}
	for _, p := range players {
		selectedNames = append(selectedNames, fmt.Sprintf("%s - %d", p.Name, p.Score))
	}
	if len(selectedNames) == 0 {
		selectedNames = []string{"No defenders yet"}
	}
	bdw.List.Rows = selectedNames
}

// Constructor for a new best attackers window
func NewBestDefendersWindow(playerService player.IService) *BestDefendersWindow {
	return &BestDefendersWindow{
		playerService: playerService,
	}
}
