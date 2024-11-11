package mock_repositories

import (
	"github.com/open-wm/blockehr/pkg/domain/entities"
	"github.com/open-wm/blockehr/pkg/interfaces"
)

type InMemoryBlockchain struct {
}

func NewInMemoryBlockchain() interfaces.BlockchainClient {
	return &InMemoryBlockchain{}
}

// AddPermission implements interfaces.BlockchainClient.
func (i *InMemoryBlockchain) AddPermission(doctorHash string, patientHash string, permissionType string, message string) ([]byte, error) {
	return nil, nil
}

// QueryPermissions implements interfaces.BlockchainClient.
func (i *InMemoryBlockchain) QueryPermissions(doctorHash string, patientHash string, message string) ([]entities.Permission, error) {
	permission := entities.Permission{
		CreatedAt: "2021-09-01",
		Type:      "granted",
	}
	return []entities.Permission{permission}, nil
}
