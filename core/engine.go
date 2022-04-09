package core

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/riltech/centurion/core/bus"
	"github.com/riltech/centurion/core/engine"
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
	ctrl   engine.IConroller
	router *httprouter.Router
	server *http.Server
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
	log.Fatal(e.server.ListenAndServe())
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
	return &Engine{bus, engine.NewController(bus), nil, nil}
}
