package http

import (
	"encoding/json"
	"net/http"

	"github.com/chrikar/chatheon/application"
	"github.com/chrikar/chatheon/application/ports"
)

type UserHandler struct {
	userService ports.UserService
}

func NewUserHandler(userService ports.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

type registerRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type loginResponse struct {
	Token string `json:"token"`
}

func (h *UserHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var req registerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	// Input validation first!
	if req.Username == "" || req.Password == "" {
		http.Error(w, "username and password cannot be empty", http.StatusBadRequest)
		return
	}

	err := h.userService.Register(req.Username, req.Password)
	if err != nil {
		if err == application.ErrUsernameTaken {
			http.Error(w, "username is already taken", http.StatusBadRequest)
		} else {
			http.Error(w, "failed to register user", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *UserHandler) LoginUser(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	token, err := h.userService.Login(req.Username, req.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(loginResponse{Token: token})
	if err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}
