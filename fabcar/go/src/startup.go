package src

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/hyperledger/fabric-sdk-go/pkg/gateway"
)

func StartServer(handler *Handler) {
	initLedger(handler.Wallet)

	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/enroll", handler.Enroll).Methods("POST")
	router.HandleFunc("/users/{user}", handler.GetUser).Methods("GET")
	router.HandleFunc("/bank/{bank}", handler.GetBank).Methods("GET")
	router.HandleFunc("/accounts/{account}", handler.GetBankAccount).Methods("GET")
	router.HandleFunc("/accounts", handler.CreateBankAccount).Methods("POST")
	router.HandleFunc("/assets/payment", handler.MakePayment).Methods("POST")
	router.HandleFunc("/assets/withdrawal", handler.MakeWithdrawal).Methods("POST")
	router.HandleFunc("/assets/transfer", handler.TransferAssets).Methods("POST")

	router.HandleFunc("/banks", handler.QueryAllBanks).Methods("GET")

	println("Server started")
	log.Fatal(http.ListenAndServe(":3000", router))
}

func initLedger(wallet *gateway.Wallet) {
	if !wallet.Exists("app") {
		err := populateWallet(wallet, "1", "app")
		if err != nil {
			log.Panicf("Failed to populate wallet contents: %s\n", err.Error())
			os.Exit(1)
		}
	}

	contract, err := getContract(wallet, "1", "app")

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
