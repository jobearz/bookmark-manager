package store

import "github.com/jobearz/bookmark-manager/internal/models"

type MockStore struct{}

func (m *MockStore) Create(bookmark models.Bookmark) (models.Bookmark, error) {
	bookmark.ID = "test-id"
	return bookmark, nil
}

func (m *MockStore) GetAll(userID string) ([]models.Bookmark, error) {
	return []models.Bookmark{}, nil
}

func (m *MockStore) GetByID(id string) (models.Bookmark, error) {
	return models.Bookmark{ID: id}, nil
}

func (m *MockStore) DeleteBookmark(id string) error {
	return nil
}

func (m *MockStore) EditBookmark(id string, bookmark models.Bookmark) (models.Bookmark, error) {
	return models.Bookmark{ID: id}, nil
}

func (m *MockStore) CreateCategory(category models.Category) (models.Category, error) {
	category.ID = "test-id"
	return category, nil
}

func (m *MockStore) DeleteCategory(id string) error {
	return nil
}

func (m *MockStore) EditCategory(id string, category models.Category) (models.Category, error) {
	return models.Category{ID: id}, nil
}

func (m *MockStore) CreateUser(user models.User) (models.User, error) {
	user.ID = "test-id"
	return user, nil
}

func (m *MockStore) GetUserByEmail(email string) (models.User, error) {
	return models.User{ID: "test-id", Email: email}, nil
}

func (m *MockStore) GetUserByID(id string) (models.User, error) {
	return models.User{ID: id}, nil
}
