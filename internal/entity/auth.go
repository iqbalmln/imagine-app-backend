package entity

import (
	"time"

	jwt "github.com/golang-jwt/jwt/v4"
)

type Login struct {
	ID       int64  `json:"id" db:"id"`
	Username string `json:"username" db:"username"`
	Email    string `json:"email" db:"email"`
	Password string `json:"password" db:"password"`
	// Role      string    `json:"role" db:"role"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

var PIC int64

type MyClaims struct {
	jwt.RegisteredClaims
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Role     string `json:"role"`
}
