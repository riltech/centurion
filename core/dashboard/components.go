package dashboard

import (
	"fmt"
	"time"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
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
	l.SetRect(0, 0, 25, 8)
	return &LogWindow{createdAt, l}
}

// Describes a clock window widget
type ClockWindow struct {
	CreatedAt time.Time
	widget    *widgets.Paragraph
}

// returns the widget
func (cw *ClockWindow) GetWidget() *widgets.Paragraph {
	if cw == nil {
		return nil
	}
	if cw.widget == nil {
		clock := widgets.NewParagraph()
		clock.Text = GetTimePassedSince(cw.CreatedAt, false)
		clock.BorderStyle.Fg = ui.ColorYellow
		clock.TextStyle.Modifier = ui.ModifierBold
		cw.widget = clock
	}
	return cw.widget
}

// Refreshes the time on the clock
func (cw *ClockWindow) Refresh() {
	cw.widget.Text = GetTimePassedSince(cw.CreatedAt, false)
}
