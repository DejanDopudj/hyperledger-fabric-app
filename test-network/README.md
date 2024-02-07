# hyperledger-fabric-app

## Running instructions
1. Install dependencies
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