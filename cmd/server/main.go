package main

import (
	"dorayaki/configs/database"
	"dorayaki/internal/handlers"
	"dorayaki/internal/models"
	"log"
	"net"
	"net/http"
	"os"
)

const defaultPort = "8080"

// POST /api/v1/orders

func main() {
	if err := database.ConnectDB(); err != nil {
		log.Fatal("Error when connect to database")
	}
	log.Print("Success connect to database")
	database.DB.SetupJoinTable(&models.Toko{}, "Dorayaki", &models.TokoDorayaki{})
	database.DB.AutoMigrate(&models.Toko{})
	database.DB.AutoMigrate(&models.Dorayaki{})

	database.DB.Migrator().CreateConstraint(&models.TokoDorayaki{}, "TokoID")
	database.DB.Migrator().CreateConstraint(&models.TokoDorayaki{}, "DorayakiID")

	r := handlers.NewHandler()
	server := &http.Server{
		Handler: r,
	}

	port := os.Getenv("PORT")

	if port == "" {
		port = defaultPort
	}
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("Error occured: %s", err.Error())
	}
	server.Serve(listener)
}
