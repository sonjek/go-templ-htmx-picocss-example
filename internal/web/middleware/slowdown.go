package middleware

import (
	"net/http"
	"time"
)

func DemoMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(10 * time.Millisecond)
		next.ServeHTTP(w, r)
	})
}
