package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog"
)

type LoggerKey string

const loggerKey = LoggerKey("reqid")

//  UseLogger takes a parent logger and injects a request specific child logger into
// the request context
func UseLogger(logger *zerolog.Logger) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			reqId := GetRequestId(r)
			ctx := GetContextWithLogger(r.Context(), logger, reqId)
			l := GetLoggerFromCtx(ctx)
			l.Info().Str("method", r.Method).Str("url", r.URL.String()).Str("remote_address", r.RemoteAddr).Str("user_agent", r.UserAgent()).Msg("")
			start := time.Now()
			wrapped := wrapResponseWriter(w)
			next.ServeHTTP(wrapped, r.WithContext(ctx))
			l.Info().Int("status", wrapped.status).Str("method", r.Method).Str("path", r.URL.EscapedPath()).Str("duration", time.Since(start).String()).Msg("")
		})
	}
}

type responseWriter struct {
	http.ResponseWriter
	status      int
	wroteHeader bool
}

func wrapResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{ResponseWriter: w}
}

func (rw *responseWriter) Status() int {
	return rw.status
}

func (rw *responseWriter) WriteHeader(code int) {
	if rw.wroteHeader {
		return
	}
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
	rw.wroteHeader = true
}

func GetContextWithLogger(ctx context.Context, logger *zerolog.Logger, requestId string) context.Context {
	l := logger.With().Str("requestId", requestId).Logger()
	return context.WithValue(ctx, loggerKey, &l)
}

func GetLogger(r *http.Request) *zerolog.Logger {
	return GetLoggerFromCtx(r.Context())
}

func GetLoggerFromCtx(ctx context.Context) *zerolog.Logger {
	l := ctx.Value(loggerKey)
	// TODO: should we return default logger or panic, default logger could hide bugs, panic seems bad?
	if l == nil {
		panic("context does not have a logger, panic")
	}
	return l.(*zerolog.Logger)
}
