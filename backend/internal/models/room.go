package models

import "time"

type Room struct {
	ID        string    `bson:"_id,omitempty" json:"_id,omitempty"`
	UserID    string    `bson:"userID" json:"userID" validate:"required"`
	Topic     string    `bson:"topic" json:"topic" validate:"required"`
	RoomType  string    `bson:"roomType" json:"roomType" validate:"required"`
	CreatedAt time.Time `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time `bson:"updatedAt" json:"updatedAt"`
}
