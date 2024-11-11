package mock_repositories

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/open-wm/blockehr/pkg/domain/entities"
	"github.com/open-wm/blockehr/pkg/interfaces"
)

type InMemoryAppointmentRepository struct {
	appointments []*entities.Appointment
}

func NewInMemoryAppointmentRepository() interfaces.AppointmentRepository {
	tomorrow := time.Now().AddDate(0, 0, 1)
	yesterday := time.Now().AddDate(0, 0, -1)
	appointmentsDefault := []*entities.Appointment{
		{
			ID:        5,
			Specialty: "cardiologia",
			Status:    entities.PAID,
			StartsAt:  tomorrow,
			DoctorId:  2,
			Doctor:    entities.NewFakeUser(),
			PatientId: 3,
		},
		{
			ID:        6,
			Specialty: "cardiologia",
			Status:    entities.PAID,
			StartsAt:  yesterday,
			DoctorId:  2,
			Doctor:    entities.NewFakeUser(),
			PatientId: 3,
		},
		{
			ID:        7,
			Specialty: "Psicólogo",
			Status:    entities.PENDING,
			StartsAt:  tomorrow,
			DoctorId:  2,
			Doctor:    entities.NewFakeUser(),
			PatientId: 3,
		},
	}
	return &InMemoryAppointmentRepository{
		appointments: appointmentsDefault,
	}
}
func (m *InMemoryAppointmentRepository) Add(appointment *entities.Appointment) error {
	m.appointments = append(m.appointments, appointment)
	return nil
}
func (m *InMemoryAppointmentRepository) Update(appointment *entities.Appointment) error {
	for i, a := range m.appointments {
		if a.ID == appointment.ID {
			m.appointments[i] = appointment
			return nil
		}
	}
	return errors.New("appointment not found")
}
func (m *InMemoryAppointmentRepository) GetAll() ([]*entities.Appointment, error) {
	if len(m.appointments) == 0 {
		return nil, errors.New("no appointments found")
	}
	return m.appointments, nil
}
func (m *InMemoryAppointmentRepository) GetById(id int) (*entities.Appointment, error) {
	for _, appointment := range m.appointments {
		if appointment.ID == id {
			return appointment, nil
		}
	}
	return nil, errors.New("appointment not found for ID")
}
func (m *InMemoryAppointmentRepository) GetByDoctorID(doctorID int) ([]*entities.Appointment, error) {
	var result []*entities.Appointment
	for _, appointment := range m.appointments {
		if appointment.DoctorId == doctorID {
			result = append(result, appointment)
		}
	}
	return result, nil
}
func (m *InMemoryAppointmentRepository) GetByPatientID(patientID int) ([]*entities.Appointment, error) {
	var result []*entities.Appointment
	for _, appointment := range m.appointments {
		if appointment.PatientId == patientID {
			result = append(result, appointment)
		}
	}
	return result, nil
}

/* -------------------------------------------------------------------
	Files
------------------------------------------------------------------- */

type InMemoryFileRepository struct {
	files []*entities.File
}

func NewInMemoryFileRepository() *InMemoryFileRepository {

	files := []*entities.File{
		{
			ID:       2,
			Url:      "https://medlineplus.gov/images/Xray_share.jpg",
			Name:     "Radiografía",
			FileSize: "455 KB",
			MimeType: ".jpg",
			RecordId: 3,
		},
	}

	return &InMemoryFileRepository{files: files}
}
func (m *InMemoryFileRepository) Add(file *entities.File) error {
	m.files = append(m.files, file)
	return nil
}
func (m *InMemoryFileRepository) GetByRecordId(id int) ([]*entities.File, error) {
	var filesRecordId []*entities.File

	for _, file := range m.files {
		if file.RecordId == id {
			fmt.Print(file.Name)
			filesRecordId = append(filesRecordId, file)
		}
	}

	if len(filesRecordId) > 0 {
		return filesRecordId, nil
	}

	return nil, errors.New("files not found for record ID")
}

