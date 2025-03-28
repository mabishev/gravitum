package restapi

import (
	"encoding/json"
	"errors"
	"gravitum/internal/domain/user"
	"gravitum/internal/service"
	"net/http"
)

type CreateUserRequest struct {
	Login     string `json:"login"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type CreateUserResponse struct {
	UserID string `json:"user_id"`
}

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		h.logger.Error("decode json", "error", err)
		return
	}

	result, err := h.service.CreateUser(r.Context(), service.CreateUserRequest{
		Login:     req.Login,
		FirstName: req.FirstName,
		LastName:  req.LastName,
	})

	if err != nil {
		switch {
		case errors.Is(err, user.ErrID), errors.Is(err, user.ErrLogin), errors.Is(err, user.ErrFirstName), errors.Is(err, user.ErrLastName):
			http.Error(w, err.Error(), http.StatusBadRequest)
		default:
			http.Error(w, "something went wrong", http.StatusInternalServerError)
		}

		h.logger.Error("create user", "error", err)
		return
	}

	resp := CreateUserResponse{
		UserID: result.ID.String(),
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "encode json error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

type User struct {
	ID        string `json:"id"`
	Login     string `json:"login"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
	u, err := h.service.GetUserByID(r.Context(), r.PathValue("id"))
	if errors.Is(err, user.ErrNotFound) {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "something went wrong", http.StatusInternalServerError)
		h.logger.Error("get user", "error", err)
		return
	}

	resp := User{
		ID:        u.ID().String(),
		Login:     u.Login(),
		FirstName: u.FirsName(),
		LastName:  u.LastName(),
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "encode json error", http.StatusInternalServerError)
		return
	}
}

type UpdateUserRequest struct {
	Login     string `json:"login"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

func (h *Handler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	var req UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	err := h.service.UpdateUser(r.Context(), service.UpdateUserRequest{
		ID:        r.PathValue("id"),
		Login:     req.Login,
		FirstName: req.FirstName,
		LastName:  req.LastName,
	})
	if errors.Is(err, user.ErrNotFound) {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "something went wrong", http.StatusInternalServerError)
		h.logger.Error("update user", "error", err)
		return
	}
}

func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	err := h.service.DeleteUserByID(r.Context(), r.PathValue("id"))
	if err != nil {
		http.Error(w, "something went wrong", http.StatusInternalServerError)
		h.logger.Error("delete user", "error", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
