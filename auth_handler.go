package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

type AuthHandler struct {
	users *UserStore
	cfg   *Config
}

func NewAuthHandler(users *UserStore, cfg *Config) *AuthHandler {
	return &AuthHandler{users: users, cfg: cfg}
}

func (h *AuthHandler) register(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		badRequest("invalid JSON body").respond(w)
		return
	}

	// validate input
	errs := make(map[string]string)
	if strings.TrimSpace(input.Username) == "" {
		errs["username"] = "username is required"
	}
	if len(input.Password) < 6 {
		errs["password"] = "password must be at least 6 characters"
	}
	if len(errs) > 0 {
		validationError(errs).respond(w)
		return
	}

	user, err := h.users.Create(input.Username, input.Password)
	if err != nil {
		// username already taken
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			badRequest("username already taken").respond(w)
			return
		}
		internalError("failed to create user").respond(w)
		return
	}

	writeJSON(w, http.StatusCreated, user)
}

func (h *AuthHandler) login(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		badRequest("invalid JSON body").respond(w)
		return
	}

	// find user
	user, ok, err := h.users.GetByUsername(input.Username)
	if err != nil {
		internalError("failed to find user").respond(w)
		return
	}
	if !ok {
		// use same message for not found and wrong password
		// never reveal which one failed — security best practice
		badRequest("invalid credentials").respond(w)
		return
	}

	// check password
	if err := h.users.CheckPassword(user, input.Password); err != nil {
		badRequest("invalid credentials").respond(w)
		return
	}

	// generate JWT token
	token, err := GenerateToken(user)
	if err != nil {
		internalError("failed to generate token").respond(w)
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{
		"token": token,
	})
}
