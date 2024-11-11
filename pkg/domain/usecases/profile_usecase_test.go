package usecase_test

import (
	"testing"

	"github.com/open-wm/blockehr/pkg/domain/entities"
	usecase "github.com/open-wm/blockehr/pkg/domain/usecases"
	"github.com/open-wm/blockehr/pkg/interfaces"
	moc "github.com/open-wm/blockehr/pkg/mocks"
)

func getProfileUC() (usecase.ProfileUsecase, interfaces.ProfileRepository) {
	profileRepo := moc.NewInMemoryProfileRepository()
	appointmentsRepo := moc.NewInMemoryAppointmentRepository()
	userRepo := moc.NewInMemoryUserRepository()
	uc := usecase.NewProfileUsecase(profileRepo, appointmentsRepo, userRepo)
	return uc, profileRepo
}

func TestAddProfile(t *testing.T) {
	//GIVEN
	uc, profileRepo := getProfileUC()

	//doctor
	doctorProfile := entities.NewFakeProfile()
	doctorProfile.Role = entities.DOCTOR
	doctorActor := entities.NewUser(1, "italo@blockehr.pe", "password", 1, doctorProfile.ID, doctorProfile)
	//staff
	staffProfile := entities.NewFakeProfile()
	staffProfile.Role = entities.STAFF
	staffActor := entities.NewUser(2, "italo@blockehr.pe", "password", 1, staffProfile.ID, staffProfile)
	//profile
	newProfile := entities.Profile{ID: 2, FirstName: "Jorge", FatherLastName: "De la Flor", MotherLastName: "Valle", DocumentNumber: "40058778", Gender: "Masculino", Phone: "94785987", DateOfBirth: "2024-03-01", Cmp: "1", Specialty: "Psicologia", Role: "DOCTOR"}
	//invalid

	invProfileNoFirstName := entities.NewProfile(1, "", "Kunimoto", "Kunimoto", "40058778", "Masculino", "94785987", "email", "2024-03-01", "5", "Psicologia", "DOCTOR")
	invProfileNoMotherLastName := entities.NewProfile(1, "Rodrigo", "", "Kunimoto", "40058778", "Masculino", "94785987", "email", "2024-03-01", "5", "Psicologia", "DOCTOR")
	invProfileNoFatherLastName := entities.NewProfile(1, "Rodrigo", "Kunimoto", "", "40058778", "Masculino", "94785987", "email", "2024-03-01", "5", "Psicologia", "DOCTOR")
	invProfileNoDocumentNumber := entities.NewProfile(1, "Rodrigo", "Kunimoto", "Kunimoto", "", "Masculino", "94785987", "email", "2024-03-01", "5", "Psicologia", "DOCTOR")
	invProfileNoPhone := entities.NewProfile(1, "Rodrigo", "Kunimoto", "Kunimoto", "40058778", "Masculino", "", "2024-03-01", "email", "5", "Psicologia", "DOCTOR")
	invProfileNoDateOfBirth := entities.NewProfile(1, "Rodrigo", "Kunimoto", "Kunimoto", "40058778", "Masculino", "94785987", "email", "", "5", "Psicologia", "DOCTOR")
	invProfileNoCmp := entities.NewProfile(1, "Rodrigo", "Kunimoto", "Kunimoto", "40058778", "Masculino", "94785987", "email", "2024-03-01", "", "Psicologia", "DOCTOR")
	invProfileNoSpecialty := entities.NewProfile(1, "Rodrigo", "Kunimoto", "Kunimoto", "40058778", "Masculino", "94785987", "email", "2024-03-01", "5", "", "DOCTOR")

	// invProfileNoGender := entities.NewProfile(1, "Rodrigo", "Kunimoto", "Kunimoto", "40058778", "", "94785987", "2024-03-01", "5", "Psicologia", "DOCTOR")

	testCases := []struct {
		name    string
		actor   entities.User
		profile *entities.Profile
		want    error
		wantErr bool
	}{
		{
			name:    "valid profile with staff",
			actor:   staffActor,
			profile: &newProfile,
			wantErr: false,
			want:    nil,
		},
		{
			name:    "invalid profile - profile made by  non-staff",
			actor:   doctorActor,
			profile: &newProfile,
			wantErr: true,
			want:    nil,
		},
		{
			name:    "invalid profile - first name is empty",
			actor:   staffActor,
			profile: &invProfileNoFirstName,
			wantErr: true,
			want:    nil,
		},
		{
			name:    "invalid profile - mother last name is empty",
			actor:   staffActor,
			profile: &invProfileNoMotherLastName,
			wantErr: true,
			want:    nil,
		},
		{
			name:    "invalid profile - father last name is empty",
			actor:   staffActor,
			profile: &invProfileNoFatherLastName,
			wantErr: true,
			want:    nil,
		},
		{
			name:    "invalid profile - document number is empty",
			actor:   staffActor,
			profile: &invProfileNoDocumentNumber,
			wantErr: true,
			want:    nil,
		},
		{
			name:    "invalid profile - phone is empty",
			actor:   staffActor,
			profile: &invProfileNoPhone,
			wantErr: true,
			want:    nil,
		},
		{
			name:    "invalid profile - date of birth is empty",
			actor:   staffActor,
			profile: &invProfileNoDateOfBirth,
			wantErr: true,
			want:    nil,
		},
		{
			name:    "invalid profile - cmp is empty",
			actor:   staffActor,
			profile: &invProfileNoCmp,
			wantErr: true,
			want:    nil,
		},
		{
			name:    "invalid profile - specialty is empty",
			actor:   staffActor,
			profile: &invProfileNoSpecialty,
			wantErr: true,
			want:    nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			//WHEN
			err := uc.AddDoctor(&tc.actor, tc.profile)

			// THEN
			if tc.wantErr {
				if err == nil {
					t.Errorf("Expecting error not to be nil")
				}
			} else {
				if err != nil {
					t.Errorf("Error adding profile: %v", err)
				}

				profiles, _ := profileRepo.GetAll()
				found := false

				for _, profile := range profiles {
					if profile.ID == tc.profile.ID {
						found = true
						break
					}
				}
				if !found {
					t.Error("profile not added to repository")

				}

			}
		})

	}

}

