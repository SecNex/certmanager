package logger

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

type HTTPLogEntry struct {
	RemoteAddr   string
	RequestTime  time.Time
	Method       string
	Path         string
	Protocol     string
	StatusCode   int
	ResponseSize int64
	ResponseTime time.Duration
	UserAgent    string
	Referer      string
}

func FormatHTTPLog(entry HTTPLogEntry) string {
	return fmt.Sprintf("%s - - \"%s %s %s\" %d %d \"%s\" \"%s\" %.3f",
		entry.RemoteAddr,
		entry.Method,
		entry.Path,
		entry.Protocol,
		entry.StatusCode,
		entry.ResponseSize,
		entry.Referer,
		entry.UserAgent,
		entry.ResponseTime.Seconds(),
	)
}

func LogHTTPRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		rw := &responseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		next.ServeHTTP(rw, r)

		entry := HTTPLogEntry{
			RemoteAddr:   r.RemoteAddr,
			Method:       r.Method,
			Path:         r.URL.Path,
			Protocol:     r.Proto,
			StatusCode:   rw.statusCode,
			ResponseSize: rw.size,
			ResponseTime: time.Since(start),
			UserAgent:    r.UserAgent(),
			Referer:      r.Referer(),
		}

		log.Println(FormatHTTPLog(entry))
	})
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
	size       int64
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	size, err := rw.ResponseWriter.Write(b)
	rw.size += int64(size)
	return size, err
}
