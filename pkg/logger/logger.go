package logger

import (
	"log"
	"net/http"
	"time"
)

type Logger struct{}

func New() *Logger {
	return &Logger{}
}

func (l *Logger) Debug(v ...interface{}) {
	log.Println(append([]interface{}{"[DEBUG]"}, v...)...)
}

func (l *Logger) Info(v ...interface{}) {
	log.Println(append([]interface{}{"[INFO]"}, v...)...)
}

func (l *Logger) Warn(v ...interface{}) {
	log.Println(append([]interface{}{"[WARN]"}, v...)...)
}

func (l *Logger) Error(v ...interface{}) {
	log.Println(append([]interface{}{"[ERROR]"}, v...)...)
}

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (l *Logger) RequestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sr := &statusRecorder{ResponseWriter: w, status: 200}
		start := time.Now()
		next.ServeHTTP(sr, r)
		log.Printf("%s %s %d %s", r.Method, r.URL.Path, sr.status, time.Since(start))
	})
}
