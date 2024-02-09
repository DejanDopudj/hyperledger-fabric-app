package src

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/hyperledger/fabric-sdk-go/pkg/gateway"
)

type ErrorMessage struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

type Handler struct {
	Validator *validator.Validate
	Wallet    *gateway.Wallet
}

func (handler *Handler) CreateUser(writer http.ResponseWriter, req *http.Request) {
	org := req.URL.Query().Get("org")

	if org == "" {
		responseError(writer, "Org enroll is required", http.StatusBadRequest)
		return
	}

	var userDto RegisterUser
	err := json.NewDecoder(req.Body).Decode(&userDto)

	if err != nil || handler.Validator.Struct(userDto) != nil {
		responseError(writer, "Invalid Body", http.StatusBadRequest)
		return
	}

	contract, err := getContract(handler.Wallet, org)

	response, err := contract.SubmitTransaction(
		"CreateUser",
		userDto.BankID,
		userDto.UserID,
		userDto.FirstName,
		userDto.LastName,
		userDto.Email,
	)

	responseOK(writer, response)
}

func (handler *Handler) GetUser(writer http.ResponseWriter, req *http.Request) {
	user := mux.Vars(req)["user"]

	if user == "" {
		responseError(writer, "User id is required", http.StatusBadRequest)
		return
	}

	org := req.URL.Query().Get("org")

	if org == "" {
		responseError(writer, "Org enroll is required", http.StatusBadRequest)
		return
	}

	contract, err := getContract(handler.Wallet, org)

	if err != nil {
		responseError(writer, err.Error(), http.StatusBadRequest)
		return
	}

	response, err := contract.EvaluateTransaction("QueryUser", user)

	if err != nil {
		responseError(writer, err.Error(), http.StatusBadRequest)
		return
	}

	responseOK(writer, response)
}

func (handler *Handler) GetBank(writer http.ResponseWriter, req *http.Request) {
	bank := mux.Vars(req)["bank"]

	if bank == "" {
		responseError(writer, "Bank id is required", http.StatusBadRequest)
		return
	}

	org := req.URL.Query().Get("org")

	if org == "" {
		responseError(writer, "Org enroll is required", http.StatusBadRequest)
		return
	}

	contract, err := getContract(handler.Wallet, org)

	if err != nil {
		responseError(writer, err.Error(), http.StatusBadRequest)
		return
	}

	response, err := contract.EvaluateTransaction("QueryBank", bank)

	if err != nil {
		responseError(writer, err.Error(), http.StatusBadRequest)
		return
	}

	responseOK(writer, response)
}

func (handler *Handler) GetBankAccount(writer http.ResponseWriter, req *http.Request) {
	org := req.URL.Query().Get("org")

	if org == "" {
		responseError(writer, "Org enroll is required", http.StatusBadRequest)
		return
	}

	account := mux.Vars(req)["account"]

	if account == "" {
		responseError(writer, "Account id is required", http.StatusBadRequest)
		return
	}

	contract, err := getContract(handler.Wallet, org)

	if err != nil {
		responseError(writer, err.Error(), http.StatusBadRequest)
		return
	}

	response, err := contract.EvaluateTransaction("QueryAccount", account)

	if err != nil {
		responseError(writer, err.Error(), http.StatusBadRequest)
		return
	}

	responseOK(writer, response)
}

