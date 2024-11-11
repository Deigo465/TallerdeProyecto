package usecase_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/open-wm/blockehr/pkg/domain/entities"
	usecase "github.com/open-wm/blockehr/pkg/domain/usecases"
	mock "github.com/open-wm/blockehr/pkg/mocks"
)

func getValidAppointment() entities.Appointment {
	specialty := "Cardiología"
	status := entities.PENDING
	tomorrow := time.Now().AddDate(0, 0, 1)

	doctor := entities.NewFakeUser()
	doctor.Profile.Role = entities.DOCTOR
	patient := entities.NewFakeProfile()
	patient.Role = entities.PATIENT
	description := "toi malito"

	newAppointment := entities.NewAppointment(1, specialty, status, tomorrow, doctor.ID, doctor, patient.ID, patient, description)

	return newAppointment
}

func TestAddAppointment(t *testing.T) {
	//GIVEN
	appointmentRepo := mock.NewInMemoryAppointmentRepository()
	userRepo := mock.NewInMemoryUserRepository()
	profileRepo := mock.NewInMemoryProfileRepository()
	blockchain := mock.NewInMemoryBlockchain()
	uc := usecase.NewAppointmentUsecase(appointmentRepo, userRepo, profileRepo, blockchain)

	//staff
	staffProfile := entities.NewFakeProfile()
	staffProfile.Role = entities.STAFF

	staffActor := entities.NewUser(1, "italo@blockehr.pe", "password", 1, staffProfile.ID, staffProfile)

	//doctor
	doctorProfile := entities.NewFakeProfile()
	doctorProfile.Role = entities.DOCTOR
	doctorUser := entities.NewUser(4, "italo@blockehr.pe", "password", 2, doctorProfile.ID, doctorProfile)

	doctorProfile2 := entities.NewFakeProfile()
	doctorProfile2.Role = entities.DOCTOR
	doctorUser2 := entities.NewUser(4, "italo@blockehr.pe", "password", 2, doctorProfile2.ID, doctorProfile2)

	//patient
	patient := entities.NewFakeProfile()
	patient.Role = entities.PATIENT

	tomorrow := time.Now().AddDate(0, 0, 1)

	// Valid
	validAppointment := getValidAppointment()
	appointment := entities.NewAppointment(1, "Cardiología", entities.PENDING, tomorrow, 1, doctorUser, patient.ID, patient, "toi malito")

	// Invalid
	invAppointmentNoDoctor := entities.NewAppointment(1, "Cardiología", entities.PENDING, tomorrow, 0, staffActor, patient.ID, patient, "toi malito")
	// yesterday := time.Now().AddDate(0, 0, -1)
	// invAppointPastDate := entities.NewAppointment(1, "Cardiología", entities.PENDING, yesterday, 1, staffActor, patient.ID, patient, "toi malito")

	invAppointNoSpecialty := entities.NewAppointment(1, "", entities.PENDING, tomorrow, 1, staffActor, patient.ID, patient, "toi malito")
	invAppointStatusInvalid := entities.NewAppointment(1, "Cardiología", -1, tomorrow, 1, staffActor, patient.ID, patient, "toi malito")

	invAppointNoSamehealthCenter := entities.NewAppointment(1, "Cardiología", entities.PENDING, tomorrow, doctorUser.ID, doctorUser, patient.ID, patient, "toi malito")
	doctorUser.HealthCenterId = 1
	invAppointNoDescription := entities.NewAppointment(1, "Cardiología", entities.PENDING, tomorrow, doctorUser.ID, doctorUser, patient.ID, patient, "")

	invAppointmentNoPatient := entities.NewAppointment(1, "Cardiología", entities.PENDING, tomorrow, staffActor.ID, staffActor, 0, patient, "toi malito")

	invAppointmentAppointmentAlreadyHappening := entities.NewAppointment(2, "Cardiología", entities.PENDING, tomorrow, 2, doctorUser2, patient.ID, patient, "toi malito")

	testCases := []struct {
		name        string
		appointment *entities.Appointment
		actor       entities.User
		want        error
		wantErr     bool
	}{
		{
			name:        "valid appointment with staff",
			actor:       staffActor,
			appointment: &validAppointment,
			wantErr:     false,
			want:        nil,
		},
		{
			name:        "invalid appointment - nil appointment",
			actor:       staffActor,
			appointment: nil,
			wantErr:     true,
			want:        nil,
		},
		{
			name:        "invalid appointment - appointment made by non-staff",
			actor:       doctorUser,
			appointment: &appointment,
			wantErr:     true,
			want:        nil,
		},
		{
			name:        "invalid appointment - appointment doctor not a doctor",
			actor:       staffActor,
			appointment: &invAppointmentNoDoctor,
			wantErr:     true,
			want:        nil,
		},
		{
			name:        "invalid appointment - specialty is empty",
			actor:       staffActor,
			appointment: &invAppointNoSpecialty,
			wantErr:     true,
			want:        nil,
		},
		{
			name:        "invalid appointment - status is not valid",
			actor:       staffActor,
			appointment: &invAppointStatusInvalid,
			wantErr:     true,
			want:        nil,
		},
		{
			name:        "invalid appointment - doctor is no same health center",
			actor:       staffActor,
			appointment: &invAppointNoSamehealthCenter,
			wantErr:     true,
			want:        nil,
		},
		{
			name:        "invalid appointment - description is empty",
			actor:       staffActor,
			appointment: &invAppointNoDescription,
			wantErr:     true,
			want:        nil,
		},
		{
			name:        "invalid appointment - Patient must exist",
			actor:       staffActor,
			appointment: &invAppointmentNoPatient,
			wantErr:     true,
			want:        nil,
		},
		{
			name:        "invalid appointment - Patient cant have 2 appointments with 2 doctors",
			actor:       staffActor,
			appointment: &invAppointmentAppointmentAlreadyHappening,
			wantErr:     true,
			want:        nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			//WHEN
			err := uc.Add(&tc.actor, tc.appointment)

			// THEN
			if tc.wantErr {
				if err == nil {
					t.Error("Expecting error but got nil")
				}
			} else {
				if err != nil {
					t.Fatalf("Error adding appointment: %v", err)
				}
				appointments, _ := appointmentRepo.GetAll()
				found := false

				for _, appointment := range appointments {
					if appointment.ID == tc.appointment.ID {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("Appointment not added to repository")
				}
			}
		})
	}

}

