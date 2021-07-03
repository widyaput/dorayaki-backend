package models

import (
	"fmt"
	"net/http"
	"os"

	"github.com/dgrijalva/jwt-go"
)

var JwtKEY = os.Getenv("JWTKEY")
var Users = map[string]string{
	"Admin": os.Getenv("PWD_ADMIN"),
}

type Credentials struct {
	Password string `json:"password"`
	Username string `json:"username"`
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
