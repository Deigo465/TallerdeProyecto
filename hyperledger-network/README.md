# Hyperledger Fabric Network for BlockeHR
To run the is project you need `docker` server installed in your system.
- Docker (tested on v26.1)

## Get started
Install `docker` and `docker compose`, and Hypereldger fabric binaries

```bash
brew install --cask docker
brew install docker-compose
```

```bash
sudo ./network up
```

Compile the chaincode
```bash
./network.sh compile # compile the chaincode
```
This command will do the following:
1. Generate the crypto material for the network

Start the newtork 
```bash
./network.sh setup
```

This command will do the following:
1. Create the artifacts for the network (genesis block and channel.tx)
2. Start the network
3. Create and join the channel on Org1 and Org2
4. Install and instantiate the chaincode on the channel
5. Commit the chaincode on the channel


The blockchain for this network is found on the hyperledger-chaincode directory


To start again the network (after a correct setup)
```bash
./network.sh up
```

To shut down the network (after a correct setup)
```bash
./network.sh down
```

To check the network logs run
```bash
 ./monitordocker.sh
```

To shutdown and remove the containers run:
```bash
./network.sh teardown
```

To test the blockchain you can run the following command:
```bash
./network.sh invoke
```

## Troubleshooting the client

Run the following command to get the ip addresses of the containers:
```bash
docker inspect -f '{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' peer0.org1.example.com
docker inspect -f '{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' peer0.org2.example.com
docker inspect -f '{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' orderer.example.com
```
Example output:
```
172.28.0.4
172.28.0.2
172.28.0.3
```


Then add these IP addresses to your /etc/hosts file (or equivalent) with the following format:

```
127.0.0.1 orderer.example.com
127.0.0.1 peer0.org1.example.com
127.0.0.1 peer0.org2.example.com

# 172.28.0.4 peer0.org1.example.com
# 172.28.0.3 peer0.org2.example.com
# 172.28.0.2 orderer.example.com
```



## Blockchain architecture:

- Org1:
  - Peer0
- Org2:
  - Peer0
- Orderer:


## Acknowledgements
- [Fabric Test Network](https://hyperledger-fabric.readthedocs.io/en/latest/test_network.html)
- [Fabric Deployment guid](https://hyperledger-fabric.readthedocs.io/en/latest/deployment_guide_overview.html#step-one-decide-on-your-network-configuration)
- [Fabric Learn repo - Chinese](https://github-com.translate.goog/wefantasy/FabricLearn?_x_tr_sl=auto&_x_tr_tl=es&_x_tr_hl=es&_x_tr_pto=wapp)
- [Fabric Gateway](https://github.com/hyperledger/fabric-gateway/blob/main/pkg/client/README.md)
