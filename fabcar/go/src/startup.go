package src

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/hyperledger/fabric-sdk-go/pkg/gateway"
)

func StartServer(handler *Handler) {
	initWallets(handler.Wallet)
	initLedger(handler.Wallet)

	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/users", handler.CreateUser).Methods("POST")
	router.HandleFunc("/users/{user}", handler.GetUser).Methods("GET")
	router.HandleFunc("/banks/{bank}", handler.GetBank).Methods("GET")
	router.HandleFunc("/accounts/{account}", handler.GetBankAccount).Methods("GET")
	router.HandleFunc("/accounts", handler.CreateBankAccount).Methods("POST")
	router.HandleFunc("/assets/payment", handler.MakePayment).Methods("POST")
	router.HandleFunc("/assets/withdrawal", handler.MakeWithdrawal).Methods("POST")
	router.HandleFunc("/assets/transfer", handler.TransferAssets).Methods("POST")

	router.HandleFunc("/banks", handler.QueryAllBanks).Methods("GET")
	router.HandleFunc("/accounts", handler.QueryAllAccounts).Methods("GET")
	router.HandleFunc("/users", handler.QueryAllUsers).Methods("GET")

	router.HandleFunc("/bank/full", handler.QueryBankFull).Methods("GET")
	router.HandleFunc("/account/full", handler.QueryAccountFull).Methods("GET")
	router.HandleFunc("/user/full", handler.QueryUserFull).Methods("GET")

	println("Server started")
	log.Fatal(http.ListenAndServe(":3000", router))
}

func initWallets(wallet *gateway.Wallet) {
	if !wallet.Exists("org1") {
		err := populateWallet(wallet, "1")
		if err != nil {
			log.Panicf("Failed to populate wallet contents: %s\n", err.Error())
			os.Exit(1)
		}
	}

	if !wallet.Exists("org2") {
		err := populateWallet(wallet, "2")
		if err != nil {
			log.Panicf("Failed to populate wallet contents: %s\n", err.Error())
			os.Exit(1)
		}
	}

	if !wallet.Exists("org3") {
		err := populateWallet(wallet, "3")
		if err != nil {
			log.Panicf("Failed to populate wallet contents: %s\n", err.Error())
			os.Exit(1)
		}
	}

	if !wallet.Exists("org4") {
		err := populateWallet(wallet, "4")
		if err != nil {
			log.Panicf("Failed to populate wallet contents: %s\n", err.Error())
			os.Exit(1)
		}
	}

}

func initLedger(wallet *gateway.Wallet) {
	contract, err := getContract(wallet, "1")

	if err != nil {
		wallet.Remove("app")
		log.Fatal("Failed to get contract")
		os.Exit(1)
	}

	_, err = contract.SubmitTransaction("InitLedger")

	if err != nil {
		wallet.Remove("app")
		log.Fatal("Failed to Init Ledger")
		os.Exit(1)
	}

	log.Println("Init Ledger Completed!")
}
