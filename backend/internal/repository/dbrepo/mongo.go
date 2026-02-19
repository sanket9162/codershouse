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
