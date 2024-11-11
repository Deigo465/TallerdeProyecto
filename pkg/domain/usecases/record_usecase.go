package usecase

import (
	"errors"
	"fmt"
	"log"

	"github.com/open-wm/blockehr/pkg/domain/entities"
	"github.com/open-wm/blockehr/pkg/interfaces"
)

type RecordUsecase interface {
	Add(actor *entities.User, record *entities.Record) error
	GetAllForPatient(actor *entities.User, patientId int) ([]*entities.Record, error)
	GetById(actor *entities.User, id int) (*entities.Record, error)
	UpdateByPatientId(actor *entities.User, id int, body string) error
}

type recordUsecase struct {
	profileRepo interfaces.ProfileRepository
	recordRepo  interfaces.RecordRepository
	filesRepo   interfaces.FileRepository
	client      interfaces.BlockchainClient
}

func NewRecordUsecase(repo interfaces.RecordRepository, profileRepo interfaces.ProfileRepository, filesRepo interfaces.FileRepository, client interfaces.BlockchainClient) RecordUsecase {
	return &recordUsecase{
		recordRepo:  repo,
		profileRepo: profileRepo,
		filesRepo:   filesRepo,
		client:      client,
	}
}

func (uc *recordUsecase) Add(actor *entities.User, record *entities.Record) error {
	// check if role of actor is doctor

	if actor.Profile.Role != "DOCTOR" {
		return errors.New("only doctor can add records")
	}
	// check if the fields are not empty

	if record.Body == "" {
		return errors.New("Expecting body to not be empty")
	}
	if record.CreatedAt == "" {
		return errors.New("Expecting created at to not be empty")
	}
	if record.UpdatedAt == "" {
		return errors.New("Expecting update at to not be empty")
	}
	if record.PatientId == 0 {
		return errors.New("Expecting patient id to not be empty")
	}
	if record.DoctorId == 0 {
		return errors.New("Expecting doctor id to not be empty")
	}
	record.Specialty = actor.Profile.Specialty

	newRecord, err := uc.recordRepo.Add(record)
	if err != nil {
		return err
	}

	for _, file := range record.Files {
		if file == nil {
			continue
		}
		log.Printf("Adding file: %s with record id %d", file.Name, record.ID)
		file.RecordId = newRecord.ID
		err := uc.filesRepo.Add(file)
		if err != nil {
			return err
		}
	}

	return nil
}

func checkIfGranted(actor *entities.User, patientId int, uc *recordUsecase) (bool, error) {
	// check on the blockchain if the doctor has permission to access the patient's records
	doctor := fmt.Sprintf("doctor%d", actor.ID)
	patient := fmt.Sprintf("patient%d", patientId)
	message := fmt.Sprintf("%s intentó ver la HC  del paciente: %d", actor.Profile.FirstName+" "+actor.Profile.FatherLastName, patientId)
	permissions, err := uc.client.QueryPermissions(doctor, patient, message)

	if err != nil {
		return false, err
	}

	latest := entities.Permission{Type: "granted"}

	for _, permission := range permissions {
		if permission.CreatedAt > latest.CreatedAt {
			latest = permission
		}
	}

	if latest.Type != "granted" {
		return false, errors.New("doctor does not have permission to access patient records")
	}

	return true, nil
}

func (uc *recordUsecase) GetAllForPatient(actor *entities.User, patientId int) ([]*entities.Record, error) {
	// Verificar si el actor es un doctor
	if actor.Profile.Role != "DOCTOR" {
		return nil, errors.New("only doctor can get records")
	}

	if patientId == 0 {
		return nil, errors.New("patient is nil")
	}

	granted, err := checkIfGranted(actor, patientId, uc)
	if err != nil || !granted {
		return nil, err
	}

	records, err := uc.recordRepo.GetAll()
	if err != nil {
		return nil, err
	}

	var patientRecords []*entities.Record
	for _, record := range records {
		if record.PatientId == patientId {
			// get doctor info as well
			doctor, err := uc.profileRepo.GetById(record.DoctorId)
			if err != nil {
				return nil, err
			}

			record.Doctor = doctor
			patientRecords = append(patientRecords, record)
		}
	}

	return patientRecords, nil
}
func (uc *recordUsecase) GetById(actor *entities.User, id int) (*entities.Record, error) {
	// check if role of actor is doctor

	if actor.Profile.Role != "DOCTOR" {
		return nil, errors.New("only doctor can get records")
	}
	record, err := uc.recordRepo.GetById(id)
	if err != nil {

		return nil, err
	}
	return record, nil
}

func (uc *recordUsecase) UpdateByPatientId(actor *entities.User, id int, body string) error {
	// check if role of actor is doctor

	if actor.Profile.Role != "DOCTOR" {
		return errors.New("only doctor can get records")
	}
	err := uc.recordRepo.UpdateByPatientID(id, body)
	if err != nil {

		return err
	}
	return nil
}
