package events

import (
	"context"
	"fmt"

	"github.com/StevenRojas/goaccess/pkg/entities"
	"github.com/StevenRojas/goaccess/pkg/repository"
)

type AccessListener interface {
	RegisterAccessListener() error
}

type access struct {
	sf          SubscriberFeed
	ch          chan *entities.RoleEvent
	modulesRepo repository.ModulesRepository
	rolesRepo   repository.RolesRepository
}

func NewAccessListener(
	modulesRepo repository.ModulesRepository,
	rolesRepo repository.RolesRepository,
	sf SubscriberFeed,
) AccessListener {
	return &access{
		modulesRepo: modulesRepo,
		rolesRepo:   rolesRepo,
		sf:          sf,
	}
}

func (l *access) RegisterAccessListener() error {
	l.ch = make(chan *entities.RoleEvent)
	sub := l.sf.Subscribe(entities.EventTypeAccess, l.ch)
	defer sub.Unsubscribe(entities.EventTypeAccess)
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
	ctx := context.Background()
	// TODO: error retry
	users, err := l.rolesRepo.UsersByRole(ctx, message.RoleID)
	if err != nil {
		l.processAccessError(err)
	}
	if len(users) == 0 { // Means the role is not assigned to other users
		// Check if the user has other roles
		roles, err := l.rolesRepo.RolesByUser(ctx, message.UserID)
		if err != nil {
			l.processAccessError(err)
		}
		if len(roles) == 0 { // remove actions for the user
			err = l.modulesRepo.RemoveAccessByUser(ctx, message.UserID)
			if err != nil {
				l.processAccessError(err)
			}
		}
	} else {
		for _, userID := range users {
			err := l.modulesRepo.SetAccessList(ctx, userID)
			l.processAccessError(err)
		}
	}
}

func (l *access) processAccessError(err error) {
	fmt.Printf("error: %v", err)
}
