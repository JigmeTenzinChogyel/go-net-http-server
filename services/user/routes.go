package user

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/JigmeTenzinChogyel/go-net-http-server/database/generated"
	"github.com/JigmeTenzinChogyel/go-net-http-server/middleware"
	"github.com/JigmeTenzinChogyel/go-net-http-server/types"
	"github.com/JigmeTenzinChogyel/go-net-http-server/utils"
	"golang.org/x/crypto/bcrypt"
)

type Handler struct {
	db *generated.Queries
}

func NewHandler(db *generated.Queries) *Handler {
	return &Handler{db: db}
}

func (h *Handler) RegisterRoutes(router *http.ServeMux) {
	router.HandleFunc("GET /todos", h.handleTodos)
	router.HandleFunc("POST /login", h.handleLogin)
	router.HandleFunc("POST /register", h.handleRegister)
	router.HandleFunc("POST /todo", h.handleCreateTodo)

}

func (h *Handler) handleTodos(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id, ok := ctx.Value(middleware.UserKey).(int32)
	if !ok {
		http.Error(w, "Invalid user ID in context", http.StatusInternalServerError)
		return
	}
	todos, err := h.db.ListTodos(ctx, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return the token in a JSON response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todos)
}

func (h *Handler) handleLogin(w http.ResponseWriter, r *http.Request) {
	var input types.LoginUser

	// Parse JSON body
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	user, err := h.db.GetUserByEmailWithPass(r.Context(), input.Email)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get user: %v", err.Error()), http.StatusInternalServerError)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		http.Error(w, fmt.Sprintf("wrong password: %s", err.Error()), http.StatusBadRequest)
		return
	}

	token, err := utils.CreateToken(user.ID)
	if err != nil {
		http.Error(w, fmt.Sprintf("error generating token: %s", err.Error()), http.StatusBadRequest)
		return
	}

	// Return the token in a JSON response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": token})

}

// register or create user
func (h *Handler) handleRegister(w http.ResponseWriter, r *http.Request) {
	var input generated.CreateUserParams

	// Parse JSON body
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Check if user already exists
	existingUser, err := h.db.CheckUserExists(r.Context(), input.Email)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to check existing user: %v", err.Error()), http.StatusInternalServerError)
		return
	}

	if existingUser {
		http.Error(w, "User already exists", http.StatusConflict)
		return
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}

	input.Password = string(passwordHash)

	// Create user using sqlc generated function
	user, err := h.db.CreateUser(r.Context(), input)
	if err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	// Return the created user
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (h *Handler) handleCreateTodo(w http.ResponseWriter, r *http.Request) {
	var input types.CreateTodoInput

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	id, ok := ctx.Value(middleware.UserKey).(int32)
	if !ok {
		http.Error(w, "Invalid user ID in context", http.StatusInternalServerError)
		return
	}

	newTodo := generated.CreateTodoParams{
		UserID:      id,
		Title:       input.Title,
		Description: sql.NullString{String: input.Description, Valid: input.Description != ""},
		Completed:   input.Completed,
	}

	todo, err := h.db.CreateTodo(ctx, newTodo)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to create todo: %v", err.Error()), http.StatusInternalServerError)
		return
	}

	// Return the created todo
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todo)
}
