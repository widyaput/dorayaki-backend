package handlers

import (
	"github.com/go-chi/chi"
)

func apiv1(router chi.Router) {
	// router.Route("/assets", assets)
	router.Route("/shops", shops)
	router.Route("/dorayakis", dorayakis)
}
