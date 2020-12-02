package repository

import (
	"context"
	"testing"

	"github.com/StevenRojas/goaccess/pkg/entities"
	"github.com/StevenRojas/goaccess/pkg/utils"

	"github.com/stretchr/testify/assert"
)

func TestGetUserById(t *testing.T) {
	r := new(UsersRepoMock)
	expected := &entities.User{
		ID:    "1",
		Email: "srojas@gmail.com",
		Name:  "steven rojas",
	}
	r.M.On("GetUserByID", "22").Return(expected, nil)
	user, e := r.GetUserByID(context.TODO(), "22")
	assert.Equal(t, expected, user)
	assert.Nil(t, e)
}

func TestGetUserByEmail(t *testing.T) {
	r := new(UsersRepoMock)
	expected := &entities.User{
		ID:    "1",
		Email: "srojas@gmail.com",
		Name:  "steven rojas",
	}
	r.M.On("GetUserByEmail", "srojas@gmail.com").Return(expected, nil)
	user, e := r.GetUserByEmail(context.TODO(), "srojas@gmail.com")
	assert.Equal(t, expected, user)
	assert.Nil(t, e)
}

func TestGetUserByToken(t *testing.T) {
	r := new(UsersRepoMock)
	expected := &entities.User{
		ID:    "1",
		Email: "srojas@gmail.com",
		Name:  "steven rojas",
	}
	r.M.On("GetUserByToken", "token_here").Return(expected, nil)
	user, e := r.GetUserByToken(context.TODO(), "token_here")
	assert.Equal(t, expected, user)
	assert.Nil(t, e)
}

func TestStoreTokens(t *testing.T) {
	r := new(UsersRepoMock)
	token := &utils.StoredToken{
		ID:             "1",
		AccessToken:    "a_jwt",
		AccessUUID:     "a_uuid",
		AccessExpires:  10,
		RefreshToken:   "r_jwt",
		RefreshUUID:    "r_uuid",
		RefreshExpires: 20,
	}
	r.M.On("StoreTokens", token).Return(nil)
	e := r.StoreTokens(context.TODO(), token)
	assert.Nil(t, e)
}

func TestDeleteToken(t *testing.T) {
	r := new(UsersRepoMock)
	r.M.On("DeleteToken", "token_key").Return(nil)
	e := r.DeleteToken(context.TODO(), "token_key")
	assert.Nil(t, e)
}
