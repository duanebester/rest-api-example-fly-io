package models

import (
	// "strconv"

	// "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       uuid.UUID `json:"id" gorm:"type:uuid;primary_key"`
	Email    string    `json:"email" gorm:"type:varchar(100);unique_index"`
	Password string    `json:"-" gorm:"type:varchar(100)"`
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

type NewUserInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func NewUser(opts *NewUserInput) *User {
	hashedPw, err := hashPassword(opts.Password)
	if err != nil {
		panic(err)
	}
	return &User{
		ID:       uuid.New(),
		Email:    opts.Email,
		Password: hashedPw,
	}
}

type LoginInput struct {
	Identity string `json:"identity"`
	Password string `json:"password"`
}
