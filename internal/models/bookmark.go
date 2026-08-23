package models

import (
	"time"
)

type Bookmark struct {
	ID         string    `json:"id"`
	UserID     string    `json:"user_id"`
	CategoryID string    `json:"category_id"`
	Title      string    `json:"title"`
	URL        string    `json:"url"`
	Notes      string    `json:"notes"`
	CreatedAt  time.Time `json:"created_at"`
}
