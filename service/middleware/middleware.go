package middleware

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

type Middleware func(http.Handler) http.Handler

func LogStack(xs ...Middleware) Middleware {
	return func(next http.Handler) http.Handler {
		for i := len(xs) - 1; i > 0; i-- {
			x := xs[i]
			next = x(next)
		}
		return next
	}
}

func LogRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("%s - %s %s %s\n", r.RemoteAddr, r.Proto, r.Method, r.URL.RequestURI())

		next.ServeHTTP(w, r)
	})
}

func RecoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")
				http.Error(w, fmt.Sprintf("Internal Server Error: %s", err), http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func LogResponse(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rw := &responseLogger{ResponseWriter: w, statusCode: http.StatusOK}
		next.ServeHTTP(rw, r)

		// Logging after the request has been processed
		log.Printf("%s - %s %s %s - Status: %d, Size: %d, Duration: %s\n",
			r.RemoteAddr, r.Proto, r.Method, r.URL.RequestURI(),
			rw.statusCode, rw.size, time.Since(start))
	})
}

// Custom ResponseWriter to capture status code and response size
type responseLogger struct {
	http.ResponseWriter
	statusCode int
	size       int
}

func (rl *responseLogger) WriteHeader(statusCode int) {
	rl.statusCode = statusCode
	rl.ResponseWriter.WriteHeader(statusCode)
}

func (rl *responseLogger) Write(data []byte) (int, error) {
	size, err := rl.ResponseWriter.Write(data)
	rl.size += size
	return size, err
}
