package events

import (
	"sync"

	"github.com/StevenRojas/goaccess/pkg/entities"
)

// SubscriberFeed interface to handle subscribers and messages
type SubscriberFeed interface {
	Subscribe(eventType string, l chan *entities.RoleEvent) Subscription
	Send(message *entities.RoleEvent)
}

// Feed struct
type Feed struct {
	lock      sync.Mutex
	listeners map[string]chan *entities.RoleEvent
}

type sub struct {
	feed    *Feed
	channel chan *entities.RoleEvent
	once    sync.Once
	err     chan error
}

// Subscription interface to allows unsubscribe from the events and get errors
type Subscription interface {
	Unsubscribe(eventType string)
	Err() <-chan error
}

func NewSubscriber() SubscriberFeed {
	return &Feed{}
}

// Subscribe method to subscribe listeners
func (f *Feed) Subscribe(eventType string, l chan *entities.RoleEvent) Subscription {
	f.lock.Lock()
	defer f.lock.Unlock()
	if f.listeners == nil {
		f.listeners = make(map[string]chan *entities.RoleEvent)
	}
	f.listeners[eventType] = l
	return &sub{
		feed:    f,
		channel: l,
		err:     make(chan error, 1),
	}
}

// Send method to send a message to the listeners
func (f *Feed) Send(message *entities.RoleEvent) {
	f.lock.Lock()
	defer f.lock.Unlock()
	if l, ok := f.listeners[message.EventType]; ok {
		l <- message
	}
}

func (f *Feed) remove(eventType string) {
	f.lock.Lock()
	f.listeners[eventType] = nil
	f.lock.Unlock()
}

// Unsubscribe method to unsubscribe from the event
func (s *sub) Unsubscribe(eventType string) {
	s.once.Do(func() {
		s.feed.remove(eventType)
		close(s.err)
	})
}

// Err method which returns the error channel
func (s *sub) Err() <-chan error {
	return s.err
}
