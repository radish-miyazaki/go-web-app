package entity

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

type UserID int64

type User struct {
	ID        UserID    `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	Password  string    `json:"password" db:"password"`
	Role      string    `json:"role" db:"role"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" db:"updated_at"`
}

func (u *User) ComparePassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
}
