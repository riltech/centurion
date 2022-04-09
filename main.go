package main

import (
	"github.com/riltech/centurion/core"
	"github.com/riltech/centurion/core/bus"
)

func main() {
	exitHandler := core.NewExitHandler()
	bus := bus.NewBus()
	engine, dashboard := core.NewEngine(bus), core.NewDashboard(bus)
	exitHandler.On(func() {
		engine.Stop()
		dashboard.Stop()
		bus.Stop()
	})
	dashboard.Start()
	go engine.Start()
}
