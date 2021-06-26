package main

import (
	"dorayaki/configs/database"
	"dorayaki/internal/models"
	"log"
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
}
