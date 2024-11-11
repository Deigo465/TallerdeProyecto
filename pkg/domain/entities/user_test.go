package entities

import "testing"

func TestNewUser(t *testing.T) {
	//given
	email := "New Email"
	password := "New Password"
	healthCenterId := 1
	profileUser := NewFakeProfile()
	//when
	user := NewUser(1, email, password, healthCenterId, profileUser.ID, profileUser)

	//then
	if user.Email != email {
		t.Fatalf("Expecting email to be %s, got %s", email, user.Email)
	}
	if user.Password != password {
		t.Fatalf("Expecting password to be %s, got %s", password, user.Password)
	}
	if user.HealthCenterId != healthCenterId {
		t.Fatalf("Expecting healthCenter id to be %d, got %d", healthCenterId, user.HealthCenterId)
	}
	if user.ProfileId != profileUser.ID {
		t.Fatalf("Expecting profile id to be %d, got %d", profileUser.ID, user.ProfileId)
	}
}

func TestNewFakeUser(t *testing.T) {

	user := NewFakeUser()

	if user.Email == "" {
		t.Fatalf("Expecting email to not be empty")
	}
	if user.Password == "" {
		t.Fatalf("Expecting password to not be empty")
	}
	if user.HealthCenterId == 0 {
		t.Fatalf("Expecting healthCenter id to not be empty")
	}
	if user.ProfileId == 0 {
		t.Fatalf("Expecting profile id to not be empty")
	}
}
