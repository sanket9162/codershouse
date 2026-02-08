package dbrepo

import (
	"github.com/sanket9162/codershouse/internal/config"
	"github.com/sanket9162/codershouse/internal/repository"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type MongoRepo struct {
	App *config.Config
	DB  *mongo.Database
}

func NewMongoRepo(app *config.Config, db *mongo.Database) repository.DatabaseRepo {
	return &MongoRepo{
		App: app,
		DB:  db,
	}
}
