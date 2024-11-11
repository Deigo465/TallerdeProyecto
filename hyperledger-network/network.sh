#!/bin/bash -eu

function up(){
  docker compose -f docker-compose.yaml up -d
}
function down(){
  docker compose -f docker-compose.yaml down -v
}

# Function to setup the network and deploy chaincode
function setup() {
  CHANNEL_NAME=mychannel
  ORDERER_ADDRESS=orderer.example.com:7050
  ORDERER_CA=/etc/hyperledger/orderer/msp/tlscacerts/tlsca.example.com-cert.pem

  # Setup the crypto needed and create the channel genesis block and channel.tx
  echo "Generating crypto material"
  bin/cryptogen generate --config=./crypto-config.yaml
  mkdir -p config
  bin/configtxgen -profile TwoOrgsOrdererGenesis -channelID system-channel -outputBlock ./config/genesis.block
  bin/configtxgen -profile TwoOrgsChannel -outputCreateChannelTx ./config/channel.tx -channelID "$CHANNEL_NAME"

  docker compose -f docker-compose.yaml up -d

  echo "Joining channel and deploying chaincode"
  # Setup the org 1
  docker exec -it peer0.org1.example.com sh -c '
    export CORE_PEER_MSPCONFIGPATH=/etc/hyperledger/admin &&

    # Create channel
    echo "====================================" &&
    echo "Creating channel..." &&
    peer channel create -o orderer.example.com:7050 -c '"$CHANNEL_NAME"' -f /etc/hyperledger/config/channel.tx --outputBlock /etc/hyperledger/config/'"$CHANNEL_NAME"'.block --tls --cafile '"$ORDERER_CA"' &&
    
    echo "====================================" &&
    echo "Channel created successfully" &&
    # Join peer to channel
    peer channel join -b /etc/hyperledger/config/'"$CHANNEL_NAME"'.block --tls --cafile '"$ORDERER_CA"' &&
    
    echo "====================================" &&
    echo "Intalling chaincode..." &&
    # Install chaincode
    peer lifecycle chaincode install /etc/hyperledger/peercfg/basic.tar.gz &&
    
    echo "====================================" &&
    echo "Getting package ID" &&
    # Query installed chaincode to get package ID
    export PACKAGE_ID=$(peer lifecycle chaincode queryinstalled | awk -F "[, ]+" "/basic/{print \$3}") &&

    echo "====================================" &&
    echo "approving chaincode" &&
    # Approve chaincode for Org1
    peer lifecycle chaincode approveformyorg -o orderer.example.com:7050 --channelID '"$CHANNEL_NAME"' --name basic --version 1.0 --package-id $PACKAGE_ID --sequence 1 --tls --cafile '"$ORDERER_CA"'
  '

  # Setup for org 2
  docker exec -it peer0.org2.example.com sh -c '
    export CORE_PEER_MSPCONFIGPATH=/etc/hyperledger/admin
    
    echo "====================================" &&
    echo "Joining peer channel" &&
    # Join peer to channel
    peer channel join -b /etc/hyperledger/config/'"$CHANNEL_NAME"'.block --tls --cafile '"$ORDERER_CA"' &&
    
    echo "====================================" &&
    echo "Intalling chaincode..." &&
    # Install chaincode
    peer lifecycle chaincode install /etc/hyperledger/peercfg/basic.tar.gz &&
    
    # Query installed chaincode to get package ID
    export PACKAGE_ID=$(peer lifecycle chaincode queryinstalled | awk -F "[, ]+" "/basic/{print \$3}") &&

    echo "====================================" &&
    echo "approving chaincode" &&
    # Approve chaincode for Org2
    peer lifecycle chaincode approveformyorg -o orderer.example.com:7050 --channelID '"$CHANNEL_NAME"' --name basic --version 1.0 --package-id $PACKAGE_ID --sequence 1 --tls --cafile '"$ORDERER_CA"'
  '

  # Commit chaincode definition and invoke initialization from Org1
  docker exec -it peer0.org1.example.com sh -c '
    export CORE_PEER_MSPCONFIGPATH=/etc/hyperledger/admin &&
    
    echo "====================================" &&
    echo "Checking commit readiness" &&
    # Check commit readiness
    peer lifecycle chaincode checkcommitreadiness --channelID '"$CHANNEL_NAME"' --name basic --version 1.0 --sequence 1 --output json &&
    
    echo "====================================" &&
    echo "Committing chaincode" &&
    # Commit chaincode definition
    peer lifecycle chaincode commit -o orderer.example.com:7050  --tls --cafile '"$ORDERER_CA"' --channelID '"$CHANNEL_NAME"' --name basic --version 1.0 --sequence 1 --peerAddresses peer0.org1.example.com:7051 --peerAddresses peer0.org2.example.com:9051 --tlsRootCertFiles /etc/hyperledger/fabric/tls/ca.crt --tlsRootCertFiles /etc/hyperledger/peer2/tls/ca.crt
  '
}

