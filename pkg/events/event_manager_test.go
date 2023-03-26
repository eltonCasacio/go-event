package events

import (
	"sync"
	"testing"
	"time"

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

func (e *TestEvent) GetPayload() interface{} {
	return e.Payload
}

func (e *TestEvent) GetDateTime() time.Time {
	return time.Now()
}

type TestEventHandler struct {
	ID int
}

func (h *TestEventHandler) Handle(event EventInterface, wg *sync.WaitGroup) {}

type EventManagerTestSuite struct {
	suite.Suite
	event        TestEvent
	event2       TestEvent
	handler      TestEventHandler
	handler2     TestEventHandler
	handler3     TestEventHandler
	eventManager *EventManager
}

func (suite *EventManagerTestSuite) SetupSuite() {
	suite.eventManager = NewEventManager()
	suite.handler = TestEventHandler{ID: 1}
	suite.handler2 = TestEventHandler{ID: 2}
	suite.handler3 = TestEventHandler{ID: 3}
	suite.event = TestEvent{Name: "event", Payload: "payload"}
	suite.event2 = TestEvent{Name: "event2", Payload: "payload2"}
}

func (suite *EventManagerTestSuite) TestEventManager_Register() {
	err := suite.eventManager.Register(suite.event.GetName(), &suite.handler)
	suite.Nil(err)
	suite.Equal(len(suite.eventManager.handlers[suite.event.GetName()]), 1)
	suite.NotNil(suite.eventManager.handlers[suite.event.GetName()])

	err = suite.eventManager.Register(suite.event.Name, &suite.handler2)
	suite.Nil(err)
	suite.Equal(len(suite.eventManager.handlers[suite.event.GetName()]), 2)

	suite.Equal(&suite.handler, suite.eventManager.handlers[suite.event.GetName()][0])
	suite.Equal(&suite.handler2, suite.eventManager.handlers[suite.event.GetName()][1])
}

func (suite *EventManagerTestSuite) TestEventManager_WhenHandlerAlreadyExists() {
	suite.eventManager.Register(suite.event.Name, &suite.handler)
	err := suite.eventManager.Register(suite.event.Name, &suite.handler)
	suite.NotNil(err)
	suite.EqualError(err, ErrorHandlerAlreadyRegistered.Error())
}

func (suite *EventManagerTestSuite) TestEventManager_Clear() {
	// EVENTO 1
	suite.eventManager.Register(suite.event.GetName(), &suite.handler)
	suite.eventManager.Register(suite.event.GetName(), &suite.handler2)
	suite.Equal(len(suite.eventManager.handlers[suite.event.GetName()]), 2)

	// EVENTO 2
	suite.eventManager.Register(suite.event2.GetName(), &suite.handler)
	suite.Equal(len(suite.eventManager.handlers[suite.event2.GetName()]), 1)

	suite.Equal(len(suite.eventManager.handlers), 2)

	err := suite.eventManager.RemoveAll()
	suite.Nil(err)
	suite.Equal(len(suite.eventManager.handlers), 0)
}

func (suite *EventManagerTestSuite) TestEventManager_Has() {
	// EVENTO 1
	suite.eventManager.Register(suite.event.GetName(), &suite.handler)
	suite.eventManager.Register(suite.event.GetName(), &suite.handler2)
	// EVENTO 2
	suite.eventManager.Register(suite.event2.GetName(), &suite.handler)

	has := suite.eventManager.Has(suite.event.GetName(), &suite.handler)
	suite.True(has)

	has = suite.eventManager.Has(suite.event.GetName(), &suite.handler2)
	suite.True(has)

	has = suite.eventManager.Has(suite.event.GetName(), &suite.handler3)
	suite.False(has)
}

func (suite *EventManagerTestSuite) TestEventManager_HasNoHandle() {
	suite.eventManager.Register(suite.event.GetName(), &suite.handler)
	has := suite.eventManager.Has(suite.event.GetName(), &suite.handler3)
	suite.False(has)
}

func (suite *EventManagerTestSuite) TestEventManager_Has_InvalidEventName() {
	has := suite.eventManager.Has("", &suite.handler3)
	suite.False(has)
}

type MockHandler struct{ mock.Mock }

func (m *MockHandler) Handle(event EventInterface, wg *sync.WaitGroup) {
	m.Called(event)
	wg.Done()
}

func (suite *EventManagerTestSuite) TestEventManager_Dispatch() {
	eh := &MockHandler{}
	eh.On("Handle", &suite.event)
	suite.eventManager.Register(suite.event.GetName(), eh)

	eh2 := &MockHandler{}
	eh2.On("Handle", &suite.event2)
	suite.eventManager.Register(suite.event2.GetName(), eh2)

	suite.eventManager.Dispatch(&suite.event)
	suite.eventManager.Dispatch(&suite.event2)

	eh.AssertExpectations(suite.T())
	eh.AssertNumberOfCalls(suite.T(), "Handle", 1)

	eh2.AssertExpectations(suite.T())
	eh2.AssertNumberOfCalls(suite.T(), "Handle", 1)
}

func (suite *EventManagerTestSuite) TestEventManager_Remove() {
	suite.eventManager.Register(suite.event.GetName(), &suite.handler)
	suite.eventManager.Register(suite.event.GetName(), &suite.handler2)
	suite.eventManager.Register(suite.event.GetName(), &suite.handler3)
	suite.Equal(3, len(suite.eventManager.handlers[suite.event.GetName()]))

	err := suite.eventManager.Remove(suite.event.GetName(), &suite.handler2)
	suite.Nil(err)
	suite.Equal(suite.eventManager.handlers[suite.event.GetName()][0], &suite.handler)
	suite.Equal(suite.eventManager.handlers[suite.event.GetName()][1], &suite.handler3)
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(EventManagerTestSuite))
}
