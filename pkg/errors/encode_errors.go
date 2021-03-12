package errors

import (
	"fmt"

	"github.com/StevenRojas/goaccess/pkg/configuration"
)

// ToHTTP parse an error and return a HTTP code and message
func ToHTTP(err error, logger configuration.LoggerWrapper) (int, interface{}) {
	switch e := err.(type) {
	case HTTPError:
		logger.Info(fmt.Sprintf("HTTP Error handled: %d %s", e.Code, e.Message), e.Err)
		return e.Code, e
	case ValidationError:
		logger.Info(fmt.Sprintf("Validation Error handled: %d [%v]", 400, e))
		return 400, e
	}
	return 500, err.Error()
}

// ToGRPC parse an error and return an error gRCP representation
func ToGRPC(err error, logger configuration.LoggerWrapper) (int, string) {
	return 0, err.Error()
}
