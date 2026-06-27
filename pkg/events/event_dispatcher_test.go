package events

import (
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type TestEvent struct {
	Name    string
	Payload interface{}
}

func (e *TestEvent) GetName() string {
	return e.Name
}

func (e *TestEvent) GetDataTime() time.Time {
	return time.Now()
}

func (e *TestEvent) GetPayload() interface{} {
	return e.Payload
}

type TestEventHandler struct {
	ID string
}

func (h *TestEventHandler) Handle(event EventInterface) {
}

type EventDispatcherTestSuite struct {
	suite.Suite
	event           TestEvent
	event2          TestEvent
	handler         TestEventHandler
	handler2        TestEventHandler
	handler3        TestEventHandler
	eventDispatcher *EventDispatcher
}

func (s *EventDispatcherTestSuite) SetupTest() {
	s.eventDispatcher = NewEventDispatcher()
	s.handler = TestEventHandler{
		ID: "1",
	}
	s.handler2 = TestEventHandler{
		ID: "2",
	}
	s.handler3 = TestEventHandler{
		ID: "3",
	}
	s.event = TestEvent{
		Name:    "test",
		Payload: "test",
	}
	s.event2 = TestEvent{
		Name:    "test2",
		Payload: "test2",
	}
}

func (s *EventDispatcherTestSuite) TestEventDispatcher_Register() {
	err := s.eventDispatcher.Register(s.event.GetName(), &s.handler)
	s.Nil(err)
	s.True(s.eventDispatcher.Has(s.event.GetName()))
	s.Equal(1, len(s.eventDispatcher.handlers[s.event.GetName()]))

	err = s.eventDispatcher.Register(s.event.GetName(), &s.handler2)
	s.Nil(err)
	s.True(s.eventDispatcher.Has(s.event.GetName()))
	s.Equal(2, len(s.eventDispatcher.handlers[s.event.GetName()]))

	err = s.eventDispatcher.Register(s.event.GetName(), &s.handler3)
	s.Nil(err)
	s.True(s.eventDispatcher.Has(s.event.GetName()))
	s.Equal(3, len(s.eventDispatcher.handlers[s.event.GetName()]))

	assert.Equal(s.T(), &s.handler, s.eventDispatcher.handlers[s.event.GetName()][0])
	assert.Equal(s.T(), &s.handler2, s.eventDispatcher.handlers[s.event.GetName()][1])
	assert.Equal(s.T(), &s.handler3, s.eventDispatcher.handlers[s.event.GetName()][2])
}

func (s *EventDispatcherTestSuite) TestEventDispatcher_Register_WithSameHandler() {
	err := s.eventDispatcher.Register(s.event.GetName(), &s.handler)
	s.Nil(err)
	s.Equal(1, len(s.eventDispatcher.handlers[s.event.GetName()]))

	err = s.eventDispatcher.Register(s.event.GetName(), &s.handler)
	s.Equal(errors.New("handler already registered for event "+s.event.GetName()), err)
	s.Equal(1, len(s.eventDispatcher.handlers[s.event.GetName()]))
}

func (s *EventDispatcherTestSuite) TestEventDispatcher_Clear() {
	s.eventDispatcher.Register(s.event.GetName(), &s.handler)
	s.eventDispatcher.Register(s.event.GetName(), &s.handler2)
	s.eventDispatcher.Register(s.event.GetName(), &s.handler3)
	s.True(s.eventDispatcher.Has(s.event.GetName()))
	s.Equal(3, len(s.eventDispatcher.handlers[s.event.GetName()]))
	s.eventDispatcher.Clear()
	s.False(s.eventDispatcher.Has(s.event.GetName()))
}

func (s *EventDispatcherTestSuite) TestEventDispatcher_Has() {
	s.False(s.eventDispatcher.Has(s.event.GetName()))
	s.eventDispatcher.Register(s.event.GetName(), &s.handler)
	s.True(s.eventDispatcher.Has(s.event.GetName()))
}

type MockHandler struct {
	mock.Mock
}

func (m *MockHandler) Handle(event EventInterface) {
	m.Called(event)
}

func (s *EventDispatcherTestSuite) TestEventDispatcher_Dispatch() {
	handler := &MockHandler{}
	var wg sync.WaitGroup
	wg.Add(1)

	handler.On("Handle", &s.event).Run(func(args mock.Arguments) {
		wg.Done()
	}).Return(nil)

	s.eventDispatcher.Register(s.event.GetName(), handler)
	s.eventDispatcher.Dispatch(&s.event)

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(time.Second):
		s.Fail("handler was not called")
	}

	handler.AssertNumberOfCalls(s.T(), "Handle", 1)
	handler.AssertCalled(s.T(), "Handle", &s.event)
	handler.AssertExpectations(s.T())
}

func (s *EventDispatcherTestSuite) TestEventDispatcher_Remove() {
	err := s.eventDispatcher.Register(s.event.GetName(), &s.handler)
	s.Nil(err)
	s.Equal(1, len(s.eventDispatcher.handlers[s.event.GetName()]))

	err = s.eventDispatcher.Register(s.event.GetName(), &s.handler2)
	s.Nil(err)
	s.Equal(2, len(s.eventDispatcher.handlers[s.event.GetName()]))

	err = s.eventDispatcher.Register(s.event2.GetName(), &s.handler3)
	s.Nil(err)
	s.Equal(1, len(s.eventDispatcher.handlers[s.event2.GetName()]))

	err = s.eventDispatcher.Remove(s.event.GetName(), &s.handler)
	s.Nil(err)
	s.Equal(1, len(s.eventDispatcher.handlers[s.event.GetName()]))
	assert.Equal(s.T(), &s.handler2, s.eventDispatcher.handlers[s.event.GetName()][0])

	err = s.eventDispatcher.Remove(s.event2.GetName(), &s.handler3)
	s.Nil(err)
	s.Equal(0, len(s.eventDispatcher.handlers[s.event2.GetName()]))
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(EventDispatcherTestSuite))
}
