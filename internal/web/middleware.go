package web

import (
	"log"
	"net/http"
	"time"
)

type wrappedWriter struct {
	http.ResponseWriter
	statusCode int
}

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

func (w *wrappedWriter) WriteHeader(sc int) {
	w.ResponseWriter.WriteHeader(sc)
	w.statusCode = sc
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		wrapped := &wrappedWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		next.ServeHTTP(wrapped, r)

		log.Println(wrapped.statusCode, r.Method, r.URL.Path, time.Since(start))
	})
}

func DemoMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(10 * time.Millisecond)
		next.ServeHTTP(w, r)
	})
}
