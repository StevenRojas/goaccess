package codec

import (
	"context"
	"encoding/json"
	"net/http"

	conf "github.com/StevenRojas/goaccess/pkg/configuration"
	e "github.com/StevenRojas/goaccess/pkg/errors"
	gokitHTTP "github.com/go-kit/kit/transport/http"
)

// EmptyResponse Empty response, just OK or Error
type EmptyResponse struct {
	Err error `json:"error,omitempty"`
}

type errorWrapper struct {
	Error interface{} `json:"error"`
}

type errorer interface {
	error() error
}

// EncodeResponse generic encoder
func EncodeResponse(w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}

// DecodeEmptyRequest generic decoder for request with no parameteres, like GET /users
func DecodeEmptyRequest(r *http.Request) (interface{}, error) {
	return nil, nil
}

// HTTPErrorEncoder generic error handler for unhandled exceptions
func HTTPErrorEncoder(logger conf.LoggerWrapper) gokitHTTP.ErrorEncoder {
	return func(ctx context.Context, err error, w http.ResponseWriter) {
		if headers, ok := err.(gokitHTTP.Headerer); ok {
			for k, values := range headers.Headers() {
				for _, v := range values {
					w.Header().Add(k, v)
				}
			}
		}
		code, message := e.ToHTTP(err, logger)
		w.WriteHeader(code)
		EncodeResponse(w, errorWrapper{Error: message})

	}
}

// JSONEncoder json encoder that handle errors and valid responses
func JSONEncoder(logger conf.LoggerWrapper) gokitHTTP.EncodeResponseFunc {
	return func(ctx context.Context, w http.ResponseWriter, response interface{}) error {
		if err, ok := response.(errorer); ok && err.error() != nil {
			logger.Info("handling JSONEncoder http error", "error", err.error())
			code, message := e.ToHTTP(err.error(), logger)
			w.WriteHeader(code)
			EncodeResponse(w, errorWrapper{Error: message})
			return nil
		}
		return EncodeResponse(w, response)
	}
}
