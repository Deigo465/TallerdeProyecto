# Guide to Deploy BlockeHR

0. Prerequisites
1. Install the dependencies (Git, Go)
2. Compile chaincode
3. Cloud server instantiation and configuration
4. Run the deploy script
5. Start app on server


## 0. Prerequisites
In order to deploy this app you need WSL2 or a Linux machine with the following dependencies:

- a local machine with Windows + WSL2 or Linux
- Internet connection
- SSH configured locally 

## 1. Install the deployment dependencies
1. Install Git (to clone the repository)
2. Install Go  (to build the project)
2. Install cURL (to install Hyperledger Network binaries, also usually included in Linux and WSL2)
3. Install RSync (usually installed by default on Linux and WSL2)


## 2. Compile the chaincode 
On the local machine:

Navigate to the `hyperledger-network` directory

```bash
cd hyperledger-network
```

Install the Hyperledger Fabric binaries
```bash
curl -sSL https://bit.ly/2ysbOFE | sudo bash -s
sudo mv fabric-samples/bin .
# Taken from https://hyperledger-fabric.readthedocs.io/es/latest/install.html
```

Compile the chaincode
```bash
./network-sh compile
```

## 3. Cloud server instantiation and configuration

1. Create a VPS on Digital Ocean or any other cloud provider
2. For the OS, choose the Ubuntu 24.10 image or similar.
3. For the Hardware, choose at least 2GB of RAM and 2 vCPUs, and min 25 GB of storage.
4. Setup SSH authentication (or password authentication) (if the cloud provider does not provide a console)
5. Log in to the server using SSH `ssh root@<server-ip>`
6. Setup the folders
```bash
    sudo mkdir /var/www
    sudo chown -R $USER:$USER /var/www
    mkdir /var/www/pkg/
    mkdir -p /var/www/hyperledger-network/config
```
7. Install Docker on the server and docker-compose
Install docker on the virtual machine (Ubuntu)
https://docs.docker.com/engine/install/ubuntu/

```bash
# Add Docker's official GPG key:
sudo apt-get update
sudo apt-get install ca-certificates curl
sudo install -m 0755 -d /etc/apt/keyrings
sudo curl -fsSL https://download.docker.com/linux/ubuntu/gpg -o /etc/apt/keyrings/docker.asc
sudo chmod a+r /etc/apt/keyrings/docker.asc

# Add the repository to Apt sources:
echo \
  "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.asc] https://download.docker.com/linux/ubuntu \
  $(. /etc/os-release && echo "$VERSION_CODENAME") stable" | \
  sudo tee /etc/apt/sources.list.d/docker.list > /dev/null
sudo apt-get update

sudo apt-get install -y docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin
```

Add to /etc/hosts file the following lines

```bash
127.0.0.1 orderer.example.com
127.0.0.1 peer0.org1.example.com
127.0.0.1 peer0.org2.example.com
```


## 4. Run the deploy script

Replace line `12` with the correct server IP address

```bash
# e.g.
SERVER="root@165.232.151.138";
```

Then run the script with

```bash
./deploy/deploy-script initial
```

For subsequent deploys you can run:

```bash
./deploy/deploy-script
```

## 5. Start the app on the server

On the local machine run:

```bash
./deploy/deploy-script blockchain setup
```
You can test the blockchain by calling the invoke method:

```bash
./deploy/deploy-script blockchain invoke
# e.g => Hello World! blockehr
```

Monitor blockchain with:

```bash
./deploy/deploy-script blockchain monitor
```


On the server, start the app with TMUX for persistent sessions

```bash
tmux new
```

Then run the app

```bash
cd /var/www/
./blockehr-app-linux-amd64
```

You should be able to navigate to the server IP address in port :3000 to see the app running.

Then, to detach from the tmux session press `Ctrl+b` then `d`

To reattach to the tmux session run

```bash
tmux attach
```