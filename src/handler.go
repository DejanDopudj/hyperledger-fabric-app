package src

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

type Handler struct {
	Validator *validator.Validate
}

func (handler *Handler) GetUser(writer http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"]

	if id == "" {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	//TODO: check if user exists
	result, err := mockResponse("GetUser")

	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	responseOK(&writer, result)
}

func (handler *Handler) RegisterUser(writer http.ResponseWriter, req *http.Request) {
	var userDto RegisterUser
	err := json.NewDecoder(req.Body).Decode(&userDto)

	if err != nil || handler.Validator.Struct(userDto) != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	//TODO: call fabcar
	result, err := mockResponse("CreateUser")

	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	responseOK(&writer, result)
}

func (handler *Handler) CreateBankAccount(writer http.ResponseWriter, req *http.Request) {
	var accountDto CreateBankAccount
	err := json.NewDecoder(req.Body).Decode(&accountDto)

	if err != nil || handler.Validator.Struct(accountDto) != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	//TODO: HEHE
	result, err := mockResponse("CreateAccount")

	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	responseOK(&writer, result)
}

func (handler *Handler) TransferAssets(writer http.ResponseWriter, req *http.Request) {
	var transferDto CreateBankAccount
	err := json.NewDecoder(req.Body).Decode(&transferDto)

	if err != nil || handler.Validator.Struct(transferDto) != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	//TODO: HEHE
	result, err := mockResponse("TransferAssets")

	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	responseOK(&writer, result)

}

func (handler *Handler) AddAssets(writer http.ResponseWriter, req *http.Request) {
	var fundsDto ManageFunds
	err := json.NewDecoder(req.Body).Decode(&fundsDto)

	if err != nil || handler.Validator.Struct(fundsDto) != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	result, err := mockResponse("AddAssets")

	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	responseOK(&writer, result)
}

func (handler *Handler) PullAssets(writer http.ResponseWriter, req *http.Request) {
	var fundsDto ManageFunds
	err := json.NewDecoder(req.Body).Decode(&fundsDto)

	if err != nil || handler.Validator.Struct(fundsDto) != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	result, err := mockResponse("PullAssets")

	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	responseOK(&writer, result)
}

func (handler *Handler) GetBanks(writer http.ResponseWriter, req *http.Request) {
	//TODO: Call fabcar
	result, err := mockResponse("QueryAllBanks")

	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	responseOK(&writer, result)
}

func (handler *Handler) GetBank(writer http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"]

	log.Println(id)

	//TODO: Call fabcar
	result, err := mockResponse("QueryBank")

	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	responseOK(&writer, result)
}

func (handler *Handler) GetUsers(writer http.ResponseWriter, req *http.Request) {
	//TODO: Call fabcar
	result, err := mockResponse("QueryAllUsers")

	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	responseOK(&writer, result)
}

func (handler *Handler) GetBankAccounts(writer http.ResponseWriter, req *http.Request) {
	//TODO: Call fabcar
	result, err := mockResponse("QueryAllAccounts")

	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	responseOK(&writer, result)
}

type Result struct {
}

func (handler *Handler) GetBankAccount(writer http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"]

	log.Println(id)

	//TODO: Call fabcar
	result, err := mockResponse("QueryAccount")

	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	writer.WriteHeader(http.StatusOK)
	writer.Header().Set("Content-Type", "application/json")

	json.NewEncoder(writer).Encode(result)
}

func mockResponse(call string) ([]byte, error) {
	return json.Marshal(Result{})
}

func responseOK(writer *http.ResponseWriter, result []byte) {
	(*writer).WriteHeader(http.StatusOK)
	(*writer).Header().Set("Content-Type", "application/json")

	json.NewEncoder((*writer)).Encode(result)
}