func TestUpdateStatus(t *testing.T) {

	//GIVEN
	appointmentRepo := mock.NewInMemoryAppointmentRepository()
	userRepo := mock.NewInMemoryUserRepository()
	profileRepo := mock.NewInMemoryProfileRepository()
	blockchain := mock.NewInMemoryBlockchain()
	uc := usecase.NewAppointmentUsecase(appointmentRepo, userRepo, profileRepo, blockchain)

	//staff
	staffProfile := entities.NewFakeProfile()
	staffProfile.Role = entities.STAFF

	staffActor := entities.NewUser(1, "italo@blockehr.pe", "password", 1, staffProfile.ID, staffProfile)

	//doctor
	doctorProfile := entities.NewFakeProfile()
	doctorProfile.Role = entities.DOCTOR
	doctorActor := entities.NewUser(1, "italo@blockehr.pe", "password", 1, doctorProfile.ID, doctorProfile)
	//patient
	patientProfile := entities.NewFakeProfile()
	patientProfile.Role = entities.PATIENT

	PatientActor := entities.NewUser(1, "italo@blockehr.pe", "password", 1, patientProfile.ID, patientProfile)

	//apointment
	tomorrow := time.Now().AddDate(0, 0, 1)
	validAppointment := entities.Appointment{ID: 5, Specialty: "Cardiología", Status: entities.PAID, StartsAt: tomorrow, DoctorId: doctorActor.ID, Doctor: doctorActor}
	//invalid

	testCases := []struct {
		name        string
		appointment *entities.Appointment
		actor       entities.User
		want        error
		wantErr     bool
	}{
		{
			name:        "valid appointment with staff",
			actor:       staffActor,
			appointment: &validAppointment,
			wantErr:     false,
			want:        nil,
		},

		{
			name:        "valid appointment with doctor",
			actor:       doctorActor,
			appointment: &validAppointment,
			wantErr:     false,
			want:        nil,
		},
		{
			name:        "invalid appointment - appointment updated by non-staff or doctor",
			actor:       PatientActor,
			appointment: &validAppointment,
			wantErr:     true,
			want:        nil,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			//WHEN

			err := uc.UpdateStatus(&tc.actor, tc.appointment.ID, tc.appointment)

			// THEN
			if tc.wantErr {
				if err == nil {
					t.Errorf("Error updating appointment: %v", err)
				}
			} else {
				if err != nil {
					t.Errorf("Error updating appointment: %v", err)
				}
			}
		})
	}

}

