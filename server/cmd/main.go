package main

import (
	"html/template"
	"log"
	"net/http"

	"github.com/jobearz/bookmark-manager/db"
	"github.com/jobearz/bookmark-manager/internal/handler"
	"github.com/jobearz/bookmark-manager/internal/middleware"
	"github.com/jobearz/bookmark-manager/internal/store"
)

func main() {
	database, err := db.Connect()
	if err != nil {
		log.Fatal("failed to connect to database:", err)
	}
	defer database.Close()

	// load all templates
	tmpl := template.Must(template.ParseGlob("templates/**/*.html"))
	tmpl = template.Must(tmpl.ParseGlob("templates/*.html"))

	// init store and handlers
	pgStore := store.NewPostgresStore(database)
	bookmarkHandler := handler.NewBookmarkHandler(pgStore, tmpl)
	authHandler := handler.NewAuthHandler(pgStore, tmpl)

	mux := http.NewServeMux()

	// static files
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// auth routes
	mux.HandleFunc("/login", authHandler.LoginPage)
	mux.HandleFunc("/auth/login", authHandler.Login)
	mux.HandleFunc("/register", authHandler.RegisterPage)
	mux.HandleFunc("/auth/register", authHandler.Register)
	mux.HandleFunc("/logout", authHandler.Logout)

	// protected routes
	mux.HandleFunc("/", middleware.RequireAuth(func(w http.ResponseWriter, r *http.Request) {
		tmpl.ExecuteTemplate(w, "index", nil)
	}))
	mux.HandleFunc("/bookmarks", middleware.RequireAuth(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			bookmarkHandler.GetAll(w, r)
		case http.MethodPost:
			bookmarkHandler.Create(w, r)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	}))
	mux.HandleFunc("/bookmarks/", middleware.RequireAuth(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			bookmarkHandler.GetByID(w, r)
		case http.MethodDelete:
			bookmarkHandler.Delete(w, r)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	}))

	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
