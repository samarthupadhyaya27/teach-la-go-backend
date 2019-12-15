package middleware

import (
	"net/http"
	"strings"
)

/*
 * CORS
 * This is middleware written in accordance with the W3C
 * Cross-Origin Resource Sharing specification.
 * Specifically, this is written along with the request
 * flow described here:
 * https://www.w3.org/TR/cors/#resource-requests
 */

// CORSConfig is a struct holding all relevant
// information for a proper CORS configuration.
// Fields correspond to CORS headers.
type CORSConfig struct {
	AllowedOrigins      []string
	AllowedMethods      []string
	AllowedHeaders      []string
	SupportsCredentials bool
	MaxAge              uint32
}

// GetOriginsStr returns the comma delimited
// string array of a CORSConfig struct's
// AllowedOrigins field.
func (c *CORSConfig) GetOriginsStr() string {
	return strings.Join(c.AllowedOrigins[:], ",")
}

// GetMethodsStr returns the comma delimited string
// array of a CORSConfig struct's AllowedMethods field.
func (c *CORSConfig) GetMethodsStr() string {
	return strings.Join(c.AllowedMethods[:], ",")
}

// OriginSupported returns whether the provided request origin resides
// in the list of allowed origins for a given CORSConfig.
func (c *CORSConfig) OriginSupported(requestOrigin string) bool {
	// an Origin header is mandatory.
	if requestOrigin == "" {
		return false
	}

	// the value of the Origin header must be a case-sensitive
	// match for a supported origin.
	for _, o := range c.AllowedOrigins {
		if requestOrigin == o {
			return true
		}
	}
	return false
}

// MethodSupported returns whether the provided request method is a
// case-sensitive match for any supported methods for a given CORSConfig.
func (c *CORSConfig) MethodSupported(requestMethod string) bool {
	// an Access-Control-Request-Method is mandatory.
	if requestMethod == "" {
		return false
	}

	// the value of the header must be a case-sensitive
	// match for a supported method.
	for _, m := range c.AllowedMethods {
		if requestMethod == m {
			return true
		}
	}
	return false
}

// HeadersSupported returns whether each member of a list of request header field names
// has an ASCII case-insensitive match for any AllowedHeaders value in a given CORSConfig.
func (c *CORSConfig) HeadersSupported(requestHeaderFieldNames []string) bool {
	// the empty list is permissible.
	if len(requestHeaderFieldNames) == 0 {
		return true
	}

	// the header field names must be a case-insensitive match for any of
	// the values in the list of supported headers.
	for _, requestHeader := range requestHeaderFieldNames {
		supported := false

		for _, supportedHeader := range c.AllowedHeaders {
			if strings.ToUpper(requestHeader) == strings.ToUpper(supportedHeader) {
				supported = true
				break
			}
		}

		if !supported {
			return false
		}
	}

	return true
}

// WithCORSConfig is middleware that handles your CORS preflight requests quickly
// and effectively with the supplied configuration. It is not verbose. To enable
// verbosity, please wrap it with some sort of request logging middleware.
func WithCORSConfig(next http.Handler, c CORSConfig) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// handle preflight request.
		if r.Method == http.MethodOptions {
			// acquire relevant fields.
			origin := r.Header.Get("Origin")
			method := r.Header.Get("Access-Control-Request-Method")
			headerFieldNames := strings.Split(r.Header.Get("Access-Control-Request-Headers"), ", ")

			// request origin, method must be a **case-sensitive** match
			// in those that are supported.
			// request headers must be a case **in**sensitive match.
			// if any conditions fail, we throw the request out.
			if !c.OriginSupported(origin) ||
				!c.MethodSupported(method) ||
				!c.HeadersSupported(headerFieldNames) {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			// Future consideration:
			// in addition to checking the Origin header, we should check the Host header to ensure
			// that the host name provided matches the host name on which the reosuce resides.

			// if the resource supports credentials, add a single Access-Control-Allow-Origin
			// header, with the value of the Origin header as value.
			// also add a single Access-Control-Allow-Credentials header with "true" as value.
			// otherwise, add a single Access-Control-Allow-Origin header, with the value of the
			// Origin header.
			if c.SupportsCredentials {
				w.Header().Set("Access-Control-Allow-Credentials", "true")
			}
			w.Header().Set("Access-Control-Allow-Origin", origin)

			// if provided a MaxAge, add a single Access-Control-Max-Age header.
			if c.MaxAge != 0 {
				w.Header().Set("Access-Control-Max-Age", string(c.MaxAge))
			}

			w.WriteHeader(http.StatusOK)
			return
		}

		// otherwise, serve actual request.
		next.ServeHTTP(w, r)
	})
}

// WithCORS is middleware that handles your CORS preflight requests quickly
// and effectively using default settings. It is not verbose. To enable
// verbosity, please wrap it with some sort of request logging middleware.
func WithCORS(next http.Handler) http.Handler {
	// by default we allow all methods, all origins,
	// and Content-Type headers.
	// MaxAge is omitted, and credentials are not
	// supported.
	defaultCfg := CORSConfig{
		AllowedHeaders: []string{
			"Content-Type",
		},
		AllowedMethods: []string{
			http.MethodConnect,
			http.MethodDelete,
			http.MethodGet,
			http.MethodHead,
			http.MethodOptions,
			http.MethodPatch,
			http.MethodPost,
			http.MethodPut,
			http.MethodTrace,
		},
		AllowedOrigins: []string{"*"},
	}

	return WithCORSConfig(next, defaultCfg)
}