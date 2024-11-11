package usecase

import (
	"bytes"
	"encoding/csv"
	"errors"
	"os"
	"strings"

	"github.com/open-wm/blockehr/pkg/domain/entities"
	"github.com/open-wm/blockehr/pkg/interfaces"
)

type ProfileUsecase interface {
	add(actor *entities.User, profile *entities.Profile) error
	AddDoctor(actor *entities.User, doctor *entities.Profile) error
	AddPatient(actor *entities.User, patient *entities.Profile) error
	Update(actor *entities.User, id int, profile *entities.Profile) error
	GetByDocumentNumber(actor *entities.User, documentNumber string) (*entities.Profile, error)
	GetAllDoctors(actor *entities.User) ([]*entities.Profile, error)
	GetAllPatients(actor *entities.User) ([]*entities.Profile, error)
	GetAll(actor *entities.User) ([]*entities.Profile, error)
	GetById(actor *entities.User, id int) (*entities.Profile, error)
	GetAllAppointmentsForPatient(actor *entities.User, patientId int) ([]*AppointmentWithDoctor, error)
}

type profileUsecase struct {
	profileRepo      interfaces.ProfileRepository
	appointmentsRepo interfaces.AppointmentRepository
	userRepo         interfaces.UserRepository
}

func NewProfileUsecase(repo interfaces.ProfileRepository, appointmentsRepo interfaces.AppointmentRepository, userRepo interfaces.UserRepository) ProfileUsecase {
	if strings.Contains(os.Args[0], ".test") {
		rootPath = "../../../"
	}
	return &profileUsecase{
		profileRepo:      repo,
		appointmentsRepo: appointmentsRepo,
		userRepo:         userRepo,
	}
}

var rootPath = "."
var ErrBadDoctorData error = errors.New("doctor data doesn't match")
var ErrDocumentAlreadyExists error = errors.New("document number already exists")

func (uc *profileUsecase) AddDoctor(actor *entities.User, doctor *entities.Profile) error {
	doctor.Role = entities.DOCTOR

	if doctor.Cmp == "" {
		return errors.New("Expecting cmp to not be empty")
	}
	if doctor.Specialty == "" {
		return errors.New("Expecting specialty to not be empty")
	}

	// check if there is already a doctor with the same
	if profile, _ := uc.profileRepo.GetByDocumentNumber(doctor.DocumentNumber); profile != nil {
		return ErrDocumentAlreadyExists
	}

	if !matchDoctor(doctor.Cmp, doctor.FirstName, doctor.FatherLastName, doctor.MotherLastName) {
		return ErrBadDoctorData
	}
	err := uc.add(actor, doctor)
	if err != nil {
		return err
	}

	newDoctor := entities.NewUser(0, doctor.ContactEmail, "123456", actor.HealthCenterId, doctor.ID, *doctor)
	return uc.userRepo.SaveUser(&newDoctor)
}

func matchDoctor(cmp, firstName, fatherLastName, motherLastName string) bool {
	b, err := os.ReadFile(rootPath + "/public/doctors.csv")
	if err != nil {
		panic(err)
	}
	r := bytes.NewReader(b)

	reader := csv.NewReader(r)

	reader.Comma = ';'

	lines, err := reader.ReadAll()
	if err != nil {
		panic(err)
	}
	for _, line := range lines {
		if strings.EqualFold(line[0], cmp) &&
			strings.EqualFold(line[1], fatherLastName) &&
			strings.EqualFold(line[2], motherLastName) &&
			strings.EqualFold(line[3], firstName) {
			return true
		}
	}
	return false
}

func (uc *profileUsecase) AddPatient(actor *entities.User, patient *entities.Profile) error {
	patient.Role = entities.PATIENT

	// check if there is already a doctor with the same
	if profile, _ := uc.profileRepo.GetByDocumentNumber(patient.DocumentNumber); profile != nil {
		return ErrDocumentAlreadyExists
	}

	return uc.add(actor, patient)
}

func (uc *profileUsecase) add(actor *entities.User, profile *entities.Profile) error {
	// check if role of actor is staff

	if actor.Profile.Role != "STAFF" {
		return errors.New("only staff can add profile")
	}
	// check if the fields are not empty

	if profile.FirstName == "" {
		return errors.New("Expecting first name to not be empty")
	}
	if profile.MotherLastName == "" {
		return errors.New("Expecting mother last name to not be empty")
	}
	if profile.FatherLastName == "" {
		return errors.New("Expecting father last name to not be empty")
	}
	if profile.DocumentNumber == "" {
		return errors.New("Expecting document number to not be empty")
	}
	if profile.Gender == "" {
		return errors.New("Expecting gender to not be empty")
	}
	if profile.Phone == "" {
		return errors.New("Expecting phone to not be empty")
	}
	if profile.DateOfBirth == "" {
		return errors.New("Expecting date of birth name to not be empty")
	}
	if profile.Role == "" {
		return errors.New("Expecting role to not be empty")
	}

	if err := uc.profileRepo.Add(profile); err != nil {
		return err
	}

	return nil
}

