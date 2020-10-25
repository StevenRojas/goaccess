package events

import (
	"fmt"

	"github.com/StevenRojas/goaccess/pkg/entities"
)

type ActionListener interface {
	RegisterActionListener() error
}

type action struct {
	sf SubscriberFeed
	ch chan *entities.RoleEvent
	// redis here

}

func NewActionListener(sf SubscriberFeed) ActionListener {
	fmt.Println("NewAccessListener....")
	return &action{
		sf: sf,
	}
}

func (l *action) RegisterActionListener() error {
	fmt.Println("Registering ActionListener....")
	l.ch = make(chan *entities.RoleEvent)
	sub := l.sf.Subscribe(l.ch)
	defer sub.Unsubscribe()
	for {
		select {
		case message := <-l.ch:
			l.processActionMessage(message)
		case err := <-sub.Err():
			l.processActionError(err)
		}
	}
}

func (l *action) processActionMessage(message *entities.RoleEvent) {
	fmt.Printf("action message got: %v\n\n", message.RoleID)
}

func (l *action) processActionError(err error) {
	fmt.Printf("error: %v", err)
}
