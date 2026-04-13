package dbrepo

import (
	"context"
	"time"

	"github.com/sanket9162/codershouse/internal/models"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func (m *MongoRepo) CreateUser(u *models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.DB.Collection("users").InsertOne(ctx, u)
	if err != nil {
		return err
	}

	return nil
}

func (m *MongoRepo) GetUserByPhone(phone string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var user models.User
	filter := bson.M{"phone": phone}

	err := m.DB.Collection("users").FindOne(ctx, filter).Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (m *MongoRepo) GetUserByID(id string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var user models.User
	oid, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	filter := bson.M{"_id": oid}
	err = m.DB.Collection("users").FindOne(ctx, filter).Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (m *MongoRepo) UpdateUser(u *models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	filter := bson.M{"_id": u.ID}
	update := bson.M{
		"$set": bson.M{
			"name":       u.Name,
			"avatar":     u.Avatar,
			"activated":  u.Activated,
			"updated_at": time.Now(),
		},
	}

	_, err := m.DB.Collection("users").UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	return nil
}

func (m *MongoRepo) CreateRoom(r *models.Room) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	res, err := m.DB.Collection("rooms").InsertOne(ctx, r)
	if err != nil {
		return err
	}

	// Extract the generated ObjectID and assign it to the room struct
	if oid, ok := res.InsertedID.(bson.ObjectID); ok {
		r.ID = oid
	}

	return nil
}

func (m *MongoRepo) GetAllRooms(roomType, searchQuery string) ([]models.Room, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rooms := []models.Room{}
	filter := bson.M{}

	if roomType != "all" && roomType != "" {
		filter["roomType"] = roomType
	}

	if searchQuery != "" {
		filter["topic"] = bson.M{"$regex": searchQuery, "$options": "i"}
	}

	cursor, err := m.DB.Collection("rooms").Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	err = cursor.All(ctx, &rooms)
	if err != nil {
		return nil, err
	}

	// Fetch speakers and owner for each room manually
	for i, room := range rooms {
		// Populate Owner
		owner, err := m.GetUserByID(room.OwnerID)
		if err == nil && owner != nil {
			rooms[i].Owner = owner
		}

		// Populate Speakers
		speakers := []*models.User{}
		for _, speakerID := range room.SpeakerIDs {
			speaker, err := m.GetUserByID(speakerID)
			if err == nil && speaker != nil {
				speakers = append(speakers, speaker)
			}
		}
		rooms[i].Speakers = speakers
	}

	return rooms, nil
}

func (m *MongoRepo) GetRoomById(id string) (*models.Room, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var room models.Room
	oid, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	filter := bson.M{"_id": oid}
	err = m.DB.Collection("rooms").FindOne(ctx, filter).Decode(&room)
	if err != nil {
		return nil, err
	}

	return &room, nil
}