func (uc *profileUsecase) Update(actor *entities.User, id int, profile *entities.Profile) error {
	// check if role of actor is staff

	if actor.Profile.Role != "STAFF" {
		return errors.New("only staff can update profile")
	}
	// check if the fields are not empty

	oldProfile, err := uc.profileRepo.GetById(id)
	if err != nil {
		return err
	}

	profile.Role = oldProfile.Role

	if profile.FirstName == "" {
		return errors.New("Expecting first name to not be empty")
	}
	if profile.MotherLastName == "" {
		return errors.New("Expecting mother last name to not be empty")
	}
	if profile.FatherLastName == "" {
		return errors.New("Expecting father last name to not be empty")
	}
	if profile.DocumentNumber == "" {
		return errors.New("Expecting document number to not be empty")
	}
	if profile.Gender == "" {
		return errors.New("Expecting gender to not be empty")
	}
	if profile.Phone == "" {
		return errors.New("Expecting phone to not be empty")
	}
	if profile.DateOfBirth == "" {
		return errors.New("Expecting date of birth name to not be empty")
	}
	if profile.Role == entities.DOCTOR {
		if profile.Cmp == "" {
			return errors.New("Expecting cmp to not be empty")
		}
		if profile.Specialty == "" {
			return errors.New("Expecting specialty to not be empty")
		}
	}

	profile.ID = id
	if err := uc.profileRepo.Update(profile); err != nil {
		return err
	}

	return nil
}

func (uc *profileUsecase) GetByDocumentNumber(actor *entities.User, documentNumber string) (*entities.Profile, error) {
	// check if role of actor is staff

	if actor.Profile.Role != "STAFF" {
		return nil, errors.New("only staff can get profile")
	}
	if documentNumber == "" {
		return nil, errors.New("Expecting document number to not be empty")
	}

	profile, err := uc.profileRepo.GetByDocumentNumber(documentNumber)
	if err != nil {

		return nil, err
	}
	return profile, nil
}

func (uc *profileUsecase) GetAllDoctors(actor *entities.User) ([]*entities.Profile, error) {
	// check if role of actor is staff

	if actor.Profile.Role != "STAFF" {
		return nil, errors.New("only staff can get profile")
	}

	profiles, err := uc.profileRepo.GetAll()
	if err != nil {
		return nil, err
	}

	doctors := []*entities.Profile{}
	for _, profile := range profiles {
		if profile.Role == "DOCTOR" {
			doctors = append(doctors, profile)
		}
	}

	return doctors, nil

}
func (uc *profileUsecase) GetAllPatients(actor *entities.User) ([]*entities.Profile, error) {
	// check if role of actor is staff

	if actor.Profile.Role != "STAFF" {
		return nil, errors.New("only staff can get profile")
	}
	profiles, err := uc.profileRepo.GetAll()
	if err != nil {
		return nil, err
	}

	patients := []*entities.Profile{}
	for _, profile := range profiles {
		if profile.Role == "PATIENT" {
			patients = append(patients, profile)
		}
	}

	return patients, nil
}
func (uc *profileUsecase) GetAll(actor *entities.User) ([]*entities.Profile, error) {
	// check if role of actor is staff
	if actor.Profile.Role != "STAFF" {
		return nil, errors.New("only staff can get profiles")
	}

	profiles, err := uc.profileRepo.GetAll()
	if err != nil {
		return nil, err
	}

	return profiles, nil
}

func (uc *profileUsecase) GetById(actor *entities.User, id int) (*entities.Profile, error) {

	if actor.Profile.Role != "STAFF" {
		return nil, errors.New("only staff can get profiles")
	}

	profile, err := uc.profileRepo.GetById(id)
	if err != nil {
		return nil, err
	}

	return profile, nil
}

func (uc *profileUsecase) GetAllAppointmentsForPatient(actor *entities.User, patientId int) ([]*AppointmentWithDoctor, error) {
	if actor.Profile.Role != "STAFF" {
		return nil, errors.New("only staff can get profiles")
	}

	appointments, err := uc.appointmentsRepo.GetByPatientID(patientId)
	if err != nil {
		return nil, err
	}

	appointmentsWithDoctor := []*AppointmentWithDoctor{}
	for _, appointment := range appointments {
		doctor, err := uc.profileRepo.GetById(appointment.DoctorId)
		if err != nil {
			return nil, err
		}
		appointmentsWithDoctor = append(appointmentsWithDoctor, &AppointmentWithDoctor{
			Appointment: *appointment,
			Doctor:      *doctor,
		})
	}

	return appointmentsWithDoctor, nil
}
