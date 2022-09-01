package domain

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        int       `gorm:"primary_key" json:"id"`
	Name      string    `gorm:"size:255" json:"name"`
	Email     string    `gorm:"size:255" json:"email"`
	Password  string    `gorm:"size:255" json:"password"`
	NoHP      string    `gorm:"size:255" json:"no_hp"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewUser(name, email, password string) *User {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return &User{
		Name:     name,
		Email:    email,
		Password: string(hashedPassword),
	}
}

var PrivateKey = []byte(`MyKey`)

func (u *User) GenerateToken() (string, error) {
	claims := jwt.MapClaims{
		"user_id": u.ID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
		"iss":     "isana",
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(PrivateKey)
}
