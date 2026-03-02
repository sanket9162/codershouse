package models

import "time"

type Room struct {
	ID        string    `bson:"_id,omitempty" json:"_id,omitempty"`
	OwnerID   string    `bson:"OwnerID" json:"OwnerID" validate:"required"`
	Topic     string    `bson:"topic" json:"topic" validate:"required"`
	RoomType  string    `bson:"roomType" json:"roomType" validate:"required"`
	Speakers  []string  `bson:"speakers" json:"speakers"`
	CreatedAt time.Time `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time `bson:"updatedAt" json:"updatedAt"`
}
