package events

import (
	"context"
	"fmt"

	"github.com/StevenRojas/goaccess/pkg/entities"
	"github.com/StevenRojas/goaccess/pkg/repository"
)

type ActionListener interface {
	RegisterActionListener() error
}

type action struct {
	sf          SubscriberFeed
	ch          chan *entities.RoleEvent
	rolesRepo   repository.RolesRepository
	actionsRepo repository.ActionsRepository
}

func NewActionListener(
	actionsRepo repository.ActionsRepository,
	rolesRepo repository.RolesRepository,
	sf SubscriberFeed,
) ActionListener {
	return &action{
		rolesRepo:   rolesRepo,
		actionsRepo: actionsRepo,
		sf:          sf,
	}
}

func (l *action) RegisterActionListener() error {
	l.ch = make(chan *entities.RoleEvent)
	sub := l.sf.Subscribe(entities.EventTypeAction, l.ch)
	defer sub.Unsubscribe(entities.EventTypeAction)
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
	ctx := context.Background()
	// TODO: error retry
	users, err := l.rolesRepo.UsersByRole(ctx, message.RoleID)
	if err != nil {
		l.processActionError(err)
	}
	if len(users) == 0 { // Means the role is not assigned to other users
		// Check if the user has other roles
		roles, err := l.rolesRepo.RolesByUser(ctx, message.UserID)
		if err != nil {
			l.processActionError(err)
		}
		if len(roles) == 0 { // remove actions for the user
			err = l.actionsRepo.RemoveActionsByUser(ctx, message.UserID)
			if err != nil {
				l.processActionError(err)
			}
		}
	} else {
		for _, userID := range users {
			err := l.actionsRepo.SetActionList(ctx, userID)
			if err != nil {
				l.processActionError(err)
			}
		}
		err = l.actionsRepo.UpdateActionList(ctx, message.RoleID)
		if err != nil {
			l.processActionError(err)
		}
	}
}

func (l *action) processActionError(err error) {
	fmt.Printf("error: %v", err)
}
