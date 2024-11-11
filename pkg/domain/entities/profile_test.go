package entities

import (
	"testing"
)

func TestNewProfile(t *testing.T) {
	//given
	firstName := "New First Name"
	motherLastName := "New Mother Last Name"
	fatherLastName := "New Father Last Name"
	documentNumber := "New Document Number"
	phone := "123456789"
	dateOfBirth := "2006-01-02"
	cmp := "987654321"
	gender := "Masculino"
	role := DOCTOR
	specialty := "Psicología"
	email := "email@example.com"
	//when
	profile := NewProfile(1, firstName, motherLastName, fatherLastName,
		documentNumber, gender, phone, email, dateOfBirth, cmp, specialty, role)

	//then
	if profile.FirstName != firstName {
		t.Fatalf("Expecting first name to be %s, got %s", firstName, profile.FirstName)
	}

	if profile.MotherLastName != motherLastName {
		t.Fatalf("Expecting mother last name to be %s, got %s", motherLastName, profile.MotherLastName)
	}

	if profile.FatherLastName != fatherLastName {
		t.Fatalf("Expecting father last name to be %s, got %s", fatherLastName, profile.FatherLastName)
	}
	if profile.DocumentNumber != documentNumber {
		t.Fatalf("Expecting document number to be %s, got %s", documentNumber, profile.DocumentNumber)
	}
	if profile.Gender != gender {
		t.Fatalf("Expecting gender to be %s, got %s", gender, profile.Gender)
	}
	if profile.Phone != phone {
		t.Fatalf("Expecting phone to be %s, got %s", phone, profile.Phone)
	}
	if profile.DateOfBirth != dateOfBirth {
		t.Fatalf("Expecting date of birth to be %s, got %s", dateOfBirth, profile.DateOfBirth)
	}
	if profile.Cmp != cmp {
		t.Fatalf("Expecting cmp to be %s, got %s", cmp, profile.Cmp)
	}
	if profile.Specialty != specialty {
		t.Fatalf("Expecting cmp to be %s, got %s", specialty, profile.Specialty)
	}
	if profile.Role != role {
		t.Fatalf("Expecting role to be %s, got %s", role, profile.Role)
	}
}

func TestNewFakeProfile(t *testing.T) {
	profile := NewFakeProfile()

	if profile.FirstName == "" {
		t.Fatalf("Expecting first name to not be empty")
	}
	if profile.MotherLastName == "" {
		t.Fatalf("Expecting mother last name to not be empty")
	}
	if profile.FatherLastName == "" {
		t.Fatalf("Expecting father last name to not be empty")
	}
	if profile.DocumentNumber == "" {
		t.Fatalf("Expecting document number to not be empty")
	}
	if profile.Phone == "" {
		t.Fatalf("Expecting phone to not be empty")
	}
	if profile.DateOfBirth == "" {
		t.Fatalf("Expecting date of birth name to not be empty")
	}
	if profile.Cmp == "" {
		t.Fatalf("Expecting cmp to not be empty")
	}
	if profile.Specialty == "" {
		t.Fatalf("Expecting specialty to not be empty")
	}
	if profile.Role == " " {
		t.Fatalf("Expecting role to not be empty")
	}

}
