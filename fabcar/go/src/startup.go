package src

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func StartServer(handler *Handler) {
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/enroll", handler.Enroll).Methods("POST")
	router.HandleFunc("/users/{user}", handler.GetUser).Methods("GET")
	router.HandleFunc("/bank/{bank}", handler.GetBank).Methods("GET")
	router.HandleFunc("/accounts/{account}", handler.GetBankAccount).Methods("GET")

	router.HandleFunc("/accounts", handler.CreateBankAccount).Methods("POST")
	router.HandleFunc("/assets/add", handler.AddAssets).Methods("POST")
	router.HandleFunc("/assets/pull", handler.PullAssets).Methods("POST")
	router.HandleFunc("/assets/transfer", handler.TransferAssets).Methods("POST")

	println("Server started")
	log.Fatal(http.ListenAndServe(":3000", router))
}
