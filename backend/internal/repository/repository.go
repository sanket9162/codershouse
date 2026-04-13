package repository

import "github.com/sanket9162/codershouse/internal/models"

type DatabaseRepo interface {
	CreateUser(u *models.User) error
	GetUserByPhone(phone string) (*models.User, error)
	GetUserByID(id string) (*models.User, error)
	UpdateUser(u *models.User) error
	CreateRoom(r *models.Room) error
	GetAllRooms(roomType, searchQuery string) ([]models.Room, error)
	GetRoomById(id string) (*models.Room, error)
}
