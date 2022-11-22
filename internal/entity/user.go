package entity

import (
	"time"
)

type User struct {
	ID        uint64    `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	Email     string    `json:"email" db:"email"`
	Password  string    `json:"password" db:"password"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	DeletedAt time.Time `json:"deleted_at" db:"deleted_at"`
}
