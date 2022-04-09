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
}

// Exit handler implementation
type ExitHandler struct {
	// Channel to receive os signals
	sigs chan os.Signal
	// Functions to run on exit
	cleanups []func()
	// Thread safety for accessing cleanups
	mux sync.Mutex
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

// Constructor for the exit handler
func NewExitHandler() IExitHandler {
	handler := &ExitHandler{
		sigs: make(chan os.Signal, 1),
		mux:  sync.Mutex{},
	}
	signal.Notify(handler.sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-handler.sigs
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
	}()
	return handler
}
