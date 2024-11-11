package chaincode

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type smartContract interface {
	CreatePermission(ctx contractapi.TransactionContextInterface, id string, doctorHash string, patientHash string, createdAt string, permissionType string) error
	GetPermission(ctx contractapi.TransactionContextInterface, id string) (*Permission, error)
	readState(ctx contractapi.TransactionContextInterface, id string) ([]byte, error)
}

// SmartContract provides functions for managing an Asset
type SmartContract struct {
	contractapi.Contract
}

type Permission struct {
	CreatedAt string `json:"created_at"`
	Type      string `json:"type"` // type can be "granted" or "revoked"
}

type DoctorPatient struct {
	ID          string       `json:"id"` // this is the doctor_hash+patient_hash
	Permissions []Permission `json:"permissions"`
}

// HelloWorld - returns a string
func (sc *SmartContract) HelloWorld(ctx contractapi.TransactionContextInterface) string {
	return "Hello World! blockehr"
}

func (s *SmartContract) readState(ctx contractapi.TransactionContextInterface, id string) ([]byte, error) {
	permissionJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %w", err)
	}
	if permissionJSON == nil {
		return nil, fmt.Errorf("the permission %s does not exist", id)
	}

	return permissionJSON, nil
}

func (s *SmartContract) AddPermission(ctx contractapi.TransactionContextInterface, doctorHash string, patientHash string, createdAt string, permissionType string, message string) error {
	id := doctorHash + patientHash
	existing, err := s.readState(ctx, id)
	var doctorPatient DoctorPatient
	if err == nil && existing != nil {
		err = json.Unmarshal(existing, &doctorPatient)
		if err != nil {
			return err
		}
	} else {
		doctorPatient = DoctorPatient{
			ID:          id,
			Permissions: []Permission{},
		}
	}

	permission := Permission{
		CreatedAt: createdAt,
		Type:      permissionType,
	}

	doctorPatient.Permissions = append(doctorPatient.Permissions, permission)

	doctorPatientJSON, err := json.Marshal(doctorPatient)
	if err != nil {
		return err
	}

	log.Println(message)
	ctx.GetStub().SetEvent("CreatePermission", doctorPatientJSON)
	return ctx.GetStub().PutState(id, doctorPatientJSON)
}

func (s *SmartContract) GetPermissions(ctx contractapi.TransactionContextInterface, doctorHash string, patientHash string, message string) (*DoctorPatient, error) {
	id := doctorHash + patientHash
	doctorPatientJSON, err := s.readState(ctx, id)
	if err != nil {
		log.Println("error getting permission")
		return nil, err
	}

	var doctorPatient DoctorPatient
	err = json.Unmarshal(doctorPatientJSON, &doctorPatient)
	if err != nil {
		log.Println("error unmarshalling permission")
		return nil, err
	}

	log.Println(message)
	return &doctorPatient, nil
}