func TestUpdateProfile(t *testing.T) {
	//GIVEN
	uc, _ := getProfileUC()

	//doctor
	doctorProfile := entities.NewFakeProfile()
	doctorProfile.Role = entities.DOCTOR
	doctorActor := entities.NewUser(1, "italo@blockehr.pe", "password", 1, doctorProfile.ID, doctorProfile)
	//staff
	staffProfile := entities.NewFakeProfile()
	staffProfile.Role = entities.STAFF
	staffActor := entities.NewUser(2, "italo@blockehr.pe", "password", 1, staffProfile.ID, staffProfile)
	//profile
	validProfile := entities.Profile{ID: 2, FirstName: "Jorge", FatherLastName: "De la Flor", MotherLastName: "Valle", DocumentNumber: "40058778", Gender: "Masculino", Phone: "94785987", DateOfBirth: "2024-03-01", Cmp: "1", Specialty: "Psicologia", Role: "DOCTOR"}
	//invalid

	// ID = 1 is Staff
	// ID = 2 is Doctor
	// ID = 3 is Patient
	invProfileNoFirstName := entities.Profile{ID: 2, FirstName: "", MotherLastName: "Kunimoto", FatherLastName: "Kunimoto", DocumentNumber: "40058778", Gender: "Masculino", Phone: "94785987", DateOfBirth: "2024-03-01", Cmp: "5", Specialty: "Psicologia", Role: "DOCTOR"}
	invProfileNoMotherLastName := entities.Profile{ID: 2, FirstName: "Rodrigo", MotherLastName: "", FatherLastName: "Kunimoto", DocumentNumber: "40058778", Gender: "Masculino", Phone: "94785987", DateOfBirth: "2024-03-01", Cmp: "5", Specialty: "Psicologia", Role: "DOCTOR"}
	invProfileNoFatherLastName := entities.Profile{ID: 2, FirstName: "Rodrigo", MotherLastName: "Kunimoto", FatherLastName: "", DocumentNumber: "40058778", Gender: "Masculino", Phone: "94785987", DateOfBirth: "2024-03-01", Cmp: "5", Specialty: "Psicologia", Role: "DOCTOR"}
	invProfileNoDocumentNumber := entities.Profile{ID: 2, FirstName: "Rodrigo", MotherLastName: "Kunimoto", FatherLastName: "Kunimoto", DocumentNumber: "", Gender: "Masculino", Phone: "94785987", DateOfBirth: "2024-03-01", Cmp: "5", Specialty: "Psicologia", Role: "DOCTOR"}
	invProfileNoPhone := entities.Profile{ID: 2, FirstName: "Rodrigo", MotherLastName: "Kunimoto", FatherLastName: "Kunimoto", DocumentNumber: "40058778", Gender: "Masculino", Phone: "", DateOfBirth: "2024-03-01", Cmp: "5", Specialty: "Psicologia", Role: "DOCTOR"}
	invProfileNoDateOfBirth := entities.Profile{ID: 2, FirstName: "Rodrigo", MotherLastName: "Kunimoto", FatherLastName: "Kunimoto", DocumentNumber: "40058778", Gender: "Masculino", Phone: "94785987", DateOfBirth: "", Cmp: "5", Specialty: "Psicologia", Role: "DOCTOR"}
	invProfileNoCmp := entities.Profile{ID: 2, FirstName: "Rodrigo", MotherLastName: "Kunimoto", FatherLastName: "Kunimoto", DocumentNumber: "40058778", Gender: "Masculino", Phone: "94785987", DateOfBirth: "2024-03-01", Cmp: "", Specialty: "Psicologia", Role: "DOCTOR"}
	invProfileNoSpecialty := entities.Profile{ID: 2, FirstName: "Rodrigo", MotherLastName: "Kunimoto", FatherLastName: "Kunimoto", DocumentNumber: "40058778", Gender: "Masculino", Phone: "94785987", DateOfBirth: "2024-03-01", Cmp: "5", Specialty: "", Role: "DOCTOR"}
	invProfileNoGender := entities.Profile{ID: 2, FirstName: "Rodrigo", MotherLastName: "Kunimoto", FatherLastName: "Kunimoto", DocumentNumber: "40058778", Gender: "", Phone: "94785987", DateOfBirth: "2024-03-01", Cmp: "5", Specialty: "Psicologia", Role: "DOCTOR"}

	testCases := []struct {
		name    string
		actor   entities.User
		profile *entities.Profile
		want    error
		wantErr bool
	}{
		{
			name:    "valid profile with staff",
			actor:   staffActor,
			profile: &validProfile,
			wantErr: false,
			want:    nil,
		},

		{
			name:    "invalid profile - profile made by  non-staff",
			actor:   doctorActor,
			profile: &validProfile,
			wantErr: true,
			want:    nil,
		},

		{
			name:    "invalid profile - first name is empty",
			actor:   staffActor,
			profile: &invProfileNoFirstName,
			wantErr: true,
			want:    nil,
		},

		{
			name:    "invalid profile - mother last name is empty",
			actor:   staffActor,
			profile: &invProfileNoMotherLastName,
			wantErr: true,
			want:    nil,
		},

		{
			name:    "invalid profile - father last name is empty",
			actor:   staffActor,
			profile: &invProfileNoFatherLastName,
			wantErr: true,
			want:    nil,
		},
		{
			name:    "invalid profile - document number is empty",
			actor:   staffActor,
			profile: &invProfileNoDocumentNumber,
			wantErr: true,
			want:    nil,
		},

		{
			name:    "invalid profile - phone is empty",
			actor:   staffActor,
			profile: &invProfileNoPhone,
			wantErr: true,
			want:    nil,
		},
		{
			name:    "invalid profile - date of birth is empty",
			actor:   staffActor,
			profile: &invProfileNoDateOfBirth,
			wantErr: true,
			want:    nil,
		},

		{
			name:    "invalid profile - cmp is empty",
			actor:   staffActor,
			profile: &invProfileNoCmp,
			wantErr: true,
			want:    nil,
		},
		{
			name:    "invalid profile - specialty is empty",
			actor:   staffActor,
			profile: &invProfileNoSpecialty,
			wantErr: true,
			want:    nil,
		},
		{
			name:    "invalid profile - gender is empty",
			actor:   staffActor,
			profile: &invProfileNoGender,
			wantErr: true,
			want:    nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			//WHEN
			err := uc.Update(&tc.actor, tc.profile.ID, tc.profile)

			// THEN
			if tc.wantErr {
				if err == nil {
					t.Errorf("Expecting error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("Error updating profile: %v", err)
				}

			}
		})

	}
}

