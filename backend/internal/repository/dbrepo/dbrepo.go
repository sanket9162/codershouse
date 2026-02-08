package dbrepo

import (
	"github.com/sanket9162/codershouse/internal/config"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type MonogRepo struct {
	App *config.Config
	DB  *mongo.Database
}

func NewMonogRepo(app *config.Config, db *mongo.Database) *MonogRepo {
	return &MonogRepo{
		App: app,
		DB:  db,
	}
}
