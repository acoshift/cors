// Package cors provides native http middleware for cors
// see https://developer.mozilla.org/en-US/docs/Web/HTTP/Access_control_CORS
package cors

import (
	"net/http"
	"strconv"
	"strings"
	"time"
)

// Config is the cors config
type Config struct {
	AllowAllOrigins  bool
	AllowOrigins     []string
	AllowMethods     []string
	AllowHeaders     []string
	AllowCredentials bool
	ExposeHeaders    []string
	MaxAge           time.Duration
}

const (
	headerACACredentials = "Access-Control-Allow-Credentials"
	headerACAMethods     = "Access-Control-Allow-Methods"
	headerACAHeaders     = "Access-Control-Allow-Headers"
	headerACMaxAge       = "Access-Control-Max-Age"
	headerACEHeaders     = "Access-Control-Expose-Headers"
	headerACAOrigin      = "Access-Control-Allow-Origin"
	headerVary           = "Vary"
	headerOrigin         = "Origin"
	headerACRMethos      = "Access-Control-Request-Method"
	headerACRHeaders     = "Access-Control-Request-Headers"
)

// New creates new CORS middleware
func New(config Config) func(http.Handler) http.Handler {
	preflightHeaders := make(http.Header)
	headers := make(http.Header)
	allowOrigins := make(map[string]bool)

	if config.AllowCredentials {
		preflightHeaders.Set(headerACACredentials, "true")
		headers.Set(headerACACredentials, "true")
	}
	if len(config.AllowMethods) > 0 {
		preflightHeaders.Set(headerACAMethods, strings.Join(config.AllowMethods, ","))
	}
	if len(config.ExposeHeaders) > 0 {
		headers.Set(headerACEHeaders, strings.Join(config.ExposeHeaders, ","))
	}
	if config.MaxAge > time.Duration(0) {
		preflightHeaders.Set(headerACMaxAge, strconv.FormatInt(int64(config.MaxAge/time.Second), 10))
	}
	if config.AllowAllOrigins {
		preflightHeaders.Set(headerACAOrigin, "*")
		headers.Set(headerACAOrigin, "*")
	} else {
		preflightHeaders.Add(headerVary, headerOrigin)
		preflightHeaders.Add(headerVary, headerACRMethos)
		preflightHeaders.Add(headerVary, headerACRHeaders)
		headers.Set(headerVary, headerOrigin)

		for _, v := range config.AllowOrigins {
			allowOrigins[v] = true
		}
	}

	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if origin := r.Header.Get(headerOrigin); len(origin) > 0 {
				h := w.Header()
				if !config.AllowAllOrigins {
					if allowOrigins[origin] {
						h.Set(headerACAOrigin, origin)
					} else {
						w.WriteHeader(http.StatusForbidden)
						return
					}
				}
				var hh http.Header
				if r.Method == http.MethodOptions {
					hh = preflightHeaders
				} else {
					hh = headers
				}
				for k, v := range hh {
					h[k] = v
				}
			}
			h.ServeHTTP(w, r)
		})
	}
}
