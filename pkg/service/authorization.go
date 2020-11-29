package service

import (
	"context"
	"errors"

	"github.com/StevenRojas/goaccess/pkg/entities"
	"github.com/StevenRojas/goaccess/pkg/events"
	"github.com/StevenRojas/goaccess/pkg/repository"
)

// AuthorizationService authorization service to handle modules, submodules and sections
type AuthorizationService interface {
	// AssignActions assign actions to a role
	AssignActions(ctx context.Context, roleID string, module string, submodule string, actions []string) error
	// UnassignActions unassign actions from a role
	UnassignActions(ctx context.Context, roleID string, module string, submodule string, actions []string) error
	// AssingRole assign role to a user
	AssignRole(ctx context.Context, userID string, roleID string) error
	// UnassignRole unassign role from a user
	UnassignRole(ctx context.Context, userID string, roleID string) error
	// GetAccessList get a json of modules, submodules and sections where the user has access
	GetAccessList(ctx context.Context, userID string) (string, error)
	// GetActionListByModule get a json list with the actions can be performed by a user in a module
	GetActionListByModule(ctx context.Context, module string, userID string) (string, error)
	// CheckPermission checks if a user has permission to perform an action
	CheckPermission(ctx context.Context, action string, userID string) (bool, error)
}

type authorization struct {
	modulesRepo    repository.ModulesRepository
	rolesRepo      repository.RolesRepository
	actionsRepo    repository.ActionsRepository
	usersRepo      repository.UsersRepository
	subscriberFeed events.SubscriberFeed
}

// NewAuthorizationService return a new authorization service instance
func NewAuthorizationService(
	modulesRepo repository.ModulesRepository,
	rolesRepo repository.RolesRepository,
	actionsRepo repository.ActionsRepository,
	usersRepo repository.UsersRepository,
	subscriberFeed events.SubscriberFeed,
) AuthorizationService {
	return &authorization{
		modulesRepo:    modulesRepo,
		rolesRepo:      rolesRepo,
		actionsRepo:    actionsRepo,
		usersRepo:      usersRepo,
		subscriberFeed: subscriberFeed,
	}
}

// AssignActions assign actions to a role
func (a *authorization) AssignActions(ctx context.Context, roleID string, module string, submodule string, actions []string) error {
	if ok, _ := a.rolesRepo.IsValidRole(ctx, roleID); !ok {
		return errors.New("Role not found")
	}
	err := a.actionsRepo.AssignActions(ctx, roleID, module, submodule, actions)
	if err != nil {
		return err
	}
	roleEvent := &entities.RoleEvent{RoleID: roleID, EventType: entities.EventTypeAction}
	go a.subscriberFeed.Send(roleEvent)
	return nil
}

// UnassignActions unassign actions from a role
func (a *authorization) UnassignActions(ctx context.Context, roleID string, module string, submodule string, actions []string) error {
	if ok, _ := a.rolesRepo.IsValidRole(ctx, roleID); !ok {
		return errors.New("Role not found")
	}
	err := a.actionsRepo.UnassignActions(ctx, roleID, module, submodule, actions)
	if err != nil {
		return err
	}
	roleEvent := &entities.RoleEvent{RoleID: roleID, EventType: entities.EventTypeAction}
	go a.subscriberFeed.Send(roleEvent)
	return nil
}

// AssingRole assign role to a user
func (a *authorization) AssignRole(ctx context.Context, userID string, roleID string) error {
	if ok, _ := a.usersRepo.IsValidUser(ctx, userID); !ok {
		return errors.New("User not found")
	}
	if ok, _ := a.rolesRepo.IsValidRole(ctx, roleID); !ok {
		return errors.New("Role not found")
	}
	err := a.rolesRepo.AssignRole(ctx, userID, roleID)
	if err != nil {
		return err
	}
	roleEvent := &entities.RoleEvent{RoleID: roleID, EventType: entities.EventTypeAccess}
	go a.subscriberFeed.Send(roleEvent)
	roleEvent = &entities.RoleEvent{RoleID: roleID, EventType: entities.EventTypeAction}
	go a.subscriberFeed.Send(roleEvent)
	return nil
}

// UnassignRole unassign role from a user
func (a *authorization) UnassignRole(ctx context.Context, userID string, roleID string) error {
	if ok, _ := a.usersRepo.IsValidUser(ctx, userID); !ok {
		return errors.New("User not found")
	}
	if ok, _ := a.rolesRepo.IsValidRole(ctx, roleID); !ok {
		return errors.New("Role not found")
	}
	err := a.rolesRepo.UnassignRole(ctx, userID, roleID)
	if err != nil {
		return err
	}
	go a.subscriberFeed.Send(&entities.RoleEvent{
		RoleID:    roleID,
		UserID:    userID,
		EventType: entities.EventTypeAccess,
	})

	go a.subscriberFeed.Send(&entities.RoleEvent{
		RoleID:    roleID,
		UserID:    userID,
		EventType: entities.EventTypeAction,
	})
	return nil
}

// GetAccessList get a json of modules, submodules and sections where the user has access
func (a *authorization) GetAccessList(ctx context.Context, userID string) (string, error) {
	access, err := a.modulesRepo.GetAccessList(ctx, userID)
	if err != nil {
		return "", err
	}
	if access == "" {
		return "", errors.New("User has not access defined")
	}
	return access, nil
}

// GetActionListByModule get a json list with the actions can be performed by a user in a module
func (a *authorization) GetActionListByModule(ctx context.Context, module string, userID string) (string, error) {
	actions, err := a.actionsRepo.GetActionListByModule(ctx, module, userID)
	if err != nil {
		return "", err
	}
	if actions == "" {
		return "", errors.New("User has not actions defined")
	}
	return actions, nil
}

// CheckPermission checks if a user has permission to perform an action
func (a *authorization) CheckPermission(ctx context.Context, action string, userID string) (bool, error) {
	return a.actionsRepo.CheckPermission(ctx, action, userID)
}
