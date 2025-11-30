package middlewares

import (
	"fmt"
	"net/http"
	"time"
)

func ResponseTime(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Create custom responseWriter to capture status code
		wrappedWriter := &ResponseWriter{
			ResponseWriter: w,
			status:         http.StatusOK,
		}

		// Calc duration
		duration := time.Since(start)
		wrappedWriter.Header().Set("X-Response-Time", duration.String())

		next.ServeHTTP(wrappedWriter, r)

		// log req details
		fmt.Printf("Method: %s, URL: %s, Status: %d. Duration: %v\n", r.Method, r.URL, wrappedWriter.status, duration.String())
		fmt.Println("Sent Response from Response Time Middleware ")
	})
}

type ResponseWriter struct {
	http.ResponseWriter
	status int
}

func (rw *ResponseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}