func TestGetByDocumentNumber(t *testing.T) {

	//GIVEN
	uc, _ := getProfileUC()

	//doctor
	doctorProfile := entities.NewFakeProfile()
	doctorProfile.Role = entities.DOCTOR
	doctorActor := entities.NewUser(1, "italo@blockehr.pe", "password", 1, doctorProfile.ID, doctorProfile)
	//staff
	staffProfile := entities.NewFakeProfile()
	staffProfile.Role = entities.STAFF
	staffActor := entities.NewUser(2, "italo@blockehr.pe", "password", 1, staffProfile.ID, staffProfile)
	//profile
	validProfile := entities.Profile{ID: 4, FirstName: "Italo", MotherLastName: "Kunimoto", FatherLastName: "Luna", DocumentNumber: "40058779", Phone: "94785987", DateOfBirth: "2024-03-01", Cmp: "0", Specialty: "Psicologia", Role: "PATIENT"}
	//invalid
	invProfileNoDocumentNumber := entities.Profile{ID: 4, FirstName: "Italo", MotherLastName: "Kunimoto", FatherLastName: "Luna", DocumentNumber: "", Gender: "Masculino",
		Phone: "94785987", DateOfBirth: "2024-03-01", Cmp: "0", Specialty: "Psicologia", Role: "PATIENT"}
	testCases := []struct {
		name    string
		actor   entities.User
		profile *entities.Profile
		want    error
		wantErr bool
	}{
		{
			name:    "valid profile with staff",
			actor:   staffActor,
			profile: &validProfile,
			wantErr: false,
			want:    nil,
		},

		{
			name:    "invalid profile - profile made by  non-staff",
			actor:   doctorActor,
			profile: &validProfile,
			wantErr: true,
			want:    nil,
		},

		{
			name:    "invalid profile - document number is empty",
			actor:   staffActor,
			profile: &invProfileNoDocumentNumber,
			wantErr: true,
			want:    nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			//WHEN
			_, err := uc.GetByDocumentNumber(&tc.actor, tc.profile.DocumentNumber)

			// THEN
			if tc.wantErr {
				if err == nil {
					t.Errorf("Error getting profile: %v", err)
				}
			} else {
				if err != nil {
					t.Errorf("Error getting profile: %v", err)
				}
			}
		})
	}

}

