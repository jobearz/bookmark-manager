package store

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jobearz/bookmark-manager/internal/models"
	_ "github.com/lib/pq"
)

type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStore(db *sql.DB) *PostgresStore {
	return &PostgresStore{db: db}
}

func (s *PostgresStore) Create(bookmark models.Bookmark) (models.Bookmark, error) {
	bookmark.ID = uuid.New().String()
	bookmark.CreatedAt = time.Now()

	_, err := s.db.Exec(
		"INSERT INTO bookmarks (id, user_id, category_id, title, url, notes, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7)",
		bookmark.ID, bookmark.UserID, bookmark.CategoryID, bookmark.Title, bookmark.URL, bookmark.Notes, bookmark.CreatedAt,
	)
	if err != nil {
		return models.Bookmark{}, err
	}

	return bookmark, nil
}

func (s *PostgresStore) GetAll(userID string) ([]models.Bookmark, error) {
	rows, err := s.db.Query("SELECT id, user_id, category_id, title, url, notes, created_at FROM bookmarks WHERE user_id = $1", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bookmarks []models.Bookmark
	for rows.Next() {
		var bookmark models.Bookmark
		err := rows.Scan(&bookmark.ID, &bookmark.UserID, &bookmark.CategoryID, &bookmark.Title, &bookmark.URL, &bookmark.Notes, &bookmark.CreatedAt)
		if err != nil {
			return nil, err
		}
		bookmarks = append(bookmarks, bookmark)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return bookmarks, nil
}

func (s *PostgresStore) GetByID(id string) (models.Bookmark, error) {
	row := s.db.QueryRow("SELECT id, user_id, category_id, title, url, notes, created_at FROM bookmarks WHERE id = $1", id)
	var bookmark models.Bookmark
	err := row.Scan(&bookmark.ID, &bookmark.UserID, &bookmark.CategoryID, &bookmark.Title, &bookmark.URL, &bookmark.Notes, &bookmark.CreatedAt)
	if err != nil {
		return models.Bookmark{}, err
	}
	return bookmark, nil
}

func (s *PostgresStore) EditBookmark(id string, bookmark models.Bookmark) (models.Bookmark, error) {
	_, err := s.db.Exec(
		"UPDATE bookmarks SET title = COALESCE(NULLIF($1, ''), title), url = COALESCE(NULLIF($2, ''), url), notes = COALESCE(NULLIF($3, ''), notes), category_id = COALESCE(NULLIF($4, ''), category_id) WHERE id = $5",
		bookmark.Title, bookmark.URL, bookmark.Notes, bookmark.CategoryID, id,
	)
	if err != nil {
		return models.Bookmark{}, err
	}
	return s.GetByID(id)
}

func (s *PostgresStore) DeleteBookmark(id string) error {
	_, err := s.db.Exec("DELETE FROM bookmarks WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("failed to delete bookmark: %w", err)
	}
	return nil
}

func (s *PostgresStore) CreateCategory(category models.Category) (models.Category, error) {
	category.ID = uuid.New().String()
	category.CreatedAt = time.Now()

	_, err := s.db.Exec(
		"INSERT INTO categories (id, user_id, name, created_at) VALUES ($1, $2, $3, $4)",
		category.ID, category.UserID, category.Name, category.CreatedAt,
	)
	if err != nil {
		return models.Category{}, err
	}

	return category, nil
}

func (s *PostgresStore) EditCategory(id string, category models.Category) (models.Category, error) {
	row := s.db.QueryRow(
		"UPDATE categories SET name = COALESCE(NULLIF($1, ''), name) WHERE id = $2 RETURNING id, user_id, name, created_at",
		category.Name, id,
	)
	var updated models.Category
	err := row.Scan(&updated.ID, &updated.UserID, &updated.Name, &updated.CreatedAt)
	if err != nil {
		return models.Category{}, err
	}
	return updated, nil
}

func (s *PostgresStore) DeleteCategory(id string) error {
	_, err := s.db.Exec("DELETE FROM categories WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("failed to delete category: %w", err)
	}
	return nil
}

func (s *PostgresStore) CreateUser(user models.User) (models.User, error) {
	user.ID = uuid.New().String()
	user.CreatedAt = time.Now()

	_, err := s.db.Exec(
		"INSERT INTO users (id, email, password, created_at) VALUES ($1, $2, $3, $4)",
		user.ID, user.Email, user.Password, user.CreatedAt,
	)
	if err != nil {
		return models.User{}, fmt.Errorf("insert failed: %w", err)
	}
	return user, nil
}

func (s *PostgresStore) GetUserByEmail(email string) (models.User, error) {
	row := s.db.QueryRow("SELECT id, email, password, created_at FROM users WHERE email = $1", email)
	var user models.User
	err := row.Scan(&user.ID, &user.Email, &user.Password, &user.CreatedAt)
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}

func (s *PostgresStore) GetUserByID(id string) (models.User, error) {
	row := s.db.QueryRow("SELECT id, email, created_at FROM users WHERE id = $1", id)
	var user models.User
	err := row.Scan(&user.ID, &user.Email, &user.CreatedAt)
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}
