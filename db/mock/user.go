package mock

import (
	"errors"

	"github.com/google/uuid"

	"apiyoutube/model"
)

// MockDB is use to simulate a db connection.
type MockDB struct {
	listUser map[string]*model.User
}

// New is creating a user list.
func New() *MockDB {
	var db MockDB
	db.listUser = make(map[string]*model.User)
	db.initUser()
	return &db
}

func (db MockDB) initUser() {
	u := model.User{
		FirstName: "Rob",
		LastName:  "Pike",
	}

	db.AddUser(u)

	u2 := model.User{
		FirstName: "Marie",
		LastName:  "Curie",
	}

	db.AddUser(u2)
}

// AddUser is adding a new user in the MockDB.
// this send an error if the user allready exsits in the MockDB.
func (db *MockDB) AddUser(u model.User) error {
	id := uuid.New().String()
	u.UUID = id
	db.listUser[u.UUID] = &u
	return nil
}

// UpdateUser is updating a user in the MockDB.
// this send an error if the user allready exsits in the MockDB.
func (db *MockDB) UpdateUser(uuid string, u model.User) error {
	if _, ok := db.listUser[uuid]; !ok {
		return errors.New("user don't exists")
	}
	db.listUser[uuid] = &u
	return nil
}

// GetUser retrives form the MockDB a given uuid.
// this send an error if the don't exsits in the MockDB.
func (db *MockDB) GetUser(uuid string) (*model.User, error) {
	if _, ok := db.listUser[uuid]; !ok {
		return nil, errors.New("user don't exists")
	}
	return db.listUser[uuid], nil
}

// DeleteUser is deleting the given uuid user from the userList.
// if the user don't exists in the userList this function sends an error.
func (db *MockDB) DeleteUser(uuid string) error {
	if _, ok := db.listUser[uuid]; !ok {
		return errors.New("user don't exists")
	}
	delete(db.listUser, uuid)
	return nil
}
