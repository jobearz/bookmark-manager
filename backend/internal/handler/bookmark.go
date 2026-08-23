package handler

import (
	"encoding/json"
	"net/http"

	"github.com/jobearz/bookmark-manager/internal/middleware"
	"github.com/jobearz/bookmark-manager/internal/models"
	"github.com/jobearz/bookmark-manager/internal/store"
)

type BookmarkHandler struct {
	store store.LinkStore
}

func NewBookmarkHandler(s store.LinkStore) *BookmarkHandler {
	return &BookmarkHandler{store: s}
}

func (h *BookmarkHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromToken(r)
	// check if request is POST
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var bookmark models.Bookmark
	bookmark.UserID = userID
	if err := json.NewDecoder(r.Body).Decode(&bookmark); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	created, err := h.store.Create(bookmark)
	if err != nil {
		http.Error(w, "failed to create bookmark", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(created)
}

func (h *BookmarkHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromToken(r)
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	bookmarks, err := h.store.GetAll(userID)
	if err != nil {
		http.Error(w, "failed to get bookmark", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(bookmarks)
}

func (h *BookmarkHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/bookmarks/"):]
	bookmark, err := h.store.GetByID(id)
	if err != nil {
		http.Error(w, "failed to find a bookmark with that ID", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(bookmark)
}

func (h *BookmarkHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/bookmarks/"):]
	err := h.store.DeleteBookmark(id)
	if err != nil {
		http.Error(w, "failed to delete bookmark", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
