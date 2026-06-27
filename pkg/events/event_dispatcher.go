package events

import (
	"errors"
)

type EventDispatcher struct {
	handlers map[string][]EventHandlerInterface
}

func NewEventDispatcher() *EventDispatcher {
	return &EventDispatcher{
		handlers: make(map[string][]EventHandlerInterface),
	}
}

func (ed *EventDispatcher) Register(eventName string, handler EventHandlerInterface) error {
	if _, exists := ed.handlers[eventName]; !exists {
		ed.handlers[eventName] = []EventHandlerInterface{}
	}

	for _, h := range ed.handlers[eventName] {
		if h == handler {
			return errors.New("handler already registered for event " + eventName)
		}
	}

	ed.handlers[eventName] = append(ed.handlers[eventName], handler)
	return nil
}

func (ed *EventDispatcher) Dispatch(event EventInterface) error {
	if handlers, exists := ed.handlers[event.GetName()]; exists {
		for _, handler := range handlers {
			go handler.Handle(event)
		}
	}
	return nil
}

func (ed *EventDispatcher) Remove(eventName string, handler EventHandlerInterface) error {
	if _, exists := ed.handlers[eventName]; exists {
		for i, h := range ed.handlers[eventName] {
			if h == handler {
				ed.handlers[eventName] = append(ed.handlers[eventName][:i], ed.handlers[eventName][i+1:]...)
				return nil
			}
		}
	}
	return nil
}

func (ed *EventDispatcher) Has(eventName string) bool {
	return len(ed.handlers[eventName]) > 0
}

func (ed *EventDispatcher) Clear() error {
	ed.handlers = make(map[string][]EventHandlerInterface)
	return nil
}
