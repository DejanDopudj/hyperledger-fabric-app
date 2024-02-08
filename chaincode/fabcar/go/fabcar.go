/*
SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type CurrencyRates struct {
	Currency1 string  `json:"currency1"`
	Currency2 string  `json:"currency2"`
	Rate      float64 `json:"rate"`
}

type SmartContract struct {
	contractapi.Contract
}

type Bank struct {
	ID           string   `json:"id"`
	Headquarters string   `json:"headquarters"`
	YearFounded  int      `json:"yearFounded"`
	PIB          string   `json:"pib"`
	UserIDs      []string `json:"userIds"`
}

type User struct {
	ID         string   `json:"id"`
	FirstName  string   `json:"firstName"`
	LastName   string   `json:"lastName"`
	Email      string   `json:"email"`
	AccountIDs []string `json:"accountIds"`
}

type Account struct {
	ID       string  `json:"id"`
	Amount   float64 `json:"amount"`
	Currency string  `json:"currency"`
	CardList []string `json:"cardList"`
}

type QueryResultBank struct {
	Key    string `json:"key"`
	Record *Bank
}

type QueryResultUser struct {
	Key    string `json:"key"`
	Record *User
}

type QueryResultAccount struct {
	Key    string `json:"key"`
	Record *Account
}

func toInt(s string) int {
	i, _ := strconv.Atoi(s)
	return i
}

func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	banks := []Bank{
		{ID: "1", Headquarters: "New York", YearFounded: 2000, PIB: "123456789", UserIDs: []string{"1", "2", "3"}},
		{ID: "2", Headquarters: "London", YearFounded: 1995, PIB: "987654321", UserIDs: []string{"4", "5", "6"}},
		{ID: "3", Headquarters: "Tokyo", YearFounded: 2010, PIB: "456789123", UserIDs: []string{"7", "8", "9"}},
		{ID: "4", Headquarters: "Berlin", YearFounded: 2005, PIB: "321654987", UserIDs: []string{"10", "11", "12"}},
	}

	for _, bank := range banks {
		bankAsBytes, _ := json.Marshal(bank)
		err := ctx.GetStub().PutState("BANK"+bank.ID, bankAsBytes)
		if err != nil {
			return fmt.Errorf("failed to put to world state. %s", err.Error())
		}

		for i := 1; i <= 3; i++ {
			user := User{
				ID:         strconv.Itoa((i - 1) + (3 * (toInt(bank.ID) - 1)) + 1),
				FirstName:  "User" + strconv.Itoa(i),
				LastName:   "LastName" + strconv.Itoa(i),
				Email:      "user" + strconv.Itoa(i) + "@example.com",
				AccountIDs: []string{strconv.Itoa(((i - 1) * 2) + 1), strconv.Itoa(((i - 1) * 2) + 2)},
			}
			userAsBytes, _ := json.Marshal(user)
			err := ctx.GetStub().PutState("USER"+user.ID, userAsBytes)
			if err != nil {
				return fmt.Errorf("failed to put to world state. %s", err.Error())
			}
		}

		for i := 1; i <= 12; i++ {
			account := Account{
				ID:       strconv.Itoa(i),
				Amount:   1000.0,
				Currency: "USD",
				CardList: []string{"card1", "card2"},
			}
			accountAsBytes, _ := json.Marshal(account)
			err := ctx.GetStub().PutState("ACCOUNT"+strconv.Itoa(i), accountAsBytes)
			if err != nil {
				return fmt.Errorf("failed to put to world state. %s", err.Error())
			}
		}
	}


	rates := []CurrencyRates{
		{Currency1: "EUR", Currency2: "USD", Rate: 1.22},
		{Currency1: "USD", Currency2: "EUR", Rate: 0.82},
		{Currency1: "USD", Currency2: "RSD", Rate: 99.45},
		{Currency1: "RSD", Currency2: "USD", Rate: 0.0101},
		{Currency1: "EUR", Currency2: "RSD", Rate: 120.96},
		{Currency1: "RSD", Currency2: "EUR", Rate: 0.0083},
	}

	for _, rate := range rates {
		rateAsBytes, err := json.Marshal(rate)
		if err != nil {
			return fmt.Errorf("failed to marshal rate: %v", err)
		}
		if err := ctx.GetStub().PutState("CURR_" + rate.Currency1+rate.Currency2, rateAsBytes); err != nil {
			return fmt.Errorf("failed to put rate into state: %v", err)
		}
	}

	return nil
}

func (s *SmartContract) QueryAllBanks(ctx contractapi.TransactionContextInterface) ([]QueryResultBank, error) {
	startKey := "BANK0"
	endKey := "BANKZ"

	resultsIterator, err := ctx.GetStub().GetStateByRange(startKey, endKey)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var results []QueryResultBank

	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}


		bank := new(Bank)
		_ = json.Unmarshal(queryResponse.Value, bank)

		queryResult := QueryResultBank{Key: queryResponse.Key, Record: bank}
		results = append(results, queryResult)
	}

	return results, nil
}


func (s *SmartContract) QueryAllUsers(ctx contractapi.TransactionContextInterface) ([]QueryResultUser, error) {
	startKey := "USER0"
	endKey := "USER999"

	resultsIterator, err := ctx.GetStub().GetStateByRange(startKey, endKey)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var results []QueryResultUser

	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}


		user := new(User)
		_ = json.Unmarshal(queryResponse.Value, user)

		queryResult := QueryResultUser{Key: queryResponse.Key, Record: user}
		results = append(results, queryResult)
	}

	return results, nil
}


func (s *SmartContract) QueryAllAccounts(ctx contractapi.TransactionContextInterface) ([]QueryResultAccount, error) {
	startKey := "ACCOUNT0"
	endKey := "ACCOUNT999"

	resultsIterator, err := ctx.GetStub().GetStateByRange(startKey, endKey)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var results []QueryResultAccount

	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}


		account := new(Account)
		_ = json.Unmarshal(queryResponse.Value, account)

		queryResult := QueryResultAccount{Key: queryResponse.Key, Record: account}
		results = append(results, queryResult)
	}

	return results, nil
}


func (s *SmartContract) QueryBank(ctx contractapi.TransactionContextInterface, bankId string) (*Bank, error) {

	bankAsBytes, err := ctx.GetStub().GetState("BANK"+bankId);

	if err != nil {
		return nil, fmt.Errorf("Failed to read from world state. %s", err.Error())
	}

	if bankAsBytes == nil {
		return nil, fmt.Errorf("%s does not exist", "BANK"+bankId)
	}

	bank := new(Bank)
	_ = json.Unmarshal(bankAsBytes, bank)

	return bank, nil
}


func (s *SmartContract) QueryUser(ctx contractapi.TransactionContextInterface, userId string) (*User, error) {
	userAsBytes, err := ctx.GetStub().GetState("USER"+userId);

	if err != nil {
		return nil, fmt.Errorf("Failed to read from world state. %s", err.Error())
	}

	if userAsBytes == nil {
		return nil, fmt.Errorf("%s does not exist", "USER"+userId)
	}

	user := new(User)
	_ = json.Unmarshal(userAsBytes, user)

	return user, nil
}


func (s *SmartContract) QueryAccount(ctx contractapi.TransactionContextInterface, accountId string) (*Account, error) {
	accountAsBytes, err := ctx.GetStub().GetState("ACCOUNT"+accountId);

	if err != nil {
		return nil, fmt.Errorf("Failed to read from world state. %s", err.Error())
	}

	if accountAsBytes == nil {
		return nil, fmt.Errorf("%s does not exist", "ACCOUNT"+accountId)
	}

	account := new(Account)
	_ = json.Unmarshal(accountAsBytes, account)

	return account, nil
}

func (s *SmartContract) CreateBank(ctx contractapi.TransactionContextInterface, bankID string, headquarters string, yearFounded int, pib string, userIds []string) error {
	bank := Bank{
		ID:           bankID,
		Headquarters: headquarters,
		YearFounded:  yearFounded,
		PIB:          pib,
		UserIDs:      userIds,
	}

	bankAsBytes, _ := json.Marshal(bank)

	return ctx.GetStub().PutState(bankID, bankAsBytes)
}

func (s *SmartContract) CreateUser(ctx contractapi.TransactionContextInterface, bankID string, userID string, firstName string, lastName string, email string, accountIds []string) error {
	bankAsBytes, err := ctx.GetStub().GetState(bankID)
	if err != nil {
		return fmt.Errorf("failed to read bank with id %s: %v", bankID, err)
	}
	if bankAsBytes == nil {
		return fmt.Errorf("bank with id %s does not exist", bankID)
	}

	var bank Bank
	if err := json.Unmarshal(bankAsBytes, &bank); err != nil {
		return err
	}

	bank.UserIDs = append(bank.UserIDs, userID)

	bankAsBytes, err = json.Marshal(bank)
	if err != nil {
		return err
	}
	if err := ctx.GetStub().PutState(bankID, bankAsBytes); err != nil {
		return fmt.Errorf("failed to update bank with id %s: %v", bankID, err)
	}

	user := User{
		ID:         userID,
		FirstName:  firstName,
		LastName:   lastName,
		Email:      email,
		AccountIDs: accountIds,
	}

	userAsBytes, err := json.Marshal(user)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(userID, userAsBytes)
}


func (s *SmartContract) CreateAccount(ctx contractapi.TransactionContextInterface, userID string, accountID string, amount float64, currency string, cardList []string) error {
	userAsBytes, err := ctx.GetStub().GetState(userID)
	if err != nil {
		return fmt.Errorf("failed to read user with id %s: %v", userID, err)
	}
	if userAsBytes == nil {
		return fmt.Errorf("user with id %s does not exist", userID)
	}

	var user User
	if err := json.Unmarshal(userAsBytes, &user); err != nil {
		return err
	}

	user.AccountIDs = append(user.AccountIDs, accountID)

	userAsBytes, err = json.Marshal(user)
	if err != nil {
		return err
	}
	if err := ctx.GetStub().PutState(userID, userAsBytes); err != nil {
		return fmt.Errorf("failed to update user with id %s: %v", userID, err)
	}

	account := Account{
		ID:       accountID,
		Amount:   amount,
		Currency: currency,
		CardList: cardList,
	}

	accountAsBytes, err := json.Marshal(account)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(accountID, accountAsBytes)
}


func (s *SmartContract) MakeWithdrawal(ctx contractapi.TransactionContextInterface, accountID string, amount float64) error {
	accountAsBytes, err := ctx.GetStub().GetState(accountID)
	if err != nil {
		return fmt.Errorf("failed to read account with id %s: %v", accountID, err)
	}
	if accountAsBytes == nil {
		return fmt.Errorf("account with id %s does not exist", accountID)
	}

	var account Account
	if err := json.Unmarshal(accountAsBytes, &account); err != nil {
		return err
	}

	if account.Amount < amount {
		return fmt.Errorf("insufficient balance in account %s", accountID)
	}

	account.Amount -= amount

	accountAsBytes, err = json.Marshal(account)
	if err != nil {
		return err
	}
	if err := ctx.GetStub().PutState(accountID, accountAsBytes); err != nil {
		return fmt.Errorf("failed to update account with id %s: %v", accountID, err)
	}

	return nil
}

func (s *SmartContract) MakePayment(ctx contractapi.TransactionContextInterface, accountID string, amount float64, currency string) error {
	accountAsBytes, err := ctx.GetStub().GetState(accountID)
	if err != nil {
		return fmt.Errorf("failed to read account with id %s: %v", accountID, err)
	}
	if accountAsBytes == nil {
		return fmt.Errorf("account with id %s does not exist", accountID)
	}

	var account Account
	if err := json.Unmarshal(accountAsBytes, &account); err != nil {
		return err
	}

	if account.Currency != currency {
		return fmt.Errorf("payment currency %s does not match account currency %s", currency, account.Currency)
	}

	account.Amount += amount

	accountAsBytes, err = json.Marshal(account)
	if err != nil {
		return err
	}
	if err := ctx.GetStub().PutState(accountID, accountAsBytes); err != nil {
		return fmt.Errorf("failed to update account with id %s: %v", accountID, err)
	}

	return nil
}

func (s *SmartContract) CheckCurrencyMatch(ctx contractapi.TransactionContextInterface, account1ID string, account2ID string) (bool, error) {
	account1AsBytes, err := ctx.GetStub().GetState(account1ID)
	if err != nil {
		return false, fmt.Errorf("failed to read account with id %s: %v", account1ID, err)
	}
	if account1AsBytes == nil {
		return false, fmt.Errorf("account with id %s does not exist", account1ID)
	}

	account2AsBytes, err := ctx.GetStub().GetState(account2ID)
	if err != nil {
		return false, fmt.Errorf("failed to read account with id %s: %v", account2ID, err)
	}
	if account2AsBytes == nil {
		return false, fmt.Errorf("account with id %s does not exist", account2ID)
	}

	var account1, account2 Account
	if err := json.Unmarshal(account1AsBytes, &account1); err != nil {
		return false, err
	}
	if err := json.Unmarshal(account2AsBytes, &account2); err != nil {
		return false, err
	}

	return account1.Currency == account2.Currency, nil
}

func (s *SmartContract) TransferBetweenAccounts(ctx contractapi.TransactionContextInterface, fromAccountID string, toAccountID string, amount float64, currency string) error {

	currenciesMatch, err := s.CheckCurrencyMatch(ctx, fromAccountID, toAccountID)
	if err != nil {
		return fmt.Errorf("failed to check currency match: %v", err)
	}


	fromAccountAsBytes, err := ctx.GetStub().GetState(fromAccountID)
	if err != nil {
		return fmt.Errorf("failed to read sender's account with id %s: %v", fromAccountID, err)
	}
	if fromAccountAsBytes == nil {
		return fmt.Errorf("sender's account with id %s does not exist", fromAccountID)
	}


	toAccountAsBytes, err := ctx.GetStub().GetState(toAccountID)
	if err != nil {
		return fmt.Errorf("failed to read receiver's account with id %s: %v", toAccountID, err)
	}
	if toAccountAsBytes == nil {
		return fmt.Errorf("receiver's account with id %s does not exist", toAccountID)
	}


	var fromAccount, toAccount Account
	if err := json.Unmarshal(fromAccountAsBytes, &fromAccount); err != nil {
		return err
	}
	if err := json.Unmarshal(toAccountAsBytes, &toAccount); err != nil {
		return err
	}

	if !currenciesMatch {
		conversionRateAsBytes, err := ctx.GetStub().GetState(currency + fromAccount.Currency)
		if err != nil {
			return fmt.Errorf("failed to read conversion rate for currencies %s to %s: %v", currency, fromAccount.Currency, err)
		}
		if conversionRateAsBytes == nil {
			return fmt.Errorf("conversion rate for currencies %s to %s not found", currency, fromAccount.Currency)
		}
		var conversionRate CurrencyRates
		if err := json.Unmarshal(conversionRateAsBytes, &conversionRate); err != nil {
			return fmt.Errorf("failed to unmarshal conversion rate for currencies %s to %s: %v", currency, fromAccount.Currency, err)
		}

		amount *= conversionRate.Rate
	}

	if fromAccount.Amount < amount {
		return fmt.Errorf("insufficient balance in sender's account %s", fromAccountID)
	}

	fromAccount.Amount -= amount
	toAccount.Amount += amount

	fromAccountAsBytes, err = json.Marshal(fromAccount)
	if err != nil {
		return err
	}
	if err := ctx.GetStub().PutState(fromAccountID, fromAccountAsBytes); err != nil {
		return fmt.Errorf("failed to update sender's account with id %s: %v", fromAccountID, err)
	}

	toAccountAsBytes, err = json.Marshal(toAccount)
	if err != nil {
		return err
	}
	if err := ctx.GetStub().PutState(toAccountID, toAccountAsBytes); err != nil {
		return fmt.Errorf("failed to update receiver's account with id %s: %v", toAccountID, err)
	}

	return nil
}



func main() {

	chaincode, err := contractapi.NewChaincode(new(SmartContract))

	if err != nil {
		fmt.Printf("Error create fabcar chaincode: %s", err.Error())
		return
	}

	if err := chaincode.Start(); err != nil {
		fmt.Printf("Error starting fabcar chaincode: %s", err.Error())
	}
}
