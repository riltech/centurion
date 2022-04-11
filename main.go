package main

import (
	"os"
	"sync"
	"time"

	"github.com/riltech/centurion/core"
	"github.com/riltech/centurion/core/bus"
	"github.com/riltech/centurion/core/config"
	"github.com/riltech/centurion/example"
	"github.com/sirupsen/logrus"
)

func main() {
	file, err := os.Create("logs")
	if err != nil {
		logrus.Fatal(err)
	}
	defer func() {
		file.Close()
	}()
	logrus.SetOutput(file)
	logrus.Info("Centurion is starting")
	spec, err := config.Init()
	if err != nil {
		logrus.Fatal(err)
	}
	exitHandler := core.NewExitHandler()
	bus := bus.NewBus()
	engine, dashboard := core.NewEngine(bus), core.NewDashboard(bus)
	wg := &sync.WaitGroup{}
	wg.Add(1)
	var exampleClient example.IClient
	exitHandler.On(func() {
		logrus.Info("Running exit handler")
		engine.Stop()
		bus.Stop()
		if exampleClient != nil {
			exampleClient.Stop()
		}
		wg.Done()
	})
	go engine.Start()
	if spec.ExampleEnabled {
		exampleClient = example.NewClient()
		go func() {
			// Wait a bit before startup
			<-time.After(2 * time.Second)
			exampleClient.Start()
		}()
	}
	dashboard.Start()
	wg.Wait()
}
