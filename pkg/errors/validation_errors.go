package errors

import (
	"encoding/json"
)

// ValidationError represents a validation error
type ValidationError struct {
	Message string      `json:"message"`
	Fields  interface{} `json:"errors"`
}

func (validationError ValidationError) Error() string {
	j, _ := json.Marshal(validationError)
	return string(j)
}
