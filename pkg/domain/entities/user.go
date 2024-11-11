package entities

import (
	"fmt"
	"math/rand"
)

type User struct {
	ID             int      `json:"id"`
	Email          string   `json:"email"`
	Password       string   `json:"password"`
	HealthCenterId int      `json:"health_center_id"`
	ProfileId      int      `json:"profile_id"`
	Profile        Profile  `json:"profile"`
	Session        *Session `json:"session"`
}

func NewUser(id int, email, password string, healthCenterId int, profileId int, profile Profile) User {
	{
		return User{id, email, password, healthCenterId, profileId, profile, nil}
	}

}
func NewFakeUser() User {
	randomEmail := fmt.Sprintf("randomEmail%d@gmail.com", rand.Intn(1000))
	return User{1, randomEmail, "123456", 1, 1, NewFakeProfile(), nil}

}
