package config

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var DB *mongo.Database

// ConnectDB initializes MongoDB connection and sets up the database
func ConnectDB() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file:", err)
	}

	// Get and validate environment variables
	mongoURI := strings.TrimSpace(os.Getenv("MONGODB_URI"))
	dbName := strings.TrimSpace(os.Getenv("DB_NAME"))

	if mongoURI == "" {
		log.Fatal("MONGODB_URI is required in .env file")
	}

	if dbName == "" {
		log.Fatal("DB_NAME is required in .env file")
	}

	log.Printf("Attempting to connect to MongoDB...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Connect to MongoDB
	clientOptions := options.Client().ApplyURI(mongoURI)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}

	// Ping the database to verify connection
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal("Failed to ping MongoDB:", err)
	}

	DB = client.Database(dbName)

	// Initialize database collections and indexes
	if err := initializeDatabase(); err != nil {
		log.Fatal("Failed to initialize database:", err)
	}

	log.Printf("Successfully connected to MongoDB database: %s", dbName)
}

func initializeDatabase() error {
	if DB == nil {
		return fmt.Errorf("database connection not initialized")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	log.Println("Creating users collection and indexes...")

	// Create users collection with indexes
	usersCollection := DB.Collection("users")

	// Create indexes for users collection
	indexModels := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "username", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys:    bson.D{{Key: "phone_number", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.D{{Key: "score", Value: -1}},
		},
	}

	// Create the indexes
	_, err := usersCollection.Indexes().CreateMany(ctx, indexModels)
	if err != nil {
		return fmt.Errorf("failed to create indexes: %v", err)
	}

	log.Println("Successfully initialized database collections and indexes")
	return nil
}
