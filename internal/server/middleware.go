package server

import (
	"bytes"
	"fmt"
	"github.com/go-easy-templ/internal/helper"
	"github.com/go-easy-templ/internal/logger"
	"github.com/google/uuid"
	"log/slog"
	"net/http"
)

type Middleware func(http.Handler) http.Handler

type wrappedWriter struct {
	http.ResponseWriter
	statusCode int
	body       bytes.Buffer
}

func (w *wrappedWriter) WriteHeader(statusCode int) {
	w.ResponseWriter.WriteHeader(statusCode)
	w.statusCode = statusCode
}

func (srv *Server) addMiddlewares(mx ...Middleware) Middleware {
	return func(next http.Handler) http.Handler {
		for i := len(mx) - 1; i >= 0; i-- {
			x := mx[i]
			next = x(next)
		}
		return next
	}
}

func (srv *Server) requestInterceptor() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			wrappedWriter := &wrappedWriter{
				ResponseWriter: w,
				statusCode:     0,
			}

			ctx := logger.AppendCtx(r.Context(), slog.String(ContextKeyRequestID, r.Header.Get(HeaderRequestID)))

			srv.logger.DebugContext(ctx, "incoming request", slog.Group("request_details",
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.Any("request_body", r.Body),
			))

			r = r.WithContext(ctx)

			next.ServeHTTP(wrappedWriter, r)

		})
	}
}

func (srv *Server) correlationID() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			if r.Header.Get(HeaderRequestID) == "" {
				id := uuid.New()
				r.Header.Set(HeaderRequestID, id.String())
				w.Header()[HeaderRequestID] = []string{id.String()}
			}

			next.ServeHTTP(w, r)
		})
	}
}

func (srv *Server) recovery() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			defer func() {
				if err := recover(); err != nil {
					w.Header().Set("Connection", "close")
					srv.logger.Error("Recovered from panic", fmt.Errorf("%s", err))
					helper.ServerErrorResponse(srv.logger, w, r, fmt.Errorf("%s", err))
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}
