package store

import (
	"fmt"
	"sync"
	"time"
	"uuid"

	"github.com/jobearz/bookmark-manager/internal/models"
)

type MemoryStore struct {
	bookmarks  map[string]models.Bookmark
	categories map[string]models.Category
	users      map[string]models.User
	mu         sync.RWMutex
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		bookmarks:  make(map[string]models.Bookmark),
		categories: make(map[string]models.Category),
		users:      make(map[string]models.User),
	}
}

func (s *MemoryStore) Create(bookmark models.Bookmark) (models.Bookmark, error) {
	// mutex lock before writing to the map
	s.mu.Lock()
	defer s.mu.Unlock()

	bookmark.ID = uuid.New().String()
	bookmark.CreatedAt = time.Now()
	s.bookmarks[bookmark.ID] = bookmark
	return bookmark, nil
}

func (s *MemoryStore) GetAll(userID string) ([]models.Bookmark, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	bookmarks := make([]models.Bookmark, 0, len(s.bookmarks))
	for _, bookmark := range s.bookmarks {
		if bookmark.UserID == userID {
			bookmarks = append(bookmarks, bookmark)
		}
	}
	return bookmarks, nil
}

func (s *MemoryStore) GetByID(id string) (models.Bookmark, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	bookmark, ok := s.bookmarks[id]
	if !ok {
		return models.Bookmark{}, fmt.Errorf("Bookmark with id %s not found", id)
	}
	return bookmark, nil
}

func (s *MemoryStore) CreateCategory(category models.Category) (models.Category, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	category.ID = uuid.New().String()
	category.CreatedAt = time.Now()
	s.categories[category.ID] = category
	return category, nil
}

func (s *MemoryStore) CreateUser(user models.User) (models.User, error) {
	s.mu.Unlock()
	defer s.mu.Unlock()

	user.ID = uuid.New().String()
	user.CreatedAt = time.Now()
	s.users[user.ID] = user
	return user, nil
}

func (s *MemoryStore) GetUserByEmail(email string) (models.User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, user := range s.users {
		if user.Email == email {
			return user, nil
		}
	}
	return models.User{}, fmt.Errorf("user with email %s not found", email)
}

func (s *MemoryStore) DeleteBookmark(id string) error { return nil }

func (s *MemoryStore) DeleteCategory(id string) error { return nil }

func (s *MemoryStore) EditBookmark(id string, bookmark models.Bookmark) (models.Bookmark, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	existing, ok := s.bookmarks[id]
	if !ok {
		return models.Bookmark{}, fmt.Errorf("bookmark with id %s not found", id)
	}

	// keep existing value if new value is empty
	updated := models.Bookmark{
		ID:         existing.ID,
		UserID:     existing.UserID,
		CreatedAt:  existing.CreatedAt,
		Title:      coalesce(bookmark.Title, existing.Title),
		URL:        coalesce(bookmark.URL, existing.URL),
		Notes:      coalesce(bookmark.Notes, existing.Notes),
		CategoryID: coalesce(bookmark.CategoryID, existing.CategoryID),
	}

	s.bookmarks[id] = updated
	return updated, nil
}

func (s *MemoryStore) EditCategory(id string, category models.Category) (models.Category, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	existing, ok := s.categories[id]
	if !ok {
		return models.Category{}, fmt.Errorf("category with id %s not found", category)
	}

	updated := models.Category{
		ID:        existing.ID,
		UserID:    existing.UserID,
		CreatedAt: existing.CreatedAt,
		Name:      coalesce(category.Name, existing.Name),
	}

	s.categories[id] = updated
	return updated, nil
}

func coalesce(new, existing string) string {
	if new == "" {
		return existing
	}
	return new
}
