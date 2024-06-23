package middleware

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"
)

const (
	ContextLoggerStart        = "engine-logger-start"
	ContextLoggerOriginalPath = "engine-logger-path"
	ContextLoggerStatus       = "engine-logger-status"
)

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		ctx := context.WithValue(r.Context(), ContextLoggerStart, start)
		ctx = context.WithValue(ctx, ContextLoggerOriginalPath, r.URL.Path)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func LogRequest(r *http.Request, status int) {
	ctx := r.Context()
	path, ok := ctx.Value(ContextLoggerOriginalPath).(string)
	if !ok {
		// route was called without logger middleware enabled
		log.Println("unable to type cast log start context value")
		return
	}
	start, ok := ctx.Value(ContextLoggerStart).(time.Time)
	if !ok {
		log.Println("unable to type cast log start context value")
	}
	s := fmt.Sprintf("%s\t%d\t%s\t%v", r.Method, status, path, time.Since(start))
	log.Println(s)
}
