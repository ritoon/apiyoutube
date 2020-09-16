package db

import (
	"apiyoutube/model"
)

type DB interface {
	DBUser
}

type DBUser interface {
	AddUser(u *model.User) error
	UpdateUser(uuid string, u model.User) error
	GetUser(uuid string) (*model.User, error)
	// GetListUser return all users.
	GetListUser() (map[string]*model.User, error)
	DeleteUser(uuid string) error
	GetUserByEmail(email string) (*model.User, error)
}