# Function to update the channel configuration
function configOrg1() {
  CHANNEL_NAME=mychannel
  ORDERER_ADDRESS=orderer.example.com:7050
  PEER_ID=peer0.org1.example.com
  PEER_PORT=7051
  ORG_ID=Org1MSP

  config
}

function configOrg2() {
  CHANNEL_NAME=mychannel
  ORDERER_ADDRESS=orderer.example.com:7050

  PEER_ID=peer0.org2.example.com
  PEER_PORT=9051
  ORG_ID=Org2MSP

  config
}

function config(){
  ORDERER_CA=/etc/hyperledger/orderer/msp/tlscacerts/tlsca.example.com-cert.pem

  # Fetch the latest configuration block
  docker exec -it $PEER_ID sh -c "
    export CORE_PEER_MSPCONFIGPATH=/etc/hyperledger/admin &&
    cd /etc/hyperledger/peercfg &&
    echo 'Fetching latest configuration block' &&
    peer channel fetch config config_block.pb -o $ORDERER_ADDRESS -c $CHANNEL_NAME --tls --cafile '"$ORDERER_CA"' &&
    echo 'Configuration block fetched successfully'
  "
  mkdir -p $BASEPATH

  # Decode the configuration block to JSON
  bin/configtxlator proto_decode --input ./peerCfg/config_block.pb --type common.Block --output "${BASEPATH}config_block.json"
  jq ".data.data[0].payload.data.config" "${BASEPATH}config_block.json" >  "${BASEPATH}config.json"

  # Modify config.json to include anchor peer
  jq ".channel_group.groups.Application.groups.${ORG_ID}.values.AnchorPeers |= 
    {\"mod_policy\": \"Admins\", 
    \"value\": {\"anchor_peers\": [{\"host\": \"$PEER_ID\", \"port\": $PEER_PORT}]}, 
    \"version\": \"0\"}" "${BASEPATH}config.json" > "${BASEPATH}modified_config.json"

  # Encode the original configuration JSON to protobuf
  bin/configtxlator proto_encode --input "${BASEPATH}config.json" --type common.Config --output "${BASEPATH}original_config.pb"

  # Encode the modified configuration JSON to protobuf
  bin/configtxlator proto_encode --input "${BASEPATH}modified_config.json" --type common.Config --output "${BASEPATH}modified_config.pb"

  # Compute the delta between the original and modified configurations
  bin/configtxlator compute_update --channel_id $CHANNEL_NAME --original "${BASEPATH}original_config.pb" --updated "${BASEPATH}modified_config.pb" --output "${BASEPATH}config_update.pb"

  # Decode the delta to JSON format
  bin/configtxlator proto_decode --input ${BASEPATH}config_update.pb --type common.ConfigUpdate --output ${BASEPATH}config_update.json

  # Wrap the delta in an envelope message
  echo '{"payload":{"header":{"channel_header":{"channel_id":"'$CHANNEL_NAME'", "type":2}},"data":{"config_update":'$(cat ${BASEPATH}config_update.json)'}}}' | jq . > "${BASEPATH}config_update_in_envelope.json"

  # Encode the envelope message to protobuf
  bin/configtxlator proto_encode --input ${BASEPATH}config_update_in_envelope.json --type common.Envelope --output ./peerCfg/config_update_in_envelope.pb

  # Submit the update to the channel
  docker exec -it $PEER_ID sh -c "
    export CORE_PEER_MSPCONFIGPATH=/etc/hyperledger/admin &&
    cd /etc/hyperledger/peercfg &&
    peer channel update -f config_update_in_envelope.pb -c $CHANNEL_NAME -o $ORDERER_ADDRESS --tls --cafile '"$ORDERER_CA"' 
  "

  echo "Anchor peer update completed successfully."
}

function invoke(){
  ORDERER_CA=/etc/hyperledger/orderer/msp/tlscacerts/tlsca.example.com-cert.pem
  CHANNEL_NAME=mychannel
  # Make a random query
  docker exec -it peer0.org1.example.com sh -c '
    export CORE_PEER_MSPCONFIGPATH=/etc/hyperledger/admin &&

    # Perform a sample chaincode invocation
    peer chaincode invoke -o orderer.example.com:7050 --tls --cafile '"$ORDERER_CA"' -C '"$CHANNEL_NAME"' -n basic -c '\''{"function": "HelloWorld", "Args":[]}'\'' --peerAddresses peer0.org1.example.com:7051 --peerAddresses peer0.org2.example.com:9051 --tlsRootCertFiles /etc/hyperledger/fabric/tls/ca.crt --tlsRootCertFiles /etc/hyperledger/peer2/tls/ca.crt
    # peer chaincode invoke -o orderer.example.com:7050 --tls --cafile '"$ORDERER_CA"' -C '"$CHANNEL_NAME"' -n basic -c '\''{"function": "AddPermission", "Args":["permission1","doctorHash1","patientHash1","2024-06-04","granted"]}'\'' --peerAddresses peer0.org1.example.com:7051 --peerAddresses peer0.org2.example.com:9051  --tlsRootCertFiles /etc/hyperledger/fabric/tls/ca.crt --tlsRootCertFiles /etc/hyperledger/peer2/tls/ca.crt
    # peer chaincode invoke -o orderer.example.com:7050 --tls --cafile '"$ORDERER_CA"' -C '"$CHANNEL_NAME"' -n basic -c '\''{"function": "GetPermission", "Args":["permission1"]}'\'' --peerAddresses peer0.org1.example.com:7051 --peerAddresses peer0.org2.example.com:9051  --tlsRootCertFiles /etc/hyperledger/fabric/tls/ca.crt --tlsRootCertFiles /etc/hyperledger/peer2/tls/ca.crt
  '
}

function teardown(){
  docker compose -f docker-compose.yaml down -v
  docker volume prune -f
  rm -rf crypto-config config config-artifacts
}

# Compile chaincode
function compile(){
  echo "Compiling chaincode"
  FABRIC_CFG_PATH_PATH=$PWD/peerCfg/
  export FABRIC_CFG_PATH="$FABRIC_CFG_PATH_PATH"
  bin/peer lifecycle chaincode package peerCfg/basic.tar.gz --path ../hyperledger-chaincode/ --lang golang --label basic
  echo "Chaincode compiled successfully"
}
function check_install(){
  # Check if the binaries are installed
  if [[ ! -f bin/cryptogen ]]; then
    echo "cryptogen not found. Will download binaries"
    curl -sSL https://bit.ly/2ysbOFE | sudo bash -s
    sudo mv fabric-samples/bin .
  fi

}

# Main function to invoke functions based on user input
function main() {
  BASEPATH="$PWD/config-artifacts/"

  if [[ $1 == "up" ]]; then
    up
  elif [[ $1 == "down" ]]; then
    down
  elif [[ $1 == "setup" ]]; then
    check_install
    setup
    configOrg1
    configOrg2
    sleep 2
    invoke
  elif [[ $1 == "compile" ]]; then
    compile
  elif [[ $1 == "config" ]]; then
    configOrg1
    configOrg2
  elif [[ $1 == "teardown" ]]; then
    teardown
  elif [[ $1 == "invoke" ]]; then
    invoke
  elif [[ $1 == "troubleshoot" ]]; then
    docker inspect -f '{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' peer0.org1.example.com
    docker inspect -f '{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' peer0.org2.example.com
    docker inspect -f '{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' orderer.example.com
  elif [[ $1 == "monitor" ]]; then
  # Run the monitor script
  ./monitordocker.sh
  else
    echo "Invalid command. Available commands are: up, down, setup, compile, config, teardown, invoke"
  fi
}

# Run the main function with the provided argument
main "$1"