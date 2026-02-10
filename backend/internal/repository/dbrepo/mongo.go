package dbrepo

import (
	"context"
	"time"

	"github.com/sanket9162/codershouse/internal/models"
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
