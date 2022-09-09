package utils

import "net/http"

// ref: https://gist.github.com/Boerworz/b683e46ae0761056a636
type LoggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func NewLoggingResponseWriter(w http.ResponseWriter) *LoggingResponseWriter {
	// WriteHeader(int) is not called if our response implicitly returns 200 OK, so
	// we default to that status code.
	return &LoggingResponseWriter{w, http.StatusOK}
}

func (lrw *LoggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func (lrw *LoggingResponseWriter) StatusCode() int {
	return lrw.statusCode
}
