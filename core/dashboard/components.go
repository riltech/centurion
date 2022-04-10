package dashboard

import (
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

// Header component for the dashboard
func GetHeader() []interface{} {
	welcome := widgets.NewParagraph()
	welcome.Text = "Welcome to Riltech's Centurion!"
	welcome.BorderStyle.Fg = ui.ColorYellow
	welcome.TextStyle.Modifier = ui.ModifierBold

	info := widgets.NewParagraph()
	info.Text = "https://github.com/riltech/centurion"
	info.Title = "More information"
	info.BorderStyle.Fg = ui.ColorYellow
	return []interface{}{
		ui.NewCol(0.5, welcome),
		ui.NewCol(0.5, info),
	}
}

// LogWindow is a wrapper class over the lists
// to provide handy high level functionality
// for rendering logs
type LogWindow struct {
	List *widgets.List
}

// Pushes a given number of item into the stack
func (lw *LogWindow) Push(item string) *LogWindow {
	if lw == nil || lw.List == nil {
		panic("LogWindow or underlying list is nil")
	}
	lw.List.Rows = append([]string{item}, lw.List.Rows...)
	return lw
}

// Returns a new list for event logs
func GetEventLog() *LogWindow {
	l := widgets.NewList()
	l.Title = "Event logs"
	l.Rows = []string{}
	l.TextStyle = ui.NewStyle(ui.ColorYellow)
	l.WrapText = false
	l.SetRect(0, 0, 25, 8)
	return &LogWindow{l}
}