func (m *InMemoryFileRepository) GetById(id int) (*entities.File, error) {
	for _, file := range m.files {
		if file.ID == id {
			return file, nil
		}
	}
	return nil, errors.New("file not found for ID")
}

/* -------------------------------------------------------------------
	PROFILE
------------------------------------------------------------------- */

type InMemoryProfileRepository struct {
	profiles []*entities.Profile
}

func NewInMemoryProfileRepository() *InMemoryProfileRepository {
	profiles := []*entities.Profile{
		{
			ID:             1,
			FirstName:      "Rodrigo",
			MotherLastName: "Kunimoto",
			FatherLastName: "Luna",
			DocumentNumber: "40058778",
			Gender:         "Masculino",
			Phone:          "94785987",
			DateOfBirth:    "2024-03-01",
			Cmp:            "10",
			Specialty:      "Psicologia",
			Role:           "STAFF",
		},
		{
			ID:             2,
			FirstName:      "Italo",
			MotherLastName: "Kunimoto",
			FatherLastName: "Aguilar",
			DocumentNumber: "72549855",
			Gender:         "Masculino",
			Phone:          "974528438",
			DateOfBirth:    "12-29-2004",
			Cmp:            "1234567",
			Specialty:      "psicologia",
			Role:           "DOCTOR",
		},

		{
			ID:             3,
			FirstName:      "Rodrigo",
			MotherLastName: "Kunimoto",
			FatherLastName: "Luna",
			DocumentNumber: "40058779",
			Gender:         "Masculino",

			Phone:       "94785987",
			DateOfBirth: "2024-03-01",
			Cmp:         "10",
			Specialty:   "Psicologia",
			Role:        "PATIENT",
		},
		{
			ID:             4,
			FirstName:      "Dr. 2",
			MotherLastName: "Kunimoto",
			FatherLastName: "Luna",
			DocumentNumber: "40058779",
			Gender:         "Masculino",
			Phone:          "94785987",
			DateOfBirth:    "2024-03-01",
			Cmp:            "10",
			Specialty:      "Psicologia",
			Role:           "DOCTOR",
		},
	}

	return &InMemoryProfileRepository{
		profiles: profiles,
	}
}

func (m *InMemoryProfileRepository) Add(profile *entities.Profile) error {
	m.profiles = append(m.profiles, profile)

	return nil
}

func (m *InMemoryProfileRepository) Update(profile *entities.Profile) error {
	for i, a := range m.profiles {
		if a.ID == profile.ID {
			m.profiles[i] = profile
			log.Printf("profile updated: %v", profile)
			return nil
		}
	}
	return errors.New("profile not found")
}

func (m *InMemoryProfileRepository) GetByDocumentNumber(documentNumber string) (*entities.Profile, error) {

	for _, profile := range m.profiles {
		if profile.DocumentNumber == documentNumber {
			if profile.Role == entities.PATIENT {
				return profile, nil

			}
		}
	}
	return nil, errors.New("profile not found")
}

func (m *InMemoryProfileRepository) GetAll() ([]*entities.Profile, error) {
	if len(m.profiles) == 0 {
		return nil, errors.New("no profiles found")
	}
	return m.profiles, nil
}
func (m *InMemoryProfileRepository) GetById(id int) (*entities.Profile, error) {
	for _, profile := range m.profiles {
		if profile.ID == id {
			return profile, nil
		}
	}
	return nil, errors.New("profile not found for ID")
}

/* -------------------------------------------------------------------
	(Health)Record
------------------------------------------------------------------- */

type InMemoryRecordRepository struct {
	records []*entities.Record
}

