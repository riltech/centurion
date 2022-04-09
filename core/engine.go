package core

import (
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/riltech/centurion/core/engine"
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
	bus    IBus
	router *httprouter.Router
}

// Interface check
var _ IEngine = (*Engine)(nil)

func (e Engine) Start() {
	e.router = engine.GetRouter()
	log.Fatal(http.ListenAndServe(":8080", e.router))
}

func (e Engine) Stop() {

}

// Constructor for Engine
func NewEngine(bus IBus) IEngine {
	return Engine{bus, nil}
}
