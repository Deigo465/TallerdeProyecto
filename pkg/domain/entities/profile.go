package entities

import (
	"fmt"
	"math/rand"
)

type Role string

const (
	DOCTOR  Role = "DOCTOR"
	STAFF   Role = "STAFF"
	ADMIN   Role = "ADMIN"
	PATIENT Role = "PATIENT"
)

type Profile struct {
	ID             int    `json:"id"`
	FirstName      string `json:"first_name"`
	MotherLastName string `json:"mother_last_name"`
	FatherLastName string `json:"father_last_name"`
	DocumentNumber string `json:"document_number"`
	Gender         string `json:"gender"`
	Phone          string `json:"phone"`
	ContactEmail   string `json:"contact_email"`
	DateOfBirth    string `json:"date_of_birth"`
	Cmp            string `json:"cmp"`
	Specialty      string `json:"specialty"`
	Role           Role   `json:"role"`
}

func NewProfile(id int, firstName, motherLastName, fatherLastName, documentNumber, gender, phone, email string, dateOfBirth string, cmp, specialty string, role Role) Profile {
	return Profile{id, firstName, motherLastName, fatherLastName, documentNumber, gender,
		phone, email, dateOfBirth, cmp, specialty, role}

}

func NewFakeProfile() Profile {
	randomName := fmt.Sprintf("B. %d", rand.Intn(1000))
	return Profile{1, randomName, "randon", "name",
		"72549855", "Masculino", "974528438", "random@example.com", "12-29-2004", "1234567", "psicologia", DOCTOR}

}
