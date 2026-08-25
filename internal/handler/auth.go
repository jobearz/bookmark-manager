package handler

import (
	"html/template"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jobearz/bookmark-manager/config"
	"github.com/jobearz/bookmark-manager/internal/models"
	"github.com/jobearz/bookmark-manager/internal/store"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	store store.LinkStore
	tmpl  *template.Template
}

func NewAuthHandler(s store.LinkStore, tmpl *template.Template) *AuthHandler {
	return &AuthHandler{store: s, tmpl: tmpl}
}

func (h *AuthHandler) LoginPage(w http.ResponseWriter, r *http.Request) {
	h.tmpl.ExecuteTemplate(w, "login", nil)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	email := r.FormValue("email")
	password := r.FormValue("password")

	user, err := h.store.GetUserByEmail(email)
	if err != nil {
		h.tmpl.ExecuteTemplate(w, "login", map[string]string{
			"Error": "Invalid email or password",
		})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		h.tmpl.ExecuteTemplate(w, "login", map[string]string{
			"Error": "Invalid email or password",
		})
		return
	}

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

	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    tokenString,
		HttpOnly: true,
		Path:     "/",
		MaxAge:   7 * 24 * 60 * 60,
	})

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *AuthHandler) RegisterPage(w http.ResponseWriter, r *http.Request) {
	h.tmpl.ExecuteTemplate(w, "register", nil)
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/register", http.StatusSeeOther)
		return
	}

	email := r.FormValue("email")
	password := r.FormValue("password")

	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "failed to hash password", http.StatusInternalServerError)
		return
	}

	user := models.User{
		Email:    email,
		Password: string(hashedBytes),
	}

	_, err = h.store.CreateUser(user)
	if err != nil {
		h.tmpl.ExecuteTemplate(w, "register", map[string]string{
			"Error": "Failed to create account. Email may already be in use.",
		})
		return
	}

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:   "token",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
