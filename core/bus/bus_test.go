package bus

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBus(t *testing.T) {
	bus := NewBus()
	ch1, ch2, ch3 := bus.Listen("foo"), bus.Listen("bar"), bus.Listen("foo")
	fooEvent, barEvent := &BusEvent{
		Type: "foo",
	}, &BusEvent{
		Type: "bar",
	}
	bus.Send(fooEvent)
	bus.Send(barEvent)
	value := <-ch1
	assert.Equal(t, fooEvent, value)
	value = <-ch3
	assert.Equal(t, fooEvent, value)
	value = <-ch2
	assert.Equal(t, barEvent, value)
	bus.Stop()
	_, ok := <-ch1
	assert.False(t, ok)
	_, ok = <-ch2
	assert.False(t, ok)
	_, ok = <-ch3
	assert.False(t, ok)
}