func NewInMemoryRecordRepository() *InMemoryRecordRepository {
	records := []*entities.Record{
		{
			ID:        2,
			Body:      "Rodrigo",
			CreatedAt: "01-06-2024",
			UpdatedAt: "01-06-2024",
			PatientId: 3,
			DoctorId:  2,
			Files: []*entities.File{
				{
					ID:       2,
					Url:      "https://medlineplus.gov/images/Xray_share.jpg",
					Name:     "Radiografía",
					FileSize: "455 KB",
					MimeType: ".jpg",
					RecordId: 2,
				},
			},
		},
		{
			ID:        3,
			Body:      "Rodrigo",
			CreatedAt: "01-06-2024",
			UpdatedAt: "01-06-2024",
			PatientId: 3,
			DoctorId:  2,
			Files: []*entities.File{
				{
					ID:       3,
					Url:      "https://medlineplus.gov/images/Xray_share.jpg",
					Name:     "Radiografía",
					FileSize: "455 KB",
					MimeType: ".jpg",
					RecordId: 3,
				},
			},
		},
	}

	return &InMemoryRecordRepository{
		records: records,
	}
}

func (m *InMemoryRecordRepository) Add(record *entities.Record) (*entities.Record, error) {
	m.records = append(m.records, record)
	record.ID = len(m.records)
	return record, nil
}

func (m *InMemoryRecordRepository) GetAll() ([]*entities.Record, error) {
	if len(m.records) == 0 {
		return nil, errors.New("no records found")
	}
	return m.records, nil

}

func (m *InMemoryRecordRepository) GetById(id int) (*entities.Record, error) {
	for _, record := range m.records {
		if record.ID == id {
			return record, nil
		}
	}
	return nil, errors.New("record not found")
}

func (m *InMemoryRecordRepository) UpdateByPatientID(patientID int, newBody string) error {
	for _, record := range m.records {
		if record.PatientId == patientID {
			record.Body = newBody
			return nil
		}
	}
	return fmt.Errorf("no record found for patient with ID %d", patientID)
}

/* -------------------------------------------------------------------
	User
------------------------------------------------------------------- */

type InMemoryUserRepository struct {
	users []*entities.User
}

func defaultUsers() []*entities.User {
	return []*entities.User{
		{

			ID:             1,
			Email:          "staff@blockehr.pe",
			Password:       "perroLoco",
			HealthCenterId: 1,
			ProfileId:      1,
			Profile: entities.Profile{
				ID:             1,
				FirstName:      "Italo",
				MotherLastName: "Castillo",
				FatherLastName: "Kunimoto",
				DocumentNumber: "123456788",
				Phone:          "123456789",
				DateOfBirth:    "01-06-2024",
				Cmp:            "123456",
				Specialty:      "Psicologia",
				Role:           entities.STAFF,
			},
		},
		{
			ID:             2,
			Email:          "doctor@blockehr.pe",
			Password:       "perroLoco",
			HealthCenterId: 1,
			ProfileId:      2,
			Profile: entities.Profile{
				ID:             2,
				FirstName:      "Rodrigo",
				MotherLastName: "Castillo",
				FatherLastName: "Kunimoto",
				DocumentNumber: "123456789",
				Phone:          "123456789",
				DateOfBirth:    "01-06-2024",
				Role:           entities.DOCTOR,
			},
		},
		{

			ID:             3,
			Email:          "patient@blockehr.pe",
			Password:       "perroLoco",
			HealthCenterId: 1,
			ProfileId:      3,
			Profile: entities.Profile{
				ID:             2,
				FirstName:      "Juan",
				MotherLastName: "Castillo",
				FatherLastName: "Kunimoto",
				DocumentNumber: "123456799",
				Phone:          "123456789",
				DateOfBirth:    "01-06-2024",
				Role:           entities.PATIENT,
			},
		},
		{
			ID:             4,
			Email:          "doctor2@blockehr.pe",
			Password:       "perroLoco",
			HealthCenterId: 2,
			ProfileId:      4,
			Profile: entities.Profile{
				ID:             4,
				FirstName:      "Juan",
				MotherLastName: "Castillo",
				FatherLastName: "Kunimoto",
				DocumentNumber: "123456799",
				Phone:          "123456789",
				DateOfBirth:    "01-06-2024",
				Cmp:            "123456",
				Specialty:      "Psicologia",
				Role:           entities.DOCTOR,
			},
		},
	}
}

