package main

import (
	"dorayaki/configs"
	"dorayaki/configs/database"
	"dorayaki/internal/models"
	"log"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
)

// const defaultPort = "8080"
// POST /api/v1/orders

func main() {
	if err := database.ConnectDB(); err != nil {
		log.Fatal("Error when connect to database")
	}
	log.Print("Success connect to database")
	database.DB.AutoMigrate(&models.TokoDorayaki{})
	database.DB.AutoMigrate(&models.Dorayaki{})
	database.DB.AutoMigrate(&models.Toko{})
	// TODO:handlers.
	r := chi.NewRouter()
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: configs.AllowedOrigins,
	}))

}
