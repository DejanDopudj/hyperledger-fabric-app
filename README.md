# Bankchain
This application supports user management, bank account creation, various querying, bank management, money withdrawal, and money transfer functionalities.

## Prerequisites
- Hyperledger Fabric 2.2.6
- Go 1.21.6
- Docker 
- Docker Compose

## Usage
1. Install dependencies:
```sh
./install-fabric.sh -f 2.2.6
```
2. Enter the `test-network`:
```sh
cd test-network
```
3. Down any previously created containers:
```sh
./network.sh down
```
4. Create components and start the containers:
```sh
./network.sh up
```
5. Create a channel:
```sh
./network.sh createChannel -ca -c channel-name
```
6. Deploy chaincode: 
```sh
./network.sh deployCC -ccn basic -ccp ../chaincode/fabcar/go/ -ccl go;
```
7. Start the server:
```sh
cd ../fabcar/go
./runfabcar.sh
```

## Authors
- Dejan Dopuđ E2 2/2023 <br>
- Sanja Petrović E2 4/2023 <br>
- Filip Milošević E2 108/2023