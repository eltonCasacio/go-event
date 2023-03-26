package events

import (
	"errors"
	"fmt"
	"sync"
)

var (
	ErrorHandlerAlreadyRegistered = errors.New("handler already registered")
)

type EventManager struct {
	handlers map[string][]EventHandlerInterface
}

func NewEventManager() *EventManager {
	return &EventManager{
		handlers: make(map[string][]EventHandlerInterface),
	}
}

func (e *EventManager) Register(eventName string, handler EventHandlerInterface) error {
	if _, ok := e.handlers[eventName]; ok {
		for _, h := range e.handlers[eventName] {
			if h == handler {
				return ErrorHandlerAlreadyRegistered
			}
		}
	}

	e.handlers[eventName] = append(e.handlers[eventName], handler)
	return nil
}

func (e *EventManager) Dispatch(event EventInterface) error {
	if handlers, ok := e.handlers[event.GetName()]; ok {
		wg := sync.WaitGroup{}
		wg.Add(len(handlers))
		for _, h := range handlers {
			go h.Handle(event, &wg)
		}
		wg.Wait()
	}
	return nil
}

func (e *EventManager) Has(eventName string, handler EventHandlerInterface) bool {
	events := e.handlers[eventName]
	fmt.Println(events)
	for _, h := range events {
		if h == handler {
			return true
		}
	}
	return false
}

func (e *EventManager) Remove(eventName string, handler EventHandlerInterface) error {
	handlers := e.handlers[eventName]
	for k, h := range handlers {
		if h == handler {
			handlers = append(handlers[:k], handlers[k+1:]...)
		}
	}
	e.handlers[eventName] = handlers
	return nil
}

func (e *EventManager) RemoveAll() error {
	e.handlers = make(map[string][]EventHandlerInterface)
	return nil
}