func NewInMemoryUserRepository() interfaces.UserRepository {
	return &InMemoryUserRepository{
		users: defaultUsers(),
	}
}
func (m *InMemoryUserRepository) Reset() {
	m.users = defaultUsers()
}

func (m *InMemoryUserRepository) SaveUser(user *entities.User) error {
	m.users = append(m.users, user)

	return nil
}

func (m *InMemoryUserRepository) UpdatePassword(id int, newPassword string) error {
	for _, user := range m.users {
		if user.ID == id {
			user.Password = newPassword
			fmt.Printf("ok")
			return nil
		}
	}
	return errors.New("user not found")
}

func (m *InMemoryUserRepository) GetUser(email string, password string) (*entities.User, error) {
	for _, user := range m.users {
		if user.Email == email && user.Password == password {
			return user, nil
		}
	}
	return nil, fmt.Errorf("user not found")
}
func (m *InMemoryUserRepository) GetUserByID(ID int) (*entities.User, error) {
	for _, user := range m.users {
		if user.ID == ID {
			return user, nil
		}
	}
	return nil, errors.New("user not found")
}
func (m *InMemoryUserRepository) GetUserByProfileID(ID int) (*entities.User, error) {
	for _, user := range m.users {
		if user.ProfileId == ID {
			return user, nil
		}
	}
	return nil, errors.New("user not found")
}

func (m *InMemoryUserRepository) GetAllDoctors() ([]*entities.User, error) {
	doctors := []*entities.User{}

	for _, user := range m.users {
		if user.Profile.Role == "DOCTOR" {
			doctors = append(doctors, user)
		}
	}

	if len(doctors) == 0 {
		return nil, errors.New("no doctors found")
	}

	return doctors, nil
}

func (m *InMemoryUserRepository) GetDoctorsForSpecialty(specialty string) ([]*entities.User, error) {
	doctors := []*entities.User{}

	for _, user := range m.users {
		if user.Profile.Specialty == specialty && user.Profile.Role == "DOCTOR" {
			doctors = append(doctors, user)
		}
	}

	return doctors, nil
}

//SESSION

type InMemorySessionRepository struct {
	sessions []*entities.Session
}

func NewInMemorySessionRepository() *InMemorySessionRepository {
	// Crear una sesión por defecto
	staffSession := &entities.Session{
		ID:        1,
		UserID:    1,
		Token:     "STAFF",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	doctorSession := &entities.Session{
		ID:        2,
		UserID:    2,
		Token:     "DOCTOR",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	expiredSession := &entities.Session{
		ID:        2,
		UserID:    2,
		Token:     "EXPIRED",
		CreatedAt: time.Now().Add(-time.Hour * 25),
		UpdatedAt: time.Now().Add(-time.Hour * 25),
	}

	// Crear el repositorio con la sesión por defecto
	repo := &InMemorySessionRepository{
		sessions: []*entities.Session{doctorSession, staffSession, expiredSession},
	}

	return repo
}

func (m *InMemorySessionRepository) GetSession(sessionToken string) *entities.Session {
	for _, session := range m.sessions {
		if session.Token == sessionToken {
			return session
		}
	}
	return nil
}

func (m *InMemorySessionRepository) DeleteSession(sessionToken string) error {
	for i, session := range m.sessions {
		if session.Token == sessionToken {
			m.sessions = append(m.sessions[:i], m.sessions[i+1:]...)
			return nil
		}
	}
	return nil
}

func (m *InMemorySessionRepository) SaveSession(session *entities.Session) (*entities.Session, error) {
	session.ID = len(m.sessions) + 1

	currentTime := time.Now()
	session.CreatedAt = currentTime
	session.UpdatedAt = currentTime

	m.sessions = append(m.sessions, session)

	return session, nil
}
