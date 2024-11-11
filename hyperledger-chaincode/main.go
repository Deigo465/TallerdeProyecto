/*
SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"log"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/open-wm/blockehr/hyperledger-chaincode/chaincode"
)

func main() {
	permissionsChaincode, err := contractapi.NewChaincode(&chaincode.SmartContract{})
	if err != nil {
		log.Panicf("Error creating permissions chaincode: %v", err)
	}

	if err := permissionsChaincode.Start(); err != nil {
		log.Panicf("Error starting permissions chaincode: %v", err)
	}
}
