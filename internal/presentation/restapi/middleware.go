package restapi

import (
	"context"
	"log/slog"
	"net/http"
)

func WithRecover(handler http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				slog.LogAttrs(context.Background(), slog.LevelError, "panic", slog.Any("cause", err))
			}
		}()

		handler.ServeHTTP(w, r)
	}
}
