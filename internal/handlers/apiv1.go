package handlers

import (
	"dorayaki/internal/models"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
)

func apiv1(router chi.Router) {
	router.Post("/signin", signIn)
	fs := http.FileServer(http.Dir(uploadPath))
	router.Handle("/files/{nameOfFile}", http.StripPrefix("/api/v1/files/", fs))
	router.Route("/shops", shops)
	router.Route("/dorayakis", dorayakis)
}

func signIn(w http.ResponseWriter, r *http.Request) {
	var cred models.Credentials
	if err := render.Bind(r, &cred); err != nil {
		render.Render(w, r, models.ErrBadRequest)
		return
	}
	expectedPassword, ok := models.Users[cred.Username]
	if !ok {
		render.Render(w, r, models.ErrUnauthorized)
		return
	}
	if expectedPassword != cred.Password {
		if expectedPassword != "" {
			render.Render(w, r, models.ErrUnauthorized)
			return
		}
		// for dev only
		if cred.Password != "dorayakidev" {
			render.Render(w, r, models.ErrUnauthorized)
			return
		}
	}
	expirationTime := time.Now().Add(1 * time.Hour)
	claims := &models.JWT{
		Username: cred.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	sign := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), claims)
	token, err := sign.SignedString([]byte(models.JwtKEY))
	if err != nil {
		render.Render(w, r, models.ServerErrorRenderer(err))
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    token,
		Expires:  expirationTime,
		HttpOnly: true,
	})
	resp := models.ResponseString{Response: *models.SuccessResponse}
	resp.Data = append(resp.Data, token)
	if err := render.Render(w, r, &resp); err != nil {
		render.Render(w, r, models.ServerErrorRenderer(err))
		return
	}
}

func auth(r *http.Request, fn ...func(r *http.Request) string) error {
	token := findTokens(r, fn...)
	if token == "" {
		return fmt.Errorf("no token found")
	}
	return parseToken(token)
}

func findTokens(r *http.Request, fn ...func(r *http.Request) string) string {
	var token string
	for _, funct := range fn {
		token = funct(r)
		if token != "" {
			break
		}
	}
	return token
}

func parseToken(tokenString string) error {
	token, err := jwt.ParseWithClaims(tokenString, &models.JWT{}, func(token *jwt.Token) (interface{}, error) {
		if jwt.GetSigningMethod("HS256") != token.Method {
			return nil, fmt.Errorf("unexpected signing method : %v", token.Header["alg"])
		}
		return []byte(models.JwtKEY), nil
	})

	if err != nil {
		return err
	}
	if token == nil || !token.Valid {
		return fmt.Errorf("unauthorized")
	}
	return nil
}

func tokenFromHeader(r *http.Request) string {
	bearer := r.Header.Get("Authorization")
	if len(bearer) > 7 && strings.ToUpper(bearer[0:6]) == "BEARER" {
		return bearer[7:]
	}
	return ""
}

func tokenFromCookie(r *http.Request) string {
	c, err := r.Cookie("token")
	if err != nil {
		return ""
	}
	return c.Value
}
func tokenFromQuery(r *http.Request) string {
	return r.URL.Query().Get("jwt")
}
