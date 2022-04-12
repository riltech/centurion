package core

import (
	"context"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/riltech/centurion/core/bus"
	"github.com/riltech/centurion/core/challenge"
	"github.com/riltech/centurion/core/combat"
	"github.com/riltech/centurion/core/engine"
	"github.com/riltech/centurion/core/player"
	"github.com/sirupsen/logrus"
)

// Describes and engine interface
type IEngine interface {
	// Starts the engine process [This is a blocking call]
	Start()
	// Stops the engine process
	Stop()
}

// Engine implementation
type Engine struct {
	bus    bus.IBus
	router *httprouter.Router
	server *http.Server

	// Internal dependencies
	ctrl engine.IConroller
}

// Interface check
var _ IEngine = (*Engine)(nil)

func (e *Engine) Start() {
	if e == nil {
		panic("Cannot start an uninitialized Engine")
	}
	e.router = e.ctrl.GetRouter()
	e.server = &http.Server{
		Handler:      e.router,
		Addr:         ":8080",
		WriteTimeout: 25 * time.Second,
		ReadTimeout:  25 * time.Second,
	}
	logrus.Infoln("Engine starts listening on 8080")
	if err := e.server.ListenAndServe(); err != nil {
		logrus.Error(err)
	}
}

func (e Engine) Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := e.server.Shutdown(ctx); err != nil {
		panic(err)
	}
}

// Constructor for Engine
func NewEngine(bus bus.IBus) IEngine {
	playerRepo := player.NewRepository()
	playerService := player.NewService(playerRepo)
	challengeRepo := challenge.NewRepository()
	challengeService := challenge.NewService(challengeRepo)
	combatRepo := combat.NewRepository()
	combatService := combat.NewService(combatRepo)
	err := challengeService.AddDefaultModules()
	if err != nil {
		logrus.Fatal(err)
	}
	engineService := engine.NewService(bus, playerService, challengeService, combatService)
	return &Engine{
		// Available after start is called
		router: nil,
		server: nil,

		// Available as the instance is created
		bus:  bus,
		ctrl: engine.NewController(bus, engineService, playerService, challengeService),
	}
}
