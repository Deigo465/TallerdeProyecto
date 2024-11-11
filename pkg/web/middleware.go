package web

import (
	"fmt"
	"log/slog"
	"net/http"
)

// responseWriter is a minimal wrapper for http.ResponseWriter that allows you
// to capture the status code for logging purposes.
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

// WriteHeader captures the status code for later logging and calls the
// underlying WriteHeader method of the http.ResponseWriter.
func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

const (
	black   = "\033[1;30m"
	red     = "\033[1;31m"
	green   = "\033[1;32m"
	yellow  = "\033[1;33m"
	blue    = "\033[1;34m"
	magenta = "\033[1;35m"
	cyan    = "\033[1;36m"
	white   = "\033[1;37m"
	reset   = "\033[0m"
)

func ColorMsg(msg any, color string) string {
	return fmt.Sprintf("%s%v%s", color, msg, reset)
}

// add a middleware to log all requests like "[GET] /endpoint 200"
// make the method have colors
func LogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Wrap the original ResponseWriter
		wrappedWriter := &responseWriter{ResponseWriter: w}
		wrappedWriter.statusCode = http.StatusOK

		// Call the next handler with the wrapped ResponseWriter
		next.ServeHTTP(wrappedWriter, r)

		// var color string
		// switch {
		// case wrappedWriter.statusCode >= 200 && wrappedWriter.statusCode < 300:
		// 	color = blue
		// case wrappedWriter.statusCode >= 300 && wrappedWriter.statusCode < 400:
		// 	color = cyan
		// case wrappedWriter.statusCode >= 400 && wrappedWriter.statusCode < 500:
		// 	color = yellow
		// case wrappedWriter.statusCode >= 500:
		// 	color = red
		// }
		// Log the method, URL, and status code
		slog.Info("Request",
			"Status Code", wrappedWriter.statusCode,
			"Method", r.Method,
			"URL", r.URL.Path,
		)
	})
}

func NewHandler() {
}

// templateEngine.AddFunc("getYear", utils.GetYear)
// templateEngine.AddFunc("formatDate", utils.FormatDate)
// templateEngine.AddFunc("formatDatetime", utils.FormatDateWithTime)
// templateEngine.AddFunc("round2", utils.Round2)
// templateEngine.AddFunc("inc", utils.Inc)
// templateEngine.AddFunc("dec", utils.Dec)
// templateEngine.AddFunc("max", utils.Max)
// templateEngine.AddFunc("min", utils.Min)
// templateEngine.AddFunc("mul", utils.Mul)
// templateEngine.AddFunc("mul", utils.Mul)
// templateEngine.AddFunc("div", utils.Div)
// templateEngine.AddFunc("Base64Image", utils.Div)
// templateEngine.AddFunc("uintToString", utils.UintToString)
// // templateEngine.AddFunc("Base64", utils.Div)
// templateEngine.AddFunc("makeArray", utils.MakeArray)
// templateEngine.AddFunc("cmpPointerUint", func(a, b *uint) bool { return *a == *b })
// templateEngine.AddFunc("rmlTypeToInt", func(a enums.RMLType) int { return int(a) })
// templateEngine.AddFunc("Deref", func(a *interface{}) interface{} { return *a })
// templateEngine.AddFunc("isBool", func(a interface{}) bool {
// 	_, ok := a.(bool)
// 	return ok
// })
// templateEngine.AddFunc("Deref", func(s *uint) uint { return *s })
// templateEngine.AddFunc("formatDateMonthYear", utils.FormatDateMonthYear)
