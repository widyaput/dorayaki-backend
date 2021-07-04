package handlers

import (
	"dorayaki/configs"
	"dorayaki/internal/models"
	"fmt"
	"net/http"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
)

var CurrentJWT *models.JWT

func NewHandler() http.Handler {
	router := chi.NewRouter()
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins: configs.AllowedOrigins,
	}))
	router.Use(middleware.Logger)
	router.Post("/signin", signin)
	fs := http.FileServer(http.Dir(uploadPath))
	router.Handle("/files/*", http.StripPrefix("/files/", fs))
	router.Route("/api/v1/", apiv1)
	router.MethodNotAllowed(methodNotAllowedHandler)
	router.NotFound(notFoundHandler)
	return router
}

func Authenticator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := Auth(r, TokenFromHeader, TokenFromCookie, TokenFromQuery); err != nil {
			fmt.Println(err.Error())
			render.Render(w, r, models.ErrUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func methodNotAllowedHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	render.Render(w, r, models.ErrMethodNotAllowed)
}
func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	render.Render(w, r, models.ErrNotFound)
}

func signin(w http.ResponseWriter, r *http.Request) {
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
		Name:    "token",
		Value:   token,
		Expires: expirationTime,
	})
	resp := models.ResponseImageURL{Response: *models.SuccessResponse}
	resp.Data = append(resp.Data, token)
	if err := render.Render(w, r, &resp); err != nil {
		render.Render(w, r, models.ServerErrorRenderer(err))
		return
	}
}

func Auth(r *http.Request, fn ...func(r *http.Request) string) error {
	token := findTokens(r, fn...)
	if token == "" {
		return fmt.Errorf("no token found")
	}
	return ParseToken(token)
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

func ParseToken(tokenString string) error {
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

func TokenFromHeader(r *http.Request) string {
	bearer := r.Header.Get("Authorization")
	if len(bearer) > 7 && strings.ToUpper(bearer[0:6]) == "BEARER" {
		return bearer[7:]
	}
	return ""
}

func TokenFromCookie(r *http.Request) string {
	c, err := r.Cookie("token")
	if err != nil {
		return ""
	}
	return c.Value
}

func TokenFromQuery(r *http.Request) string {
	return r.URL.Query().Get("jwt")
}
