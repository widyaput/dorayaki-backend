package models

import (
	"fmt"
	"net/http"
	"os"

	jwt "github.com/golang-jwt/jwt/v4"
)

var JwtKEY = os.Getenv("JWTKEY")

type Credentials struct {
	Username string `json:"username" gorm:"primaryKey"`
	Password string `json:"password" gorm:"not null"`
}

func (Credentials) TableName() string {
	return "credentials"
}

type JWT struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

func InitialiseKey() {
	if JwtKEY == "" {
		JwtKEY = "smokeWeedEveryday"
	}
}

func (c *Credentials) Bind(r *http.Request) error {
	if c.Password == "" {
		return fmt.Errorf("password is required")
	}
	if c.Username == "" {
		return fmt.Errorf("username is required")
	}
	return nil
}
