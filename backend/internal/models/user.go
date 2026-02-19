package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type User struct {
	ID        bson.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Phone     string        `json:"phone" bson:"phone"`
	Activated bool          `json:"activated" bson:"activated"`
	CreatedAt time.Time     `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time     `json:"updated_at" bson:"updated_at"`
	Name      string        `json:"name" bson:"name"`
	Avatar    string        `json:"avatar" bson:"avatar"`
}
