package model

import "encoding/json"

type User struct {
	UUID      string `json:"uuid"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Pass      string `json:"pass"`
}

func (u User) MarshalJSON() ([]byte, error) {
	aux := struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Email     string `json:"email"`
	}{
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Email:     u.Email,
	}
	return json.Marshal(aux)
}
