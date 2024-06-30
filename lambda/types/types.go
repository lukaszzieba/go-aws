package types

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type RegisterUser struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type User struct {
	Username     string `json:"username"`
	PasswordHash string `json:"password"`
}

type LoginUser struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func NewUser(registerUser RegisterUser) (User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(registerUser.Password), 10)

	if err != nil {
		return User{}, err
	}

	return User{Username: registerUser.Username, PasswordHash: string(hashedPassword)}, nil
}

func ValidatePassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))

	return err == nil
}

func CreateToken(username string) string {
	now := time.Now()

	validUntil := now.Add(time.Hour + 1).Unix()

	claims := jwt.MapClaims{
		"user":    username,
		"expires": validUntil,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims, nil)
	secret := "secret"

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return ""
	}

	return tokenString
}
