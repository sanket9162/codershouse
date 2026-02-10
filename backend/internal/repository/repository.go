package repository

import "github.com/sanket9162/codershouse/internal/models"

type DatabaseRepo interface {
	CreateUser(u *models.User) error
}
