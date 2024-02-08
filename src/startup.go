package src

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func StartServer(handler *Handler) {
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/accounts", handler.GetBankAccounts).Methods("GET")
	router.HandleFunc("/accounts/{id}", handler.GetBankAccount).Methods("GET")
	router.HandleFunc("/accounts", handler.CreateBankAccount).Methods("POST")
	router.HandleFunc("/assets/add", handler.AddAssets).Methods("POST")
	router.HandleFunc("/assets/pull", handler.PullAssets).Methods("POST")
	router.HandleFunc("/assets/transfer", handler.TransferAssets).Methods("POST")
	router.HandleFunc("/banks", handler.GetBanks).Methods("GET")
	router.HandleFunc("/banks/{id}", handler.GetBank).Methods("GET")
	router.HandleFunc("/login/{id}", handler.GetUser).Methods("GET")
	router.HandleFunc("/users", handler.RegisterUser).Methods("POST")

	println("Server started")
	log.Fatal(http.ListenAndServe(os.Getenv("PORT"), router))
}
