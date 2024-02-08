package src

import (
	"encoding/json"
	"net/http"
	"strconv"

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
		responseError(&writer, "Invalid Body", http.StatusBadRequest)
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
		responseError(&writer, "User id is required", http.StatusBadRequest)
		return
	}

	bank := req.URL.Query().Get("bank")

	if bank == "" {
		responseError(&writer, "Bank id is required", http.StatusBadRequest)
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

	responseOK(writer, response)
}

func (handler *Handler) GetBank(writer http.ResponseWriter, req *http.Request) {
	bank := mux.Vars(req)["bank"]

	if bank == "" {
		responseError(&writer, "Bank id is required", http.StatusBadRequest)
		return
	}

	user := req.URL.Query().Get("user")

	if user == "" {
		responseError(&writer, "User id is required", http.StatusBadRequest)
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

	responseOK(writer, response)
}

func (handler *Handler) GetBankAccount(writer http.ResponseWriter, req *http.Request) {
	bank := req.URL.Query().Get("bank")

	if bank == "" {
		responseError(&writer, "Bank id is required", http.StatusBadRequest)
		return
	}

	user := req.URL.Query().Get("user")

	if user == "" {
		responseError(&writer, "User id is required", http.StatusBadRequest)
		return
	}

	account := mux.Vars(req)["account"]

	if user == "" {
		responseError(&writer, "Account id is required", http.StatusBadRequest)
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

	responseOK(writer, response)
}

func (handler *Handler) CreateBankAccount(writer http.ResponseWriter, req *http.Request) {
	var accountDto CreateBankAccount
	err := json.NewDecoder(req.Body).Decode(&accountDto)

	if err != nil || handler.Validator.Struct(accountDto) != nil {
		responseError(&writer, "Invalid Body", http.StatusBadRequest)
		return
	}

	contract, err := getContract(handler.Wallet, accountDto.BankID, accountDto.UserID)

	if err != nil {
		responseError(&writer, err.Error(), http.StatusInternalServerError)
		return
	}

	cardList, _ := json.Marshal(accountDto.CardList)

	response, err := contract.SubmitTransaction(
		"CreateAccount",
		accountDto.UserID,
		accountDto.AccountID,
		strconv.FormatFloat(accountDto.Amount, 'f', -1, 64),
		accountDto.Currency,
		string(cardList),
	)

	if err != nil {
		responseError(&writer, err.Error(), http.StatusInternalServerError)
		return
	}

	responseOK(writer, response)
}

func (handler *Handler) MakePayment(writer http.ResponseWriter, req *http.Request) {
	var fundsDto Payment
	err := json.NewDecoder(req.Body).Decode(&fundsDto)

	if err != nil || handler.Validator.Struct(fundsDto) != nil {
		responseError(&writer, "Invalid Body", http.StatusBadRequest)
		return
	}

	contract, err := getContract(handler.Wallet, fundsDto.BankID, fundsDto.UserID)

	if err != nil {
		responseError(&writer, err.Error(), http.StatusInternalServerError)
		return
	}

	result, err := contract.SubmitTransaction(
		"MakePayment",
		fundsDto.AccountID,
		strconv.FormatFloat(fundsDto.Amount, 'f', -1, 64),
		fundsDto.Currency,
	)

	if err != nil {
		responseError(&writer, err.Error(), http.StatusInternalServerError)
		return
	}

	responseOK(writer, result)
}

func (handler *Handler) MakeWithdrawal(writer http.ResponseWriter, req *http.Request) {
	var fundsDto Withdrawal
	err := json.NewDecoder(req.Body).Decode(&fundsDto)

	if err != nil || handler.Validator.Struct(fundsDto) != nil {
		responseError(&writer, "Invalid Body", http.StatusBadRequest)
		return
	}

	contract, err := getContract(handler.Wallet, fundsDto.BankID, fundsDto.UserID)

	if err != nil {
		responseError(&writer, err.Error(), http.StatusInternalServerError)
		return
	}

	response, err := contract.SubmitTransaction(
		"MakeWithdrawal",
		fundsDto.AccountID,
		strconv.FormatFloat(fundsDto.Amount, 'f', -1, 64),
	)

	if err != nil {
		responseError(&writer, err.Error(), http.StatusInternalServerError)
		return
	}

	responseOK(writer, response)
}

func (handler *Handler) TransferAssets(writer http.ResponseWriter, req *http.Request) {
	var transferDto TransferFunds
	err := json.NewDecoder(req.Body).Decode(&transferDto)

	if err != nil || handler.Validator.Struct(transferDto) != nil {
		responseError(&writer, "Invalid Body", http.StatusBadRequest)
		return
	}

	contract, err := getContract(handler.Wallet, transferDto.BankID, transferDto.UserID)

	if err != nil {
		responseError(&writer, err.Error(), http.StatusInternalServerError)
		return
	}

	conversion, err := contract.EvaluateTransaction(
		"CheckCurrencyMatch",
		transferDto.FromAccountID,
		transferDto.ToAccountID,
	)

	if err != nil {
		responseError(&writer, err.Error(), http.StatusInternalServerError)
		return
	}

	currenciesMatch, err := strconv.ParseBool(string(conversion))

	if err != nil {
		responseError(&writer, err.Error(), http.StatusInternalServerError)
		return
	}

	if !currenciesMatch && !transferDto.AcceptConversion {
		responseError(&writer, "Failed to make transaction since the account currencies dont match.", http.StatusInternalServerError)
		return
	}

	response, err := contract.SubmitTransaction(
		"TransferBetweenAccounts",
		transferDto.FromAccountID,
		transferDto.ToAccountID,
		strconv.FormatFloat(transferDto.Amount, 'f', -1, 64),
	)

	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	responseOK(writer, response)
}

func (handler *Handler) QueryAllBanks(writer http.ResponseWriter, req *http.Request) {
	contract, err := getContract(handler.Wallet, "1", "app")

	if err != nil {
		responseError(&writer, err.Error(), http.StatusBadRequest)
		return
	}

	response, err := contract.EvaluateTransaction(
		"QueryAllBanks",
	)

	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	responseOK(writer, response)
}

func mockResponse(call string) ([]byte, error) {
	return json.Marshal(ErrorMessage{})
}

func responseOK(writer http.ResponseWriter, result []byte) {
	writer.WriteHeader(http.StatusOK)
	writer.Header().Set("Content-Type", "application/json")

	json.NewEncoder(writer).Encode(string(result))
}

func responseError(writer *http.ResponseWriter, message string, code int) {
	(*writer).WriteHeader(code)
	json.NewEncoder((*writer)).Encode(ErrorMessage{message, code})
}
