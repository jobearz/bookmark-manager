package store

import (
	"sync"

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
