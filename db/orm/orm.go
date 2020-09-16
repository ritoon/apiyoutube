package orm

import (
	"fmt"
	"log"

	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"apiyoutube/model"
)

// ORMDB is use to simulate a db connection.
type ORMDB struct {
	conn *gorm.DB
}

// New is creating a user list.
func New(host, user, pass, name, port string) *ORMDB {
	var db ORMDB

	dsn := fmt.Sprintf("host=%v user=%v password=%v dbname=%v port=%v sslmode=disable TimeZone=Europe/Paris", host, user, pass, name, port)
	var err error
	db.conn, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Panicf("db/orm: try to connect to the db %v", err)
	}

	db.conn.AutoMigrate(&model.User{})

	db.initUser()

	return &db
}

func (db ORMDB) initUser() {
	u := model.User{
		FirstName: "Rob",
		LastName:  "Pike",
	}

	db.AddUser(&u)

	u2 := model.User{
		FirstName: "Marie",
		LastName:  "Curie",
	}

	db.AddUser(&u2)
}

// AddUser is adding a new user in the ORMDB.
// this send an error if the user allready exsits in the ORMDB.
func (db *ORMDB) AddUser(u *model.User) error {
	u.UUID = uuid.New().String()
	return db.conn.Create(u).Error
}

// UpdateUser is updating a user in the ORMDB.
// this send an error if the user allready exsits in the ORMDB.
func (db *ORMDB) UpdateUser(uuid string, u model.User) error {
	return db.conn.Table("users").Where("uuid = ?", uuid).Updates(u).Error
}

// GetUser retrives form the ORMDB a given uuid.
// this send an error if the don't exsits in the ORMDB.
func (db *ORMDB) GetUser(uuid string) (*model.User, error) {
	var u model.User
	return &u, db.conn.Table("users").Where("uuid = ?", uuid).First(&u).Error
}

// DeleteUser is deleting the given uuid user from the userList.
// if the user don't exists in the userList this function sends an error.
func (db *ORMDB) DeleteUser(uuid string) error {
	return db.conn.Table("users").Where("uuid = ?", uuid).Delete(&model.User{}).Error
}

// GetListUser retrive all users from the db.
func (db *ORMDB) GetListUser() (map[string]*model.User, error) {
	users := map[string]*model.User{}
	usersDB, error := db.conn.Limit(-1).Model(&model.User{}).Rows()
	if error != nil {
		return nil, error
	}
	for usersDB.Next() {
		var user model.User
		error = db.conn.ScanRows(usersDB, &user)
		users[user.UUID] = &user
	}
	return users, nil
}

func (db *ORMDB) GetUserByEmail(email string) (*model.User, error) {
	var u model.User
	return &u, db.conn.Table("users").Where("email = ?", email).First(&u).Error
}
