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
	logrus.SetFormatter(&logrus.TextFormatter{
		DisableQuote: true,
	})
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
	var exampleAttacker example.IAttacker
	var exampleDefender example.IDefender
	exitHandler.On(func() {
		logrus.Info("Running exit handler")
		engine.Stop()
		bus.Stop()
		if exampleAttacker != nil {
			exampleAttacker.Stop()
		}
		if exampleDefender != nil {
			exampleDefender.Stop()
		}
		wg.Done()
	})
	go engine.Start()
	if spec.ExampleEnabled {
		exampleAttacker = example.NewAttacker()
		go func() {
			// Wait a bit before startup
			<-time.After(2 * time.Second)
			exampleAttacker.Start()
		}()
		exampleDefender = example.NewDefender()
		go func() {
			// Wait a bit less before startup
			// as the defender has to install the module
			<-time.After(1 * time.Second)
			exampleDefender.Start()
		}()
	}
	dashboard.Start()
	wg.Wait()
}
