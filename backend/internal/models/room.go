package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Room struct {
	ID         bson.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	OwnerID    string        `bson:"OwnerID" json:"ownerId" validate:"required"`
	Owner      *User         `bson:"-" json:"owner,omitempty"`
	Topic      string        `bson:"topic" json:"topic" validate:"required"`
	RoomType   string        `bson:"roomType" json:"roomType" validate:"required"`
	SpeakerIDs []string      `bson:"speakers" json:"-"`
	Speakers   []*User       `bson:"-" json:"speakers"`
	CreatedAt  time.Time     `bson:"createdAt" json:"createdAt"`
	UpdatedAt  time.Time     `bson:"updatedAt" json:"updatedAt"`
}
