package infrastructure

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/hyperledger/fabric-gateway/pkg/client"
	"github.com/open-wm/blockehr/pkg/domain/entities"
	"github.com/open-wm/blockehr/pkg/infrastructure/hyperledger"
)

const (
	mspID       = "Org1MSP"
	cryptoPath  = "hyperledger-network/crypto-config/peerOrganizations/org1.example.com"
	certPath    = cryptoPath + "/users/Admin@org1.example.com/msp/signcerts"
	keyPath     = cryptoPath + "/users/Admin@org1.example.com/msp/keystore"
	tlsCertPath = cryptoPath + "/peers/peer0.org1.example.com/tls/ca.crt"

	peerEndpoint = "localhost:7051"
	gatewayPeer  = "peer0.org1.example.com"
)

type Permission struct {
	CreatedAt string `json:"created_at"`
	Type      string `json:"type"` // type can be "granted" or "revoked"
}

type DoctorPatient struct {
	ID          string                `json:"id"` // this is the doctor_hash+patient_hash
	Permissions []entities.Permission `json:"permissions"`
}

type Blockchain struct {
	client *client.Contract
}

func NewBlockchain(basePath string) *Blockchain {
	contract := hyperledger.SetupClient(mspID,
		basePath+certPath,
		basePath+keyPath,
		basePath+tlsCertPath,
		peerEndpoint, gatewayPeer)
	return &Blockchain{client: contract}
}

func (b *Blockchain) HelloWorld() {
	fmt.Println("\n--> Evaluate Transaction: HelloWorld, function returns 'Hello, World!'")

	evaluateResult, err := b.client.SubmitTransaction("HelloWorld")
	if err != nil {
		panic(fmt.Errorf("failed to evaluate transaction: %w", err))
	}

	fmt.Printf("*** Result: %s\n", evaluateResult)
}

func (b *Blockchain) AddPermission(doctorHash, patientHash, permissionType, message string) ([]byte, error) {
	fmt.Println("\n--> Submit Transaction: AddPermission, function adds a new permission to the ledger")

	now := time.Now()
	res, err := b.client.SubmitTransaction("AddPermission", doctorHash, patientHash, now.Local().String(), permissionType, message)
	if err != nil {
		log.Println(fmt.Errorf("failed to submit transaction: %w", err))
		hyperledger.ExampleErrorHandling(err)
	}

	return res, nil
}

func (b *Blockchain) QueryPermissions(doctorHash, patientHash, message string) ([]entities.Permission, error) {
	evaluateResult, err := b.client.EvaluateTransaction("GetPermissions", doctorHash, patientHash, message)
	if err != nil {
		log.Println(fmt.Errorf("failed to evaluate transaction: %w", err))
	}
	doctorPatient := &DoctorPatient{}

	err = json.Unmarshal(evaluateResult, &doctorPatient)
	if err != nil {
		return nil, err
	}

	return doctorPatient.Permissions, nil
}
