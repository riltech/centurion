package bus

import (
	"strings"

	"github.com/davecgh/go-spew/spew"
	"github.com/sirupsen/logrus"
)

// Describes the interface of the common event bus
type IBus interface {
	// Provides a channel that can be listened for a given topic
	Listen(topic string) <-chan *BusEvent
	// Send a new event to the event bus
	Send(event *BusEvent)
	// Graceful shutdown
	Stop()
}

// Event bus implementation
type Bus struct {
	// Used for event distribution
	main chan *BusEvent
	// Distribution to individual listeners
	listeners map[string][]chan *BusEvent
	// Indicates that the bus should be stopped
	stop chan uint8
}

// Inteface check
var _ IBus = (*Bus)(nil)

// Constructor for event bus
func NewBus() IBus {
	bus := &Bus{
		main:      make(chan *BusEvent, 25),
		listeners: make(map[string][]chan *BusEvent),
		stop:      make(chan uint8, 1),
	}
	go bus.distribute()
	return bus
}

func (b Bus) Listen(topic string) <-chan *BusEvent {
	ch := make(chan *BusEvent, 2)
	arr, ok := b.listeners[strings.ToLower(topic)]
	if ok {
		b.listeners[strings.ToLower(topic)] = append(arr, ch)
		return ch
	}
	b.listeners[strings.ToLower(topic)] = []chan *BusEvent{ch}
	return ch
}

func (b Bus) Send(event *BusEvent) {
	logrus.Infof("New bus event: %s", spew.Sdump(event))
	b.main <- event
}

// Used for incoming event distribution for listeners
func (b *Bus) distribute() {
	for {
		select {
		case value := <-b.main:
			if value == nil {
				panic("Bus value cannot be null")
			}
			if listeners, ok := b.listeners[strings.ToLower(value.Type)]; ok {
				for _, listener := range listeners {
					listener <- value
				}
			}
		case <-b.stop:
			logrus.Infoln("Bus distribution is stopping!")
			return
		}
	}
}

func (b Bus) Stop() {
	b.stop <- 1
	close(b.main)
	for _, channels := range b.listeners {
		for _, ch := range channels {
			close(ch)
		}
	}
	logrus.Infoln("Bus is stopped!")
}
