package events

import (
	"sync"
	"time"
)

type EventInterface interface {
	GetName() string
	GetDateTime() time.Time
	GetPayload() interface{} // event data
}

type EventHandlerInterface interface {
	Handle(event EventInterface, wg *sync.WaitGroup) //execute operation
}

type EventManagerInterface interface {
	Register(eventName string, handler EventHandlerInterface) error
	Dispatch(event EventInterface) error
	Has(eventName string, handler EventHandlerInterface) bool
	Remove(eventName string, handler EventHandlerInterface) error
	RemoveAll() error
}
