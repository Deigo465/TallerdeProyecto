package interfaces

import (
	"github.com/open-wm/blockehr/pkg/domain/entities"
)

type AppointmentRepository interface {
	Add(appointment *entities.Appointment) error
	Update(appointment *entities.Appointment) error
	GetAll() ([]*entities.Appointment, error)
	GetById(id int) (*entities.Appointment, error)
	GetByDoctorID(doctorID int) ([]*entities.Appointment, error)
	GetByPatientID(patientID int) ([]*entities.Appointment, error)
}
type FileRepository interface {
	Add(file *entities.File) error
	GetByRecordId(id int) ([]*entities.File, error)
	GetById(id int) (*entities.File, error)
}

type HealthCenterRepository interface {
	Add(healthCenter *entities.HealthCenter) error
	Update(healthCenter *entities.HealthCenter) error
	GetAll() ([]*entities.HealthCenter, error)
	GetByID(id int) (*entities.HealthCenter, error)
}
type ProfileRepository interface {
	Add(profile *entities.Profile) error
	Update(profile *entities.Profile) error
	GetByDocumentNumber(documentNumber string) (*entities.Profile, error)
	GetAll() ([]*entities.Profile, error)
	GetById(id int) (*entities.Profile, error)
}
type RecordRepository interface {
	Add(record *entities.Record) (*entities.Record, error)
	GetAll() ([]*entities.Record, error)
	GetById(id int) (*entities.Record, error)
	UpdateByPatientID(patientID int, newBody string) error
}

type UserRepository interface {
	SaveUser(user *entities.User) error
	UpdatePassword(id int, newPassword string) error

	// needed for login
	GetUser(email string, password string) (*entities.User, error)
	GetUserByID(ID int) (*entities.User, error)
	GetUserByProfileID(ID int) (*entities.User, error)
	GetAllDoctors() ([]*entities.User, error)
	GetDoctorsForSpecialty(specialty string) ([]*entities.User, error)

	// Maybe consider the following
	// 	UpdateCreate(user *entities.User) (*entities.User, error)
	// 	GetAll() ([]*entities.User, error)
}
type SessionRepository interface {
	GetSession(sessionToken string) *entities.Session
	DeleteSession(sessionToken string) error
	SaveSession(session *entities.Session) (*entities.Session, error)
}

type BlockchainClient interface {
	AddPermission(doctorHash, patientHash, permissionType string, message string) ([]byte, error)
	QueryPermissions(doctorHash string, patientHash string, message string) ([]entities.Permission, error)
}
