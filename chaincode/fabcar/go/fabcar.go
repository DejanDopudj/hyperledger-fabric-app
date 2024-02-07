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

	return nil
}

// QueryAllBanks returns all banks found in world state
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
	// Retrieve the bank record
	bankAsBytes, err := ctx.GetStub().GetState(bankID)
	if err != nil {
		return fmt.Errorf("failed to read bank with id %s: %v", bankID, err)
	}
	if bankAsBytes == nil {
		return fmt.Errorf("bank with id %s does not exist", bankID)
	}

	// Unmarshal bank data
	var bank Bank
	if err := json.Unmarshal(bankAsBytes, &bank); err != nil {
		return err
	}

	// Add user to the bank
	bank.UserIDs = append(bank.UserIDs, userID)

	// Update the bank record
	bankAsBytes, err = json.Marshal(bank)
	if err != nil {
		return err
	}
	if err := ctx.GetStub().PutState(bankID, bankAsBytes); err != nil {
		return fmt.Errorf("failed to update bank with id %s: %v", bankID, err)
	}

	// Create user record
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
	// Retrieve the user record
	userAsBytes, err := ctx.GetStub().GetState(userID)
	if err != nil {
		return fmt.Errorf("failed to read user with id %s: %v", userID, err)
	}
	if userAsBytes == nil {
		return fmt.Errorf("user with id %s does not exist", userID)
	}

	// Unmarshal user data
	var user User
	if err := json.Unmarshal(userAsBytes, &user); err != nil {
		return err
	}

	// Add account to the user
	user.AccountIDs = append(user.AccountIDs, accountID)

	// Update the user record
	userAsBytes, err = json.Marshal(user)
	if err != nil {
		return err
	}
	if err := ctx.GetStub().PutState(userID, userAsBytes); err != nil {
		return fmt.Errorf("failed to update user with id %s: %v", userID, err)
	}

	// Create account record
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
