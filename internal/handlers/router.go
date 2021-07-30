package handlers

import (
	"dorayaki/internal/models"
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

var CurrentJWT *models.JWT

func NewHandler() http.Handler {
	router := chi.NewRouter()
	router.Use(Cors)
	router.Use(middleware.Logger)

	fs := http.FileServer(http.Dir("swaggerui/"))
	router.Handle("/docs/api/v1/*", http.StripPrefix("/docs/api/v1/", fs))
	router.Route("/api/v1/", apiv1)
	router.MethodNotAllowed(methodNotAllowedHandler)
	router.NotFound(notFoundHandler)
	return router
}

func authenticator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := auth(r, tokenFromHeader, tokenFromCookie, tokenFromQuery); err != nil {
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

func Cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "X-Requested-With, Content-Type")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		// log.Printf("Should set headers")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}
