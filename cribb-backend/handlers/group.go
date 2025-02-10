// handlers/group.go
package handlers

import (
	"context"
	"cribb-backend/config"
	"cribb-backend/models"
	"encoding/json"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// CreateGroupHandler creates a new group
// CreateGroupHandler creates a new group
func CreateGroupHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var group models.Group
	if err := json.NewDecoder(r.Body).Decode(&group); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Set timestamps
	group.CreatedAt = time.Now()
	group.UpdatedAt = time.Now()

	// Insert the new group into the database
	_, err := config.DB.Collection("groups").InsertOne(context.Background(), group)
	if err != nil {
		http.Error(w, "Failed to create group", http.StatusInternalServerError)
		return
	}

	// Return the created group
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(group)
}

// JoinGroupHandler allows a user to join a group
// JoinGroupHandler allows a user to join a group
func JoinGroupHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var request struct {
		Username  string `json:"username"`
		GroupName string `json:"group_name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Fetch the group by name
	var group models.Group
	err := config.DB.Collection("groups").FindOne(context.Background(), bson.M{"name": request.GroupName}).Decode(&group)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			http.Error(w, "Group not found", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to fetch group", http.StatusInternalServerError)
		}
		return
	}

	// Fetch the user by username
	var user models.User
	err = config.DB.Collection("users").FindOne(context.Background(), bson.M{"username": request.Username}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			http.Error(w, "User not found", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to fetch user", http.StatusInternalServerError)
		}
		return
	}

	// Update the user's group and group_id
	_, err = config.DB.Collection("users").UpdateOne(
		context.Background(),
		bson.M{"_id": user.ID},
		bson.M{"$set": bson.M{"group": group.Name, "group_id": group.ID}},
	)
	if err != nil {
		http.Error(w, "Failed to update user group", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// GetGroupMembersHandler retrieves all members of a group
// GetGroupMembersHandler retrieves all members of a group by group name
func GetGroupMembersHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	groupName := r.URL.Query().Get("group_name")
	if groupName == "" {
		http.Error(w, "Group name is required", http.StatusBadRequest)
		return
	}

	// Fetch the group by name
	var group models.Group
	err := config.DB.Collection("groups").FindOne(context.Background(), bson.M{"name": groupName}).Decode(&group)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			http.Error(w, "Group not found", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to fetch group", http.StatusInternalServerError)
		}
		return
	}

	// Fetch all users in the group
	cursor, err := config.DB.Collection("users").Find(context.Background(), bson.M{"group_id": group.ID})
	if err != nil {
		http.Error(w, "Failed to fetch users", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(context.Background())

	var users []models.User
	if err := cursor.All(context.Background(), &users); err != nil {
		http.Error(w, "Failed to decode users", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}
