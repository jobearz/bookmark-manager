package handler

import (
	"html/template"
	"net/http"

	"github.com/jobearz/bookmark-manager/internal/middleware"
	"github.com/jobearz/bookmark-manager/internal/models"
	"github.com/jobearz/bookmark-manager/internal/store"
)

type BookmarkHandler struct {
	store store.LinkStore
	tmpl  *template.Template
}

func NewBookmarkHandler(s store.LinkStore, tmpl *template.Template) *BookmarkHandler {
	return &BookmarkHandler{store: s, tmpl: tmpl}
}

func (h *BookmarkHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromToken(r)

	var bookmark models.Bookmark
	bookmark.UserID = userID
	bookmark.Title = r.FormValue("title")
	bookmark.URL = r.FormValue("url")
	bookmark.Notes = r.FormValue("notes")
	bookmark.CategoryID = r.FormValue("category_id")

	_, err := h.store.Create(bookmark)
	if err != nil {
		http.Error(w, "failed to create bookmark", http.StatusInternalServerError)
		return
	}

	// return updated bookmark list as HTML for HTMX to swap in
	bookmarks, err := h.store.GetAll(userID)
	if err != nil {
		http.Error(w, "failed to get bookmarks", http.StatusInternalServerError)
		return
	}
	h.tmpl.ExecuteTemplate(w, "bookmark-list", bookmarks)
}

func (h *BookmarkHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromToken(r)
	bookmarks, err := h.store.GetAll(userID)
	if err != nil {
		http.Error(w, "failed to get bookmarks", http.StatusInternalServerError)
		return
	}
	h.tmpl.ExecuteTemplate(w, "bookmark-list", bookmarks)
}

func (h *BookmarkHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/bookmarks/"):]
	bookmark, err := h.store.GetByID(id)
	if err != nil {
		http.Error(w, "bookmark not found", http.StatusNotFound)
		return
	}
	h.tmpl.ExecuteTemplate(w, "bookmark-card", bookmark)
}

func (h *BookmarkHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/bookmarks/"):]
	err := h.store.DeleteBookmark(id)
	if err != nil {
		http.Error(w, "failed to delete bookmark", http.StatusInternalServerError)
		return
	}
	// return empty string — HTMX will swap the card with nothing (deletes it)
	w.WriteHeader(http.StatusOK)
}