func TestGetAll(t *testing.T) {
	//GIVEN
	appointmentRepo := mock.NewInMemoryAppointmentRepository()
	userRepo := mock.NewInMemoryUserRepository()
	profileRepo := mock.NewInMemoryProfileRepository()
	blockchain := mock.NewInMemoryBlockchain()
	uc := usecase.NewAppointmentUsecase(appointmentRepo, userRepo, profileRepo, blockchain)

	//staff
	staffProfile := entities.NewFakeProfile()
	staffProfile.Role = entities.STAFF

	staffActor := entities.NewUser(2, "italo@blockehr.pe", "password", 1, staffProfile.ID, staffProfile)
	//Doctor
	doctorProfile := entities.NewFakeProfile()
	doctorProfile.Role = entities.DOCTOR
	doctorActor := entities.NewUser(5, "doctor@blockehr.pe", "perroLoco", 1, doctorProfile.ID, doctorProfile)

	testCases := []struct {
		name    string
		actor   entities.User
		want    error
		wantErr bool
	}{
		{
			name:    "valid appointment with staff",
			actor:   staffActor,
			wantErr: false,
			want:    nil,
		},
		{
			name:    "valid appointment with doctor",
			actor:   doctorActor,
			wantErr: false,
			want:    nil,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			//WHEN

			listAppointment, err := uc.GetAll(&tc.actor)

			// THEN
			if tc.wantErr {
				if err == nil {
					t.Errorf("Error getting appointment: %v", err)
				}
			} else {
				if err != nil {
					t.Errorf("Error getting appointment: %v", err)
				}
			}
			fmt.Print(len(listAppointment))
		})
	}
}

func TestGetById(t *testing.T) {
	//GIVEN
	appointmentRepo := mock.NewInMemoryAppointmentRepository()
	userRepo := mock.NewInMemoryUserRepository()
	profileRepo := mock.NewInMemoryProfileRepository()
	blockchain := mock.NewInMemoryBlockchain()
	uc := usecase.NewAppointmentUsecase(appointmentRepo, userRepo, profileRepo, blockchain)

	//staff
	staffProfile := entities.NewFakeProfile()
	staffProfile.Role = entities.STAFF

	staffActor := entities.NewUser(2, "italo@blockehr.pe", "password", 1, staffProfile.ID, staffProfile)
	//Doctor
	doctorProfile := entities.NewFakeProfile()
	doctorProfile.Role = entities.DOCTOR
	doctorActor := entities.NewUser(1, "italo@blockehr.pe", "password", 1, doctorProfile.ID, doctorProfile)
	//appointment
	tomorrow := time.Now().AddDate(0, 0, 1)
	validAppointment := entities.Appointment{ID: 5, Specialty: "Cardiología", Status: entities.PAID, StartsAt: tomorrow, DoctorId: doctorActor.ID,
		Doctor: doctorActor}
	//patient
	patientProfile := entities.NewFakeProfile()
	patientProfile.Role = entities.PATIENT
	patientActor := entities.NewUser(2, "rod@blockehr.pe", "password", 1, patientProfile.ID, patientProfile)

	testCases := []struct {
		name        string
		actor       entities.User
		appointment entities.Appointment
		want        error
		wantErr     bool
	}{
		{
			name:        "valid appointment with staff",
			actor:       staffActor,
			appointment: validAppointment,
			wantErr:     false,
			want:        nil,
		},
		{
			name:        "valid appointment with doctor",
			actor:       doctorActor,
			appointment: validAppointment,
			wantErr:     false,
			want:        nil,
		},
		{
			name:        "valid appointment with doctor",
			actor:       patientActor,
			appointment: validAppointment,
			wantErr:     true,
			want:        nil,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			//WHEN
			_, err := uc.GetById(&tc.actor, tc.appointment.ID)

			// THEN
			if tc.wantErr {
				if err == nil {
					t.Errorf("Error getting appointment: %v", err)
				}
			} else {
				if err != nil {
					t.Errorf("Error getting appointment: %v", err)
				}
			}
		})
	}
}
