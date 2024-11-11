package usecase_test

import (
	"testing"

	"github.com/open-wm/blockehr/pkg/domain/entities"
	usecase "github.com/open-wm/blockehr/pkg/domain/usecases"
	moc "github.com/open-wm/blockehr/pkg/mocks"
)

func TestSaveUser(t *testing.T) {
	//GIVEN
	userRepo := moc.NewInMemoryUserRepository()
	uc := usecase.NewUserUsecase(userRepo)

	//staff
	staffProfile := entities.NewFakeProfile()
	staffProfile.Role = entities.STAFF
	staffActor := entities.NewUser(1, "italo@blockehr.pe", "password", 1, staffProfile.ID, staffProfile)

	doctorProfile := entities.NewFakeProfile()
	doctorProfile.Role = entities.DOCTOR
	//doctor
	doctorProfile.Role = entities.DOCTOR
	doctorActor := entities.NewUser(1, "italo@blockehr.pe", "password", 1, doctorProfile.ID, doctorProfile)
	//user
	newUser := entities.NewFakeUser()
	//invalid
	invUserNoEmail := entities.NewUser(1, "", "password", 1, doctorProfile.ID, doctorProfile)
	invUserNoPassword := entities.NewUser(1, "italo@blockehr.pe", "", 1, doctorProfile.ID, doctorProfile)
	invUserNoHealthCenterId := entities.NewUser(1, "italo@blockehr.pe", "password", 0, doctorProfile.ID, doctorProfile)
	invUserNoProfileId := entities.NewUser(1, "italo@blockehr.pe", "password", 1, 0, doctorProfile)
	//invUserNoProfile := entities.NewUser(1, "italo@blockehr.pe", "password", 1, doctorProfile.ID, entities.Profile{}) //TODO

	testCases := []struct {
		name    string
		actor   entities.User
		user    *entities.User
		want    error
		wantErr bool
	}{
		{
			name:    "valid user with staff",
			actor:   staffActor,
			user:    &newUser,
			wantErr: false,
			want:    nil,
		},
		{
			name:    "invalid user - user made by non-staff",
			actor:   doctorActor,
			user:    &newUser,
			wantErr: true,
			want:    nil,
		},
		{
			name:    "invalid user - email is empty",
			actor:   staffActor,
			user:    &invUserNoEmail,
			wantErr: true,
			want:    nil,
		},
		{
			name:    "invalid user - password is empty",
			actor:   staffActor,
			user:    &invUserNoPassword,
			wantErr: true,
			want:    nil,
		},
		{
			name:    "invalid user - health center is empty",
			actor:   staffActor,
			user:    &invUserNoHealthCenterId,
			wantErr: true,
			want:    nil,
		},
		{
			name:    "invalid user - profile id is empty",
			actor:   staffActor,
			user:    &invUserNoProfileId,
			wantErr: true,
			want:    nil,
		},
		/*
			{
				name:    "invalid user - profile is empty",
				actor:   staffActor,
				user:    &invUserNoProfile,
				wantErr: true,
				want:    nil,
			},
		*/
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			//WHEN
			err := uc.SaveUser(&tc.actor, tc.user)

			// THEN
			if tc.wantErr {
				if err == nil {
					t.Errorf("Error adding user: %v", err)
				}
			} else {
				if err != nil {
					t.Errorf("Error adding user: %v", err)
				}
				_, err := userRepo.GetUserByID(tc.user.ID)

				if err != nil {
					t.Errorf("user not added to repository")
				}
			}
		})
	}

}

func TestResetDoctorPassword(t *testing.T) {

	//GIVEN
	userRepo := moc.NewInMemoryUserRepository()
	uc := usecase.NewUserUsecase(userRepo)

	//staff
	staffProfile := entities.NewFakeProfile()
	staffProfile.Role = entities.STAFF
	staffActor := entities.NewUser(1, "staff@blockehr.pe", "password", 1, staffProfile.ID, staffProfile)

	doctorProfile := entities.NewFakeProfile()
	doctorProfile.ID = 4
	doctorProfile.Role = entities.DOCTOR
	//doctor
	doctorProfile.Role = entities.DOCTOR
	doctorActor := entities.NewUser(4, "italo@blockehr.pe", "12345678", 1, doctorProfile.ID, doctorProfile)

	testCases := []struct {
		name    string
		actor   entities.User
		user    *entities.User
		want    error
		wantErr bool
	}{
		{
			name:    "valid user with staff",
			actor:   staffActor,
			user:    &doctorActor,
			wantErr: false,
			want:    nil,
		},
		{
			name:    "invalid user with doctor",
			actor:   doctorActor,
			user:    &doctorActor,
			wantErr: true,
			want:    nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			//WHEN
			_, err := uc.ResetDoctorPassword(&tc.actor, tc.user.ID)

			// THEN
			if tc.wantErr {
				if err == nil {
					t.Errorf("Error reset password : %v", err)
				}
			} else {
				if err != nil {
					t.Errorf("Error reset password: %v", err)
				}
			}
		})
	}

}
