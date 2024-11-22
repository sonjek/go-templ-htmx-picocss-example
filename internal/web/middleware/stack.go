package middleware

import (
	"net/http"
)

type Middleware func(http.Handler) http.Handler

func CreateMiddlewareStack(ms ...Middleware) Middleware {
	return func(next http.Handler) http.Handler {
		// Call from the end of the stack to the beginning
		for i := len(ms) - 1; i >= 0; i-- {
			x := ms[i]
			next = x(next)
		}

		return next
	}
}
