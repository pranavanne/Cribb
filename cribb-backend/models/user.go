package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Username    string             `bson:"username" json:"username"`
	Password    string             `bson:"password" json:"-"` // "-" ensures password isn't included in JSON responses
	Name        string             `bson:"name" json:"name"`
	PhoneNumber string             `bson:"phone_number" json:"phone_number"`
	Score       int                `bson:"score" json:"score"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updated_at"`
}
