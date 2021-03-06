package entities

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/thedevsaddam/govalidator"
)

var userRules = govalidator.MapData{
	"id":    []string{"required"},
	"email": []string{"required", "email"},
	"name":  []string{"required"},
}

var userRuleMessages = govalidator.MapData{
	"id": []string{
		"required:ID field is required",
	},
}

var roleRules = govalidator.MapData{
	"name": []string{"required"},
}

func InitValidator() {
	govalidator.AddCustomRule("is_admin", func(field string, rule string, message string, value interface{}) error {
		_, ok := value.(bool)
		if !ok {
			if message != "" {
				return errors.New(message)
			}
			return fmt.Errorf("The %s field must be a boolean value", field)
		}
		return nil
	})
}

func InitUserValidator(r *http.Request, userRequest User) *govalidator.Validator {
	opts := govalidator.Options{
		Request:  r,
		Data:     &userRequest,
		Rules:    userRules,
		Messages: userRuleMessages,
	}
	return govalidator.New(opts)
}

func InitRoleValidator(r *http.Request, roleRequest *Role) *govalidator.Validator {
	opts := govalidator.Options{
		Request: r,
		Data:    roleRequest,
		Rules:   roleRules,
	}
	return govalidator.New(opts)
}
