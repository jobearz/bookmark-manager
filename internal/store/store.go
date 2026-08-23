package store

import "github.com/jobearz/bookmark-manager/internal/models"

type LinkStore interface {
	Create(bookmark models.Bookmark) (models.Bookmark, error)
	GetAll(userID string) ([]models.Bookmark, error)
	GetByID(id string) (models.Bookmark, error)
	CreateCategory(category models.Category) (models.Category, error)
	CreateUser(user models.User) (models.User, error)
	GetUserByEmail(email string) (models.User, error)
	GetUserByID(id string) (models.User, error)
	DeleteBookmark(id string) error
	DeleteCategory(id string) error
}