func TestGetAllProfiles(t *testing.T) {

	//GIVEN
	uc, _ := getProfileUC()

	//doctor
	doctorProfile := entities.NewFakeProfile()
	doctorProfile.Role = entities.DOCTOR
	doctorActor := entities.NewUser(1, "italo@blockehr.pe", "password", 1, doctorProfile.ID, doctorProfile)
	//staff
	staffProfile := entities.NewFakeProfile()
	staffProfile.Role = entities.STAFF
	staffActor := entities.NewUser(2, "italo@blockehr.pe", "password", 1, staffProfile.ID, staffProfile)
	testCases := []struct {
		name    string
		actor   entities.User
		want    error
		wantErr bool
	}{
		{
			name:    "valid profile with staff",
			actor:   staffActor,
			wantErr: false,
			want:    nil,
		},

		{
			name:    "invalid profile - profile made by  non-staff",
			actor:   doctorActor,
			wantErr: true,
			want:    nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			//WHEN
			_, err := uc.GetAll(&tc.actor)

			// THEN
			if tc.wantErr {
				if err == nil {
					t.Errorf("Error getting profiles: %v", err)
				}
			} else {
				if err != nil {
					t.Errorf("Error getting profiles: %v", err)
				}

			}
		})
	}
}

func TestGetAllDoctors(t *testing.T) {

	//GIVEN
	uc, _ := getProfileUC()

	//doctor
	doctorProfile := entities.NewFakeProfile()
	doctorProfile.Role = entities.DOCTOR
	doctorActor := entities.NewUser(1, "italo@blockehr.pe", "password", 1, doctorProfile.ID, doctorProfile)
	//staff
	staffProfile := entities.NewFakeProfile()
	staffProfile.Role = entities.STAFF
	staffActor := entities.NewUser(2, "italo@blockehr.pe", "password", 1, staffProfile.ID, staffProfile)
	testCases := []struct {
		name    string
		actor   entities.User
		want    error
		wantErr bool
	}{
		{
			name:    "valid profile with staff",
			actor:   staffActor,
			wantErr: false,
			want:    nil,
		},

		{
			name:    "invalid profile - profile made by  non-staff",
			actor:   doctorActor,
			wantErr: true,
			want:    nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			//WHEN
			_, err := uc.GetAllDoctors(&tc.actor)

			// THEN
			if tc.wantErr {
				if err == nil {
					t.Errorf("Error getting profiles: %v", err)
				}
			} else {
				if err != nil {
					t.Errorf("Error getting profiles: %v", err)
				}

			}
		})
	}
}

