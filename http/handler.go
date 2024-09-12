package http

import (
	"encoding/json"
	"net/http"

	"lmich.com/tononkira/domain"
)

type Handler struct {
	UserService domain.UserService
}

// CreateUser handles HTTP requests for creating a user
func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user domain.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	if err := h.UserService.CreateUser(&user); err != nil {
		http.Error(w, "Error creating user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
