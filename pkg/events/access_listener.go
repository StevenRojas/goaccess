package events

import (
	"fmt"

	"github.com/StevenRojas/goaccess/pkg/entities"
)

type AccessListener interface {
	RegisterAccessListener() error
}

type access struct {
	sf SubscriberFeed
	ch chan *entities.RoleEvent
	// redis here

}

func NewAccessListener(sf SubscriberFeed) AccessListener {
	fmt.Println("NewAccessListener....")
	return &access{
		sf: sf,
	}
}

func (l *access) RegisterAccessListener() error {
	fmt.Println("Registering AccessListener....")
	l.ch = make(chan *entities.RoleEvent)
	sub := l.sf.Subscribe(l.ch)
	defer sub.Unsubscribe()
	for {
		select {
		case message := <-l.ch:
			l.processAccessMessage(message)
		case err := <-sub.Err():
			l.processAccessError(err)
		}
	}
}

func (l *access) processAccessMessage(message *entities.RoleEvent) {
	fmt.Printf("access message got: %v\n\n", message.RoleID)
}

func (l *access) processAccessError(err error) {
	fmt.Printf("error: %v", err)
}
