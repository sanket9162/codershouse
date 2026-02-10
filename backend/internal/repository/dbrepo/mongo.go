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
