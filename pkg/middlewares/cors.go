package middlewares

import (
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

const (
	AllowedHeaders       = "Authorization, Accept, Accept-Encoding, Accept-Language, Content-Length, Content-Type, Origin, Referer"
	AllowedMethods       = "POST, PUT, PATCH, OPTIONS, GET, DELETE"
	HeaderAllowedHeaders = "Access-Control-Allow-Headers"
	HeaderAllowedMethods = "Access-Control-Allow-Methods"
	HeaderAllowedOrigin  = "Access-Control-Allow-Origin"
	ContentTypeHeader    = "Content-Type"
	ContentTypeValue     = "application/json; charset=utf-8"
)

// AllowedOrigins allowed origins
var AllowedOrigins = []string{"localhost"}

// CORSPolicies middleware. Sets the 3 type of response headers for routes with Options
func CORSPolicies(methods []string) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			var allowedOrigin = ""
			origin := req.Header.Get("Origin")
			for _, ao := range AllowedOrigins {
				if strings.Contains(origin, ao) {
					allowedOrigin = origin
				}
			}
			for _, method := range methods {
				if http.MethodOptions == method {
					w.Header().Add(HeaderAllowedOrigin, allowedOrigin)
					if http.MethodOptions == req.Method {
						w.Header().Add(HeaderAllowedHeaders, AllowedHeaders)
						w.Header().Add(HeaderAllowedMethods, strings.Join(methods, ","))
					}
				}
			}
			next.ServeHTTP(w, req)
		})
	}
}

func ContentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add(ContentTypeHeader, ContentTypeValue)
		next.ServeHTTP(w, r)
	})
}