func TestGetAllPatients(t *testing.T) {

	//GIVEN
	uc, _ := getProfileUC()

	//doctor
	doctorProfile := entities.NewFakeProfile()
	doctorProfile.Role = entities.DOCTOR
	doctorActor := entities.NewUser(1, "italo@blockehr.pe", "password", 1, doctorProfile.ID, doctorProfile)
	//staff
	staffProfile := entities.NewFakeProfile()
	staffProfile.Role = entities.STAFF
	staffActor := entities.NewUser(2, "italo@blockehr.pe", "password", 1, staffProfile.ID, staffProfile)
	testCases := []struct {
		name    string
		actor   entities.User
		want    error
		wantErr bool
	}{
		{
			name:    "valid profile with staff",
			actor:   staffActor,
			wantErr: false,
			want:    nil,
		},

		{
			name:    "invalid profile - profile made by  non-staff",
			actor:   doctorActor,
			wantErr: true,
			want:    nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			//WHEN
			_, err := uc.GetAllPatients(&tc.actor)

			// THEN
			if tc.wantErr {
				if err == nil {
					t.Errorf("Error getting profiles: %v", err)
				}
			} else {
				if err != nil {
					t.Errorf("Error getting profiles: %v", err)
				}

			}
		})
	}
}
func TestGetProfileById(t *testing.T) {

	//GIVEN
	uc, _ := getProfileUC()

	//doctor
	doctorProfile := entities.NewFakeProfile()
	doctorProfile.Role = entities.DOCTOR
	doctorActor := entities.NewUser(1, "italo@blockehr.pe", "password", 1, doctorProfile.ID, doctorProfile)
	//staff
	staffProfile := entities.NewFakeProfile()
	staffProfile.Role = entities.STAFF
	staffActor := entities.NewUser(2, "italo@blockehr.pe", "password", 1, staffProfile.ID, staffProfile)
	testCases := []struct {
		name    string
		actor   entities.User
		Id      int
		want    error
		wantErr bool
	}{
		{
			name:    "valid profile with staff",
			actor:   staffActor,
			Id:      1,
			wantErr: false,
			want:    nil,
		},

		{
			name:    "invalid profile - profile made by  non-staff",
			actor:   doctorActor,
			Id:      1,
			wantErr: true,
			want:    nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			//WHEN
			_, err := uc.GetById(&tc.actor, tc.Id)

			// THEN
			if tc.wantErr {
				if err == nil {
					t.Errorf("Error getting profile: %v", err)
				}
			} else {
				if err != nil {
					t.Errorf("Error getting profile: %v", err)
				}

			}
		})
	}
}
