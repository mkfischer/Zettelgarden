package models

import "time"

type Card struct {
	ID         int           `json:"id"`
	CardID     string        `json:"card_id"`
	UserID     int           `json:"user_id"`
	Title      string        `json:"title"`
	Body       string        `json:"body"`
	Link       string        `json:"link"`
	IsDeleted  bool          `json:"is_deleted"`
	CreatedAt  time.Time     `json:"created_at"`
	UpdatedAt  time.Time     `json:"updated_at"`
	Parent     PartialCard   `json:"parent"`
	Files      []File        `json:"files"`
	Children   []PartialCard `json:"children"`
	References []PartialCard `json:"references"`
}

type PartialCard struct {
	ID        int       `json:"id"`
	CardID    string    `json:"card_id"`
	UserID    int       `json:"user_id"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type EditCardParams struct {
	CardID string `json:"card_id"`
	Title  string `json:"title"`
	Body   string `json:"body"`
	Link   string `json:"link"`
}
