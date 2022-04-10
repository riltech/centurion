package main

import (
	"os"
	"sync"

	"github.com/riltech/centurion/core"
	"github.com/riltech/centurion/core/bus"
	"github.com/sirupsen/logrus"
)

func main() {
	file, err := os.Create("logs")
	defer func() {
		file.Close()
	}()
	if err != nil {
		panic(err)
	}
	logrus.SetOutput(file)
	logrus.Info("Centurion is starting")
	exitHandler := core.NewExitHandler()
	bus := bus.NewBus()
	engine, dashboard := core.NewEngine(bus), core.NewDashboard(bus)
	wg := &sync.WaitGroup{}
	wg.Add(1)
	exitHandler.On(func() {
		logrus.Info("Running exit handler")
		engine.Stop()
		bus.Stop()
		wg.Done()
	})
	go engine.Start()
	dashboard.Start()
	wg.Wait()
}
