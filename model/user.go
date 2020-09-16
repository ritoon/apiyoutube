package model

import "encoding/json"

type User struct {
	UUID      string `json:"uuid" gorm:"primaryKey"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Pass      string `json:"pass"`
}

func (u User) MarshalJSON() ([]byte, error) {
	aux := struct {
		UUID      string `json:"uuid"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Email     string `json:"email"`
	}{
		UUID:      u.UUID,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Email:     u.Email,
	}
	return json.Marshal(aux)
}

type Login struct {
	Email string `json:"email"`
	Pass  string `json:"pass"`
}
