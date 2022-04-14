package core

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

// Describes an exit handler interface
type IExitHandler interface {
	// Subscribe a function to run when exit should happen
	On(func())
	// Trigger calls all functions at the same time
	Trigger()
	// Returns if exit is already been initiated
	IsRunning() bool
}

// Exit handler implementation
type ExitHandler struct {
	// Channel to receive os signals
	sigs chan os.Signal
	// Called when trigger happens
	stop chan uint8
	// Functions to run on exit
	cleanups []func()
	// Thread safety for accessing cleanups
	mux sync.Mutex
	// Status of running exit handlers
	isRunning bool
}

// Inteface check
var _ IExitHandler = (*ExitHandler)(nil)

func (e *ExitHandler) On(toRun func()) {
	if e == nil {
		panic("Cannot subscribe without an IExit instance")
	}
	defer func() {
		e.mux.Unlock()
	}()
	e.mux.Lock()
	if e.cleanups == nil {
		e.cleanups = []func(){toRun}
		return
	}
	e.cleanups = append(e.cleanups, toRun)
}

func (e *ExitHandler) Trigger() {
	if e == nil {
		return
	}
	e.stop <- 1
}

func (e *ExitHandler) IsRunning() bool {
	if e == nil {
		return false
	}
	return e.isRunning
}

// Constructor for the exit handler
func NewExitHandler() IExitHandler {
	handler := &ExitHandler{
		sigs: make(chan os.Signal, 1),
		stop: make(chan uint8, 1),
		mux:  sync.Mutex{},
	}
	signal.Notify(handler.sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		select {
		case <-handler.sigs:
		case <-handler.stop:
		}
		handler.isRunning = true
		for _, cleanup := range handler.cleanups {
			func() {
				defer func() {
					if err := recover(); err != nil {
						fmt.Printf("Error while cleaning up: %s\n", err)
					}
				}()
				cleanup()
			}()
		}
		close(handler.sigs)
		close(handler.stop)
	}()
	return handler
}
