package main

import (
	"cribb-backend/config"
	"cribb-backend/handlers"
	"log"
	"net/http"
)

func main() {
	// Connect to MongoDB and initialize collections
	config.ConnectDB()

	// Register routes
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Server is running!"))
	})

	// Auth routes
	http.HandleFunc("/api/register", handlers.RegisterHandler)

	// User routes
	http.HandleFunc("/api/users", handlers.GetUsersHandler)
	http.HandleFunc("/api/users/by-username", handlers.GetUserByUsernameHandler)
	http.HandleFunc("/api/users/by-score", handlers.GetUsersByScoreHandler)

	log.Println("Server starting on port 8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
