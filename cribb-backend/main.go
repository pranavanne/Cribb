package main

import (
	"cribb-backend/config"
	"log"
)

func main() {
	// Connect to MongoDB
	config.ConnectDB()

	log.Println("Successfully connected to MongoDB!")
}
