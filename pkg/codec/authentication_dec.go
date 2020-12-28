package codec

import (
	"context"
	"net/http"

	"github.com/StevenRojas/goaccess/pkg/entities"
	e "github.com/StevenRojas/goaccess/pkg/errors"
)

// DecodeRegisterUserRequest decode request
func DecodeRegisterUserRequest(_ context.Context, r *http.Request) (interface{}, error) {
	defer r.Body.Close()
	userRequest := entities.User{}

	errs := entities.InitUserValidator(r, userRequest).ValidateJSON()
	if len(errs) > 0 {
		return nil, e.ValidationError{
			Message: "Validation error for User",
			Fields:  errs,
		}
	}
	return userRequest, nil
}