func (handler *Handler) CreateBankAccount(writer http.ResponseWriter, req *http.Request) {
	org := req.URL.Query().Get("org")

	if org == "" {
		responseError(writer, "Org enroll is required", http.StatusBadRequest)
		return
	}

	var accountDto CreateBankAccount
	err := json.NewDecoder(req.Body).Decode(&accountDto)

	if err != nil || handler.Validator.Struct(accountDto) != nil {
		responseError(writer, "Invalid Body", http.StatusBadRequest)
		return
	}

	contract, err := getContract(handler.Wallet, org)

	if err != nil {
		responseError(writer, err.Error(), http.StatusInternalServerError)
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
		responseError(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	responseOK(writer, response)
}

func (handler *Handler) MakePayment(writer http.ResponseWriter, req *http.Request) {
	org := req.URL.Query().Get("org")

	if org == "" {
		responseError(writer, "Org enroll is required", http.StatusBadRequest)
		return
	}

	var fundsDto Payment
	err := json.NewDecoder(req.Body).Decode(&fundsDto)

	if err != nil || handler.Validator.Struct(fundsDto) != nil {
		responseError(writer, "Invalid Body", http.StatusBadRequest)
		return
	}

	contract, err := getContract(handler.Wallet, org)

	if err != nil {
		responseError(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	result, err := contract.SubmitTransaction(
		"MakePayment",
		fundsDto.AccountID,
		strconv.FormatFloat(fundsDto.Amount, 'f', -1, 64),
		fundsDto.Currency,
	)

	if err != nil {
		responseError(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	responseOK(writer, result)
}

func (handler *Handler) MakeWithdrawal(writer http.ResponseWriter, req *http.Request) {
	org := req.URL.Query().Get("org")

	if org == "" {
		responseError(writer, "Org enroll is required", http.StatusBadRequest)
		return
	}

	var fundsDto Withdrawal
	err := json.NewDecoder(req.Body).Decode(&fundsDto)

	if err != nil || handler.Validator.Struct(fundsDto) != nil {
		responseError(writer, "Invalid Body", http.StatusBadRequest)
		return
	}

	contract, err := getContract(handler.Wallet, org)

	if err != nil {
		responseError(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	response, err := contract.SubmitTransaction(
		"MakeWithdrawal",
		fundsDto.AccountID,
		strconv.FormatFloat(fundsDto.Amount, 'f', -1, 64),
	)

	if err != nil {
		responseError(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	responseOK(writer, response)
}

func (handler *Handler) TransferAssets(writer http.ResponseWriter, req *http.Request) {
	org := req.URL.Query().Get("org")

	if org == "" {
		responseError(writer, "Org enroll is required", http.StatusBadRequest)
		return
	}

	var transferDto TransferFunds
	err := json.NewDecoder(req.Body).Decode(&transferDto)

	if err != nil { //|| handler.Validator.Struct(transferDto) != nil {
		responseError(writer, "Invalid Body", http.StatusBadRequest)
		return
	}

	contract, err := getContract(handler.Wallet, org)

	if err != nil {
		responseError(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	conversion, err := contract.EvaluateTransaction(
		"CheckCurrencyMatch",
		transferDto.FromAccountID,
		transferDto.ToAccountID,
	)

	if err != nil {
		responseError(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	currenciesMatch, err := strconv.ParseBool(string(conversion))

	if err != nil {
		responseError(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	if !currenciesMatch && !transferDto.AcceptConversion {
		responseError(writer, "Failed to make transaction since the account currencies dont match.", http.StatusInternalServerError)
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

func (handler *Handler) QueryBankFull(writer http.ResponseWriter, req *http.Request) {
	org := req.URL.Query().Get("org")

	if org == "" {
		responseError(writer, "Org enroll is required", http.StatusBadRequest)
		return
	}

	id := req.URL.Query().Get("id")
	headquarters := req.URL.Query().Get("headquarters")
	pib := req.URL.Query().Get("pib")
	yearFounded := req.URL.Query().Get("yearFounded")
	userId := req.URL.Query().Get("userId")

	contract, err := getContract(handler.Wallet, org)

	if err != nil {
		responseError(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	response, err := contract.EvaluateTransaction(
		"QueryBankFull",
		id,
		headquarters,
		pib,
		yearFounded,
		userId,
	)

	log.Println(response)

	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	responseOK(writer, response)
}

func (handler *Handler) QueryAccountFull(writer http.ResponseWriter, req *http.Request) {
	org := req.URL.Query().Get("org")

	if org == "" {
		responseError(writer, "Org enroll is required", http.StatusBadRequest)
		return
	}

	id := req.URL.Query().Get("id")
	amount := req.URL.Query().Get("amount")
	currency := req.URL.Query().Get("currency")
	userId := req.URL.Query().Get("userId")

	contract, err := getContract(handler.Wallet, org)

	if err != nil {
		responseError(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	response, err := contract.EvaluateTransaction(
		"QueryAccountFull",
		id,
		amount,
		currency,
		userId,
	)

	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	responseOK(writer, response)
}

func (handler *Handler) QueryUserFull(writer http.ResponseWriter, req *http.Request) {
	org := req.URL.Query().Get("org")

	if org == "" {
		responseError(writer, "Org enroll is required", http.StatusBadRequest)
		return
	}

	id := req.URL.Query().Get("id")
	firstName := req.URL.Query().Get("firstName")
	lastName := req.URL.Query().Get("lastName")
	email := req.URL.Query().Get("email")
	accountId := req.URL.Query().Get("accountId")
	bankId := req.URL.Query().Get("bankId")

	contract, err := getContract(handler.Wallet, org)

	if err != nil {
		responseError(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	response, err := contract.EvaluateTransaction(
		"QueryUserFull",
		id,
		firstName,
		lastName,
		email,
		accountId,
		bankId,
	)

	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	responseOK(writer, response)
}

func (handler *Handler) QueryAllBanks(writer http.ResponseWriter, req *http.Request) {
	handler.queryAll(writer, "QueryAllBanks")
}
func (handler *Handler) QueryAllUsers(writer http.ResponseWriter, req *http.Request) {
	handler.queryAll(writer, "QueryAllUsers")
}
func (handler *Handler) QueryAllAccounts(writer http.ResponseWriter, req *http.Request) {
	handler.queryAll(writer, "QueryAllAccounts")
}

func (handler *Handler) queryAll(writer http.ResponseWriter, query string) {
	contract, err := getContract(handler.Wallet, "1")

	if err != nil {
		responseError(writer, err.Error(), http.StatusBadRequest)
		return
	}

	response, err := contract.EvaluateTransaction(query)

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
	//writer.WriteHeader(http.StatusOK)
	writer.Header().Set("Content-Type", "application/json")

	writer.Write(result)
}

func responseError(writer http.ResponseWriter, message string, code int) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(code)

	err := ErrorMessage{message, code}
	json.NewEncoder(writer).Encode(err)
}
