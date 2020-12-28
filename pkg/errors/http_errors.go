package errors

const (
	errBadRequest            = "Unable to process request"
	errUnauthorizedOperation = "Unauthorized operation"
	errNotFound              = "Resource not found"
	errForbiddenOperation    = "Forbidden operation"
	errMethodNotAllowed      = "Method Not Allowed"
	errNotAcceptable         = "Not Acceptable"
	errRequestTimeout        = "Request Timeout"
	errConflict              = "Conflict"
	errTooManyRequests       = "Too Many Requests"
	errInternalServerError   = "Internal Server Error"
	errBadGateway            = "Bad Gateway"
	errServiceUnavailable    = "Service Unavailable"
	errGatewayTimeout        = "GatewayTimeout"
)

// HTTPError represents an HTTP error
type HTTPError struct {
	Code    int    `json:"-"`
	Message string `json:"message"`
	Err     error  `json:"-"`
}

func (httpError *HTTPError) Error() string {
	return httpError.Message
}

// HTTPBadRequest wrapper to 400 error
func HTTPBadRequest(err error) *HTTPError {
	return &HTTPError{
		Code:    400,
		Message: errBadRequest,
		Err:     err,
	}
}

// HTTPUnauthorized wrapper to 401 error
func HTTPUnauthorized(err error) *HTTPError {
	return &HTTPError{
		Code:    401,
		Message: errUnauthorizedOperation,
		Err:     err,
	}
}

// HTTPForbidden wrapper to 403 error
func HTTPForbidden(err error) *HTTPError {
	return &HTTPError{
		Code:    403,
		Message: errBadRequest,
		Err:     err,
	}
}

// HTTPNotFound wrapper to 404 error
func HTTPNotFound(err error) *HTTPError {
	return &HTTPError{
		Code:    404,
		Message: errNotFound,
		Err:     err,
	}
}

// HTTPMethodNotAllowed wrapper to 405 error
func HTTPMethodNotAllowed(err error) *HTTPError {
	return &HTTPError{
		Code:    405,
		Message: errMethodNotAllowed,
		Err:     err,
	}
}

// HTTPNotAcceptable wrapper to 406 error
func HTTPNotAcceptable(err error) *HTTPError {
	return &HTTPError{
		Code:    405,
		Message: errNotAcceptable,
		Err:     err,
	}
}

// HTTPRequestTimeout wrapper to 408 error
func HTTPRequestTimeout(err error) *HTTPError {
	return &HTTPError{
		Code:    408,
		Message: errRequestTimeout,
		Err:     err,
	}
}

// HTTPConflict wrapper to 409 error
func HTTPConflict(err error) *HTTPError {
	return &HTTPError{
		Code:    409,
		Message: errConflict,
		Err:     err,
	}
}

// HTTPTooManyRequests wrapper to 429 error
func HTTPTooManyRequests(err error) *HTTPError {
	return &HTTPError{
		Code:    429,
		Message: errTooManyRequests,
		Err:     err,
	}
}

// HTTPInternalServerError wrapper to 500 error
func HTTPInternalServerError(err error) *HTTPError {
	return &HTTPError{
		Code:    500,
		Message: errInternalServerError,
		Err:     err,
	}
}

// HTTPBadGateway wrapper to 502 error
func HTTPBadGateway(err error) *HTTPError {
	return &HTTPError{
		Code:    502,
		Message: errBadGateway,
		Err:     err,
	}
}

// HTTPServiceUnavailable wrapper to 503 error
func HTTPServiceUnavailable(err error) *HTTPError {
	return &HTTPError{
		Code:    503,
		Message: errServiceUnavailable,
		Err:     err,
	}
}

// HTTPGatewayTimeout wrapper to 504 error
func HTTPGatewayTimeout(err error) *HTTPError {
	return &HTTPError{
		Code:    504,
		Message: errGatewayTimeout,
		Err:     err,
	}
}
