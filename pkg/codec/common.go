package codec

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	conf "github.com/StevenRojas/goaccess/pkg/configuration"
	e "github.com/StevenRojas/goaccess/pkg/errors"
	gokitHTTP "github.com/go-kit/kit/transport/http"
)

// EmptyResponse Empty response, just OK or Error
type EmptyResponse struct {
	Err error `json:"error,omitempty"`
}

// NoContentResponse 204 response
type NoContentResponse struct{}

// IDResponse ID response
type IDResponse struct {
	ID string `json:"id,omitempty"`
}

// StringList list of string
type StringList struct {
	List []string `json:"list"`
}

// ListMapString map of string
type ListMapString struct {
	List []map[string]string `json:"list"`
}

// MapString map of string
type MapString struct {
	List map[string]string `json:"list"`
}

// MapStringList list of string
type MapStringList struct {
	List map[string][]string `json:"list"`
}

type MapOfMapStringList struct {
	List map[string]map[string][]string `json:"list"`
}

type Pagination struct {
	Query string
	Sort  string
	Page  int
	Limit int
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
func DecodeEmptyRequest(_ context.Context, r *http.Request) (interface{}, error) {
	return nil, nil
}

// DecodePaginatedListRequest decode paginated list request
func DecodePaginatedListRequest(_ context.Context, r *http.Request) (interface{}, error) {
	pagination := Pagination{}
	pagination.Query = r.URL.Query().Get("q")
	pagination.Sort = r.URL.Query().Get("sort")

	value := r.URL.Query().Get("page")
	page, err := strconv.Atoi(value)
	if err != nil {
		pagination.Page = 0
	} else {
		pagination.Page = page
	}

	value = r.URL.Query().Get("limit")
	limit, err := strconv.Atoi(value)
	if err != nil {
		pagination.Limit = 0
	} else {
		pagination.Limit = limit
	}

	return pagination, nil
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
		if _, ok := response.(*NoContentResponse); ok {
			w.WriteHeader(http.StatusNoContent)
			return nil
		}
		return EncodeResponse(w, response)
	}
}
