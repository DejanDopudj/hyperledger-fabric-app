# hyperledger-fabric-app

## Running instructions
1. Position yourself inside the `artifacts/channel` folder:
```sh
cd artifacts/channel
```
2. Generate configuration files:
```sh
bash run_cryptogen.sh
```
```sh
bash run_configtx.sh
```
```sh
bash ccp-generate.sh
```
3. Start the containers:
```sh
docker compose up
```
### Note
Please paste the "bin" folder from fabric-samples into the root of the project to be able to run the `run_cryptogen.sh` and `run_configtx.sh` scripts.
