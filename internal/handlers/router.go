package handlers

import (
	"dorayaki/configs"
	"dorayaki/internal/models"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
)

func NewHandler() http.Handler {
	router := chi.NewRouter()
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins: configs.AllowedOrigins,
	}))

	fs := http.FileServer(http.Dir(uploadPath))
	router.Handle("/files/*", http.StripPrefix("/files/", fs))
	router.Route("/api/v1/", apiv1)
	router.MethodNotAllowed(methodNotAllowedHandler)
	router.NotFound(notFoundHandler)
	return router
}

func methodNotAllowedHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	render.Render(w, r, models.ErrMethodNotAllowed)
}
func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	render.Render(w, r, models.ErrNotFound)
}
