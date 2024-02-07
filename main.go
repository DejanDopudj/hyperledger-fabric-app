package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gorilla/mux"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/gateway"
	"github.com/joho/godotenv"
)

type Handler struct {
	contract *gateway.Contract
}

func (handler *Handler) CreateUser(writer http.ResponseWriter, req *http.Request) {
	writer.WriteHeader(http.StatusOK)
	writer.Header().Set("Content-Type", "application/json")
}

func (handler *Handler) CreateBankAccount(writer http.ResponseWriter, req *http.Request) {
	writer.WriteHeader(http.StatusOK)
	writer.Header().Set("Content-Type", "application/json")
}

func (handler *Handler) TransferAssets(writer http.ResponseWriter, req *http.Request) {
	writer.WriteHeader(http.StatusOK)
	writer.Header().Set("Content-Type", "application/json")
}

func (handler *Handler) AddAssets(writer http.ResponseWriter, req *http.Request) {
	writer.WriteHeader(http.StatusOK)
	writer.Header().Set("Content-Type", "application/json")
}

func (handler *Handler) PullAssets(writer http.ResponseWriter, req *http.Request) {
	writer.WriteHeader(http.StatusOK)
	writer.Header().Set("Content-Type", "application/json")
}

func (handler *Handler) GetBanks(writer http.ResponseWriter, req *http.Request) {
	writer.WriteHeader(http.StatusOK)
	writer.Header().Set("Content-Type", "application/json")
}

func (handler *Handler) GetUsers(writer http.ResponseWriter, req *http.Request) {
	writer.WriteHeader(http.StatusOK)
	writer.Header().Set("Content-Type", "application/json")
}

func (handler *Handler) GetBankAccounts(writer http.ResponseWriter, req *http.Request) {
	writer.WriteHeader(http.StatusOK)
	writer.Header().Set("Content-Type", "application/json")
}

func startServer(handler *Handler) {
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/accounts", handler.GetBankAccounts).Methods("GET")
	router.HandleFunc("/accounts", handler.CreateBankAccount).Methods("POST")
	router.HandleFunc("/assets", handler.AddAssets).Methods("POST")
	router.HandleFunc("/assets", handler.PullAssets).Methods("PUT")
	router.HandleFunc("/assets/transfer", handler.TransferAssets).Methods("POST")
	router.HandleFunc("/banks", handler.GetBanks).Methods("GET")
	router.HandleFunc("/users", handler.GetUsers).Methods("GET")
	router.HandleFunc("/users", handler.CreateUser).Methods("POST")

	println("Server started")
	log.Fatal(http.ListenAndServe(os.Getenv("PORT"), router))
}

func initApplication() *gateway.Contract {
	err := os.Setenv("DISCOVERY_AS_LOCALHOST", "true")
	if err != nil {
		log.Fatalf("Error setting DISCOVERY_AS_LOCALHOST environemnt variable: %v", err)
	}

	wallet, err := gateway.NewFileSystemWallet("wallet")
	if err != nil {
		log.Fatalf("Failed to create wallet: %v", err)
	}

	if !wallet.Exists("appUser") {
		//err = populateWallet(wallet)
		//if err != nil {
		//	log.Fatalf("Failed to populate wallet contents: %v", err)
		//}
	}

	ccpPath := filepath.Join(
		"artifacts",
		"channel",
		"crypto-config",
		"peerOrganizations",
		"org1.example.com",
		//TODO: Create connection file for every organization
		"connection-org1.yaml",
	)

	gw, err := gateway.Connect(
		gateway.WithConfig(config.FromFile(filepath.Clean(ccpPath))),
		gateway.WithIdentity(wallet, "appUser"),
	)
	if err != nil {
		log.Fatalf("Failed to connect to gateway: %v", err)
	}
	defer gw.Close()

	network, err := gw.GetNetwork("mychannel")
	if err != nil {
		log.Fatalf("Failed to get network: %v", err)
	}

	return network.GetContract("basic")
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	startServer(&Handler{
		contract: initApplication(),
	})
}
