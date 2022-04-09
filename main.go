package main

import (
	"github.com/riltech/centurion/core"
)

func main() {
	exitHandler := core.NewExitHandler()
	bus := core.NewBus()
	engine, dashboard := core.NewEngine(bus), core.NewDashboard(bus)
	exitHandler.On(func() {
		engine.Stop()
		dashboard.Stop()
		bus.Stop()
	})
	go engine.Start()
	dashboard.Start()
}
