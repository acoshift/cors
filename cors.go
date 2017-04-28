// Package cors provides native http middleware for cors
// see https://developer.mozilla.org/en-US/docs/Web/HTTP/Access_control_CORS
package cors

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/acoshift/header"
	"github.com/acoshift/middleware"
)

// Config is the cors config
type Config struct {
	Skipper          middleware.Skipper
	AllowAllOrigins  bool
	AllowOrigins     []string
	AllowMethods     []string
	AllowHeaders     []string
	AllowCredentials bool
	ExposeHeaders    []string
	MaxAge           time.Duration
}

// New creates new CORS middleware
func New(config Config) func(http.Handler) http.Handler {
	if config.Skipper == nil {
		config.Skipper = middleware.DefaultSkipper
	}

	preflightHeaders := make(http.Header)
	headers := make(http.Header)
	allowOrigins := make(map[string]bool)

	if config.AllowCredentials {
		preflightHeaders.Set(header.AccessControlAllowCredentials, "true")
		headers.Set(header.AccessControlAllowCredentials, "true")
	}
	if len(config.AllowMethods) > 0 {
		preflightHeaders.Set(header.AccessControlAllowMethods, strings.Join(config.AllowMethods, ","))
	}
	if len(config.AllowHeaders) > 0 {
		preflightHeaders.Set(header.AccessControlAllowHeaders, strings.Join(config.AllowHeaders, ","))
	}
	if len(config.ExposeHeaders) > 0 {
		headers.Set(header.AccessControlExposeHeaders, strings.Join(config.ExposeHeaders, ","))
	}
	if config.MaxAge > time.Duration(0) {
		preflightHeaders.Set(header.AccessControlMaxAge, strconv.FormatInt(int64(config.MaxAge/time.Second), 10))
	}
	if config.AllowAllOrigins {
		preflightHeaders.Set(header.AccessControlAllowOrigin, "*")
		headers.Set(header.AccessControlAllowOrigin, "*")
	} else {
		preflightHeaders.Add(header.Vary, header.Origin)
		preflightHeaders.Add(header.Vary, header.AccessControlRequestMethod)
		preflightHeaders.Add(header.Vary, header.AccessControlRequestHeaders)
		headers.Set(header.Vary, header.Origin)

		for _, v := range config.AllowOrigins {
			allowOrigins[v] = true
		}
	}

	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if config.Skipper(r) {
				h.ServeHTTP(w, r)
				return
			}

			if origin := r.Header.Get(header.Origin); len(origin) > 0 {
				h := w.Header()
				if !config.AllowAllOrigins {
					if allowOrigins[origin] {
						h.Set(header.AccessControlAllowOrigin, origin)
					} else {
						w.WriteHeader(http.StatusForbidden)
						return
					}
				}
				if r.Method == http.MethodOptions {
					for k, v := range preflightHeaders {
						h[k] = v
					}
					w.WriteHeader(http.StatusOK)
					return
				}
				for k, v := range headers {
					h[k] = v
				}
			}
			h.ServeHTTP(w, r)
		})
	}
}
