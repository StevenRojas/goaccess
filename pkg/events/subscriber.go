package events

import (
	"fmt"
	"sync"

	"github.com/StevenRojas/goaccess/pkg/entities"
)

// SubscriberFeed interface to handle subscribers and messages
type SubscriberFeed interface {
	Subscribe(l chan *entities.RoleEvent) Subscription
	Send(message *entities.RoleEvent)
}

// Feed struct
type Feed struct {
	lock      sync.Mutex
	listeners []chan *entities.RoleEvent
}

type sub struct {
	feed    *Feed
	index   int
	channel chan *entities.RoleEvent
	once    sync.Once
	err     chan error
}

// Subscription interface to allows unsubscribe from the events and get errors
type Subscription interface {
	Unsubscribe()
	Err() <-chan error
}

func NewSubscriber() SubscriberFeed {
	fmt.Println("NewSubscriber....")
	return &Feed{}
}

// Subscribe method to subscribe listeners
func (f *Feed) Subscribe(l chan *entities.RoleEvent) Subscription {
	fmt.Println("Subscribe....")
	f.lock.Lock()
	defer f.lock.Unlock()
	f.listeners = append(f.listeners, l)
	return &sub{
		feed:    f,
		index:   len(f.listeners) - 1,
		channel: l,
		err:     make(chan error, 1),
	}
}

// Send method to send a message to the listeners
func (f *Feed) Send(message *entities.RoleEvent) {
	fmt.Printf("Sending....%v\n", message.RoleID)
	f.lock.Lock()
	defer f.lock.Unlock()
	for _, l := range f.listeners {
		l <- message
	}
}

func (f *Feed) remove(i int) {
	f.lock.Lock()
	defer f.lock.Unlock()
	last := len(f.listeners) - 1
	f.listeners[i] = f.listeners[last]
	f.listeners = f.listeners[:last]
}

// Unsubscribe method to unsubscribe from the event
func (s *sub) Unsubscribe() {
	s.once.Do(func() {
		s.feed.remove(s.index)
		close(s.err)
	})
}

// Err method which returns the error channel
func (s *sub) Err() <-chan error {
	return s.err
}
