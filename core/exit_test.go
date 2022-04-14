package core

import (
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestExitHandler(t *testing.T) {
	funcRan := make(chan uint8, 1)
	exit := NewExitHandler()
	exit.On(func() {
		funcRan <- 1
	})
	go func() {
		exit.(*ExitHandler).sigs <- syscall.SIGINT
	}()
	select {
	case <-time.After(5 * time.Second):
		t.Fatal("Exit handler did not run cleanup in time")
	case value := <-funcRan:
		assert.Equal(t, uint8(1), value)
	}

	exit = NewExitHandler()
	exit.On(func() {
		funcRan <- 1
	})
	exit.Trigger()
	select {
	case <-time.After(5 * time.Second):
		t.Fatal("Exit handler did not run cleanup in time")
	case value := <-funcRan:
		assert.Equal(t, uint8(1), value)
	}
}
