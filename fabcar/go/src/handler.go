package src

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/hyperledger/fabric-sdk-go/pkg/gateway"
)

type ErrorMessage struct {
	message string
	code    int
}

type Handler struct {
	Validator *validator.Validate
	Wallet    *gateway.Wallet
}

func (handler *Handler) Enroll(writer http.ResponseWriter, req *http.Request) {
	var userDto RegisterUser
	err := json.NewDecoder(req.Body).Decode(&userDto)

	if err != nil || handler.Validator.Struct(userDto) != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	err = enrollUser(handler.Wallet, userDto)

	if err != nil {
		responseError(&writer, err.Error(), http.StatusBadRequest)
		return
	}

	writer.WriteHeader(http.StatusCreated)
}

func (handler *Handler) GetUser(writer http.ResponseWriter, req *http.Request) {

	user := mux.Vars(req)["user"]

	if user == "" {
		responseError(&writer, fmt.Sprintf("User id is required"), http.StatusBadRequest)
		return
	}

	bank := req.URL.Query().Get("bank")

	if bank == "" {
		responseError(&writer, fmt.Sprintf("Bank id is required"), http.StatusBadRequest)
		return
	}

	contract, err := getContract(handler.Wallet, bank, user)

	if err != nil {
		responseError(&writer, err.Error(), http.StatusBadRequest)
		return
	}

	response, err := contract.EvaluateTransaction("QueryUser", user)

	if err != nil {
		responseError(&writer, err.Error(), http.StatusBadRequest)
		return
	}

	responseOK(&writer, response)
}

func (handler *Handler) GetBank(writer http.ResponseWriter, req *http.Request) {
	bank := mux.Vars(req)["bank"]

	if bank == "" {
		responseError(&writer, fmt.Sprintf("Bank id is required"), http.StatusBadRequest)
		return
	}

	user := req.URL.Query().Get("user")

	if user == "" {
		responseError(&writer, fmt.Sprintf("User id is required"), http.StatusBadRequest)
		return
	}

	contract, err := getContract(handler.Wallet, bank, user)

	if err != nil {
		responseError(&writer, err.Error(), http.StatusBadRequest)
		return
	}

	response, err := contract.EvaluateTransaction("QueryBank", bank)

	if err != nil {
		responseError(&writer, err.Error(), http.StatusBadRequest)
		return
	}

	responseOK(&writer, response)
}

func (handler *Handler) GetBankAccount(writer http.ResponseWriter, req *http.Request) {
	bank := req.URL.Query().Get("bank")

	if bank == "" {
		responseError(&writer, fmt.Sprintf("Bank id is required"), http.StatusBadRequest)
		return
	}

	user := req.URL.Query().Get("user")

	if user == "" {
		responseError(&writer, fmt.Sprintf("User id is required"), http.StatusBadRequest)
		return
	}

	account := mux.Vars(req)["account"]

	if user == "" {
		responseError(&writer, fmt.Sprintf("Account id is required"), http.StatusBadRequest)
		return
	}

	contract, err := getContract(handler.Wallet, bank, user)

	if err != nil {
		responseError(&writer, err.Error(), http.StatusBadRequest)
		return
	}

	response, err := contract.EvaluateTransaction("QueryAccount", account)

	if err != nil {
		responseError(&writer, err.Error(), http.StatusBadRequest)
		return
	}

	responseOK(&writer, response)

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

func mockResponse(call string) ([]byte, error) {
	return json.Marshal(ErrorMessage{})
}

func responseOK(writer *http.ResponseWriter, result []byte) {
	(*writer).WriteHeader(http.StatusOK)
	(*writer).Header().Set("Content-Type", "application/json")

	json.NewEncoder((*writer)).Encode(result)
}

func responseError(writer *http.ResponseWriter, message string, code int) {
	(*writer).WriteHeader(code)
	json.NewEncoder((*writer)).Encode(ErrorMessage{message, code})
}
