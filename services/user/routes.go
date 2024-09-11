package user

import (
	"encoding/json"
	"net/http"

	"github.com/JigmeTenzinChogyel/go-net-http-server/database/generated"
)

type Handler struct {
	db *generated.Queries
}

func NewHandler(db *generated.Queries) *Handler {
	return &Handler{db: db}
}

func (h *Handler) RegisterRoutes(router *http.ServeMux) {
	router.HandleFunc("GET /hello", h.handleTest)
	router.HandleFunc("POST /login", h.handleLogin)
	router.HandleFunc("POST /register", h.handleRegister)
}

func (h *Handler) handleTest(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("User ID: "))
}

func (h *Handler) handleLogin(w http.ResponseWriter, r *http.Request) {

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
	existingUser, err := h.db.GetUserByEmail(r.Context(), input.Email)
	if err != nil {
		http.Error(w, "Failed to check existing user", http.StatusInternalServerError)
		return
	}

	if existingUser.Email == input.Email {
		http.Error(w, "User already exists", http.StatusConflict)
		return
	}

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
