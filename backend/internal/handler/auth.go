package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jobearz/bookmark-manager/config"
	"github.com/jobearz/bookmark-manager/internal/models"
	"github.com/jobearz/bookmark-manager/internal/store"
	"golang.org/x/crypto/bcrypt"
)

type AuthorizationHandler struct {
	store store.LinkStore
}

func NewAuthorizationHandler(s store.LinkStore) *AuthorizationHandler {
	return &AuthorizationHandler{store: s}
}

func (h *AuthorizationHandler) Register(w http.ResponseWriter, r *http.Request) {
	var UserRequestBody struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	// decode from request body
	if err := json.NewDecoder(r.Body).Decode(&UserRequestBody); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	email := UserRequestBody.Email
	password := UserRequestBody.Password

	// hash password using bcrypt
	hashPassword := func(password string) (string, error) {
		bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, "failed to hash password", http.StatusInternalServerError)
			return "", err
		}
		return string(bytes), err
	}

	hashedPassword, _ := hashPassword(password)

	user := models.User{
		Email:    email,
		Password: hashedPassword,
	}

	newUser, err := h.store.CreateUser(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	newUser.Email = email
	newUser.Password = hashedPassword

	// return 201 with user excluding password
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newUser)
}

func (h *AuthorizationHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	var UserRequestBody struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	// decode email + pass from body
	if err := json.NewDecoder(r.Body).Decode(&UserRequestBody); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	user, err := h.store.GetUserByEmail(UserRequestBody.Email)
	if err != nil {
		http.Error(w, "user not found", http.StatusUnauthorized)
		return
	}

	// compare pass w/ hash
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(UserRequestBody.Password)); err != nil {
		http.Error(w, "password does not match hashed password", http.StatusUnauthorized)
		return
	}

	// generate jwt token if they match
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"exp":     time.Now().Add(7 * 24 * time.Hour).Unix(),
	})

	tokenString, err := token.SignedString([]byte(config.JWTSecret()))
	if err != nil {
		http.Error(w, "failed to generate token", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"token":   tokenString,
		"user_id": user.ID,
	})
}
