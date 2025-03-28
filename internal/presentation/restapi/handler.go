package restapi

import (
	"gravitum/internal/service"
	"log/slog"
	"net/http"
)

type Handler struct {
	logger  *slog.Logger
	service *service.Service
}

func NewHandler(l *slog.Logger, s *service.Service) *Handler {
	return &Handler{
		logger:  l,
		service: s,
	}
}

func (h *Handler) Routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /users", h.CreateUser)
	mux.HandleFunc("GET /users/{id}", h.GetUser)
	mux.HandleFunc("PUT /users/{id}", h.UpdateUser)
	mux.HandleFunc("DELETE /users/{id}", h.DeleteUser)

	return WithRecover(mux)
}
