package entity

import "time"

type Shot struct {
	ID          int64     `json:"id" db:"id"`
	Title       string    `json:"title" db:"title"`
	IMG         string    `json:"img" db:"img"`
	Description string    `json:"description" db:"description"`
	Category    string    `json:"category" db:"category"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
	DeletedAt   time.Time `json:"deleted_at" db:"deleted_at"`
}
