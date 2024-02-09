/*
SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"encoding/json"
	"fmt"
    "regexp"
	"errors"
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
	BankID	   string   `json:"bankId"`
}

type Account struct {
	ID       string  `json:"id"`
	Amount   float64 `json:"amount"`
	Currency string  `json:"currency"`
	CardList []string `json:"cardList"`
	UserID	 string  `json:"userId"`
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
		{ID: "1", Headquarters: "New York", YearFounded: 2000, PIB: "123456789", UserIDs: []string{"1_1", "1_2", "1_3"}},
		{ID: "2", Headquarters: "London", YearFounded: 1995, PIB: "987654321", UserIDs: []string{"2_1", "2_2", "2_3"}},
		{ID: "3", Headquarters: "Tokyo", YearFounded: 2010, PIB: "456789123", UserIDs: []string{"3_1", "3_2", "3_3"}},
		{ID: "4", Headquarters: "Berlin", YearFounded: 2005, PIB: "321654987", UserIDs: []string{"4_1", "4_2", "4_3"}},
	}

	for _, bank := range banks {
		bankAsBytes, _ := json.Marshal(bank)
		err := ctx.GetStub().PutState("BANK"+bank.ID, bankAsBytes)
		if err != nil {
			return fmt.Errorf("failed to put to world state. %s", err.Error())
		}
		account_id := 1
		for i := 1; i <= 3; i++ {
			user := User{
				ID:         bank.ID + "_" + strconv.Itoa(i),
				FirstName:  "User" + strconv.Itoa(i),
				LastName:   "LastName" + strconv.Itoa(i),
				Email:      "user" + strconv.Itoa(i) + "@example.com",
				AccountIDs: []string{bank.ID+"_"+strconv.Itoa(account_id), bank.ID+"_"+strconv.Itoa(account_id+1)},
				BankID:     bank.ID,
			}
			userAsBytes, _ := json.Marshal(user)
			err := ctx.GetStub().PutState("USER"+user.ID, userAsBytes)
			if err != nil {
				return fmt.Errorf("failed to put to world state. %s", err.Error())
			}

			for j := 1; j <= 2; j++ {
				account := Account{
					ID:       bank.ID+"_"+strconv.Itoa(account_id),
					Amount:   1000.0,
					Currency: "USD",
					CardList: []string{"card1", "card2"},
					UserID:   user.ID,
				}
				accountAsBytes, _ := json.Marshal(account)
				err := ctx.GetStub().PutState("ACCOUNT"+account.ID, accountAsBytes)
				if err != nil {
					return fmt.Errorf("failed to put to world state. %s", err.Error())
				}
				account_id++
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


func (s *SmartContract) QueryAllBanks22(ctx contractapi.TransactionContextInterface) (string, error) {

    certificate, err := ctx.GetClientIdentity().GetMSPID()
    if err != nil {
        return "", err
    }

    return certificate, nil
}

func (s *SmartContract) QueryAllBanks23(ctx contractapi.TransactionContextInterface) (string, error) {

    certificate, err := ctx.GetClientIdentity().GetMSPID()
    if err != nil {
        return "", err
    }
	permission,_ := ExtractOrgNumber(certificate)

    return permission, nil
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
	endKey := "USERZ"

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
	endKey := "ACCOUNTZ"

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

func (s *SmartContract) compareIds(ctx contractapi.TransactionContextInterface, bankId string) bool {
	certificate, err := ctx.GetClientIdentity().GetMSPID()
	
    if err != nil {
        return false
    }
	permission,err := ExtractOrgNumber(certificate)
    if err != nil {
        return false
    }
	
	if permission != bankId {
		return false
	}
	return true
}

func (s *SmartContract) QueryBank(ctx contractapi.TransactionContextInterface, bankId string) (*Bank, error) {
	
	if(!s.compareIds(ctx,bankId)){
		return nil, fmt.Errorf("User does not have rights to acces bank %s", bankId)
	}

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

func (s *SmartContract) GetBankIdFromAccount(ctx contractapi.TransactionContextInterface, userId string) (string, error) {
    user, err := s.QueryUser(ctx, userId)
    if err != nil {
        return "", err
    }
    
    return user.BankID, nil
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


	if(!s.compareIds(ctx,user.BankID)){
		return nil, fmt.Errorf("User does not have rights to acces bank %s", user.BankID)
	}

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

	bankId, err := s.GetBankIdFromAccount(ctx, account.UserID)
	if err != nil {
		return nil, fmt.Errorf("Error %s",err.Error())
	}

	if(!s.compareIds(ctx,bankId)){
		return nil, fmt.Errorf("User does not have rights to acces bank %s", bankId)
	}

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

	return ctx.GetStub().PutState("BANK"+bankID, bankAsBytes)
}

func (s *SmartContract) CreateUser(ctx contractapi.TransactionContextInterface, bankID string, userID string, firstName string, lastName string, email string) error {
	bytes, err := ctx.GetStub().GetState("USER"+userID);

    if bytes != nil {
        return fmt.Errorf("User %s already exists", "USER"+userID)
    }

	if(!s.compareIds(ctx,bankID)){
		return fmt.Errorf("User does not have rights to acces bank %s", bankID)
	}

	bankAsBytes, err := ctx.GetStub().GetState("BANK"+bankID)
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
	if err := ctx.GetStub().PutState("BANK"+bankID, bankAsBytes); err != nil {
		return fmt.Errorf("failed to update bank with id %s: %v", bankID, err)
	}

	user := User{
		ID:         userID,
		FirstName:  firstName,
		LastName:   lastName,
		Email:      email,
		AccountIDs: []string{},
		BankID:     bankID,
	}

	userAsBytes, err := json.Marshal(user)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState("USER"+userID, userAsBytes)
}


func (s *SmartContract) CreateAccount(ctx contractapi.TransactionContextInterface, userID string, accountID string, amount float64, currency string, cardList string) error {
    bytes, err := ctx.GetStub().GetState("ACCOUNT"+accountID);

    if bytes != nil {
        return fmt.Errorf("Account %s already exists", "ACCOUNT"+accountID)
    }

	bankId, err := s.GetBankIdFromAccount(ctx, userID)
	if err != nil {
		return fmt.Errorf("Error %s",err.Error())
	}

	if(!s.compareIds(ctx,bankId)){
		return fmt.Errorf("User does not have rights to acces bank %s", bankId)
	}

	var cardIds []string
	err = json.Unmarshal([]byte(cardList), &cardIds)
	if err != nil {
		return fmt.Errorf("failed to unmarshall card ids %s: %v", userID, err)
	}
	userAsBytes, err := ctx.GetStub().GetState("USER"+userID)

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
	if err := ctx.GetStub().PutState("USER"+userID, userAsBytes); err != nil {
		return fmt.Errorf("failed to update user with id %s: %v", userID, err)
	}

	account := Account{
		ID:       accountID,
		Amount:   amount,
		Currency: currency,
		CardList: cardIds,
		UserID:   userID,
	}

	accountAsBytes, err := json.Marshal(account)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState("ACCOUNT"+accountID, accountAsBytes)
}


func (s *SmartContract) MakeWithdrawal(ctx contractapi.TransactionContextInterface, accountID string, amount float64) error {
	accountAsBytes, err := ctx.GetStub().GetState("ACCOUNT"+accountID)
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

	bankId, err := s.GetBankIdFromAccount(ctx, account.UserID)
	if err != nil {
		return fmt.Errorf("Error %s",err.Error())
	}

	if(!s.compareIds(ctx,bankId)){
		return fmt.Errorf("User does not have rights to acces bank %s", bankId)
	}

	if account.Amount < amount {
		return fmt.Errorf("insufficient balance in account %s", accountID)
	}

	account.Amount -= amount

	accountAsBytes, err = json.Marshal(account)
	if err != nil {
		return err
	}
	if err := ctx.GetStub().PutState("ACCOUNT"+accountID, accountAsBytes); err != nil {
		return fmt.Errorf("failed to update account with id %s: %v", accountID, err)
	}

	return nil
}

func (s *SmartContract) MakePayment(ctx contractapi.TransactionContextInterface, accountID string, amount float64, currency string) error {
	accountAsBytes, err := ctx.GetStub().GetState("ACCOUNT"+accountID)
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

	bankId, err := s.GetBankIdFromAccount(ctx, account.UserID)
	if err != nil {
		return fmt.Errorf("Error %s",err.Error())
	}

	if(!s.compareIds(ctx,bankId)){
		return fmt.Errorf("User does not have rights to acces bank %s", bankId)
	}


	if account.Currency != currency {
		return fmt.Errorf("payment currency %s does not match account currency %s", currency, account.Currency)
	}

	account.Amount += amount

	accountAsBytes, err = json.Marshal(account)
	if err != nil {
		return err
	}
	if err := ctx.GetStub().PutState("ACCOUNT"+accountID, accountAsBytes); err != nil {
		return fmt.Errorf("failed to update account with id %s: %v", accountID, err)
	}

	return nil
}

func (s *SmartContract) CheckCurrencyMatch(ctx contractapi.TransactionContextInterface, account1ID string, account2ID string) (bool, error) {
	account1AsBytes, err := ctx.GetStub().GetState("ACCOUNT"+account1ID)
	if err != nil {
		return false, fmt.Errorf("failed to read account with id %s: %v", account1ID, err)
	}
	if account1AsBytes == nil {
		return false, fmt.Errorf("account with id %s does not exist", account1ID)
	}

	account2AsBytes, err := ctx.GetStub().GetState("ACCOUNT"+account2ID)
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

func (s *SmartContract) TransferBetweenAccounts(ctx contractapi.TransactionContextInterface, fromAccountID string, toAccountID string, amount float64) error {

	currenciesMatch, err := s.CheckCurrencyMatch(ctx, fromAccountID, toAccountID)
	if err != nil {
		return fmt.Errorf("failed to check currency match: %v", err)
	}


	fromAccountAsBytes, err := ctx.GetStub().GetState("ACCOUNT"+fromAccountID)
	if err != nil {
		return fmt.Errorf("failed to read sender's account with id %s: %v", fromAccountID, err)
	}
	if fromAccountAsBytes == nil {
		return fmt.Errorf("sender's account with id %s does not exist", fromAccountID)
	}


	toAccountAsBytes, err := ctx.GetStub().GetState("ACCOUNT"+toAccountID)
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

	bankId, err := s.GetBankIdFromAccount(ctx, fromAccount.UserID)
	if err != nil {
		return fmt.Errorf("Error %s",err.Error())
	}

	if(!s.compareIds(ctx,bankId)){
		return fmt.Errorf("User does not have rights to acces bank %s", bankId)
	}
	amountSent := amount
	currency := toAccount.Currency
	if !currenciesMatch {
		conversionRateAsBytes, err := ctx.GetStub().GetState("CURR_"+fromAccount.Currency+currency)
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

		amountSent *= conversionRate.Rate
	}

	if fromAccount.Amount < amount {
		return fmt.Errorf("insufficient balance in sender's account %s", fromAccountID)
	}

	fromAccount.Amount -= amount
	toAccount.Amount += amountSent

	fromAccountAsBytes, err = json.Marshal(fromAccount)
	if err != nil {
		return err
	}
	if err := ctx.GetStub().PutState("ACCOUNT"+fromAccountID, fromAccountAsBytes); err != nil {
		return fmt.Errorf("failed to update sender's account with id %s: %v", fromAccountID, err)
	}

	toAccountAsBytes, err = json.Marshal(toAccount)
	if err != nil {
		return err
	}
	if err := ctx.GetStub().PutState("ACCOUNT"+toAccountID, toAccountAsBytes); err != nil {
		return fmt.Errorf("failed to update receiver's account with id %s: %v", toAccountID, err)
	}

	return nil
}


func (s *SmartContract) QueryUserFull(ctx contractapi.TransactionContextInterface, id, firstName, lastName, email, accountId, bankId string) ([]*User, error) {
    selector := map[string]interface{}{}
    if id != "" {
        selector["id"] = id
    }
    if firstName != "" {
        selector["firstName"] = firstName
    }
    if lastName != "" {
        selector["lastName"] = lastName
    }
    if email != "" {
        selector["email"] = email
    }
    if bankId != "" {
        selector["bankId"] = bankId
    }
    if accountId != "" {
		selector["accountIds"] = map[string]interface{}{
			"$elemMatch": map[string]interface{}{
				"$eq": accountId,
			},
		}
	}
    query := map[string]interface{}{
        "selector": selector,
    }

    selectorBytes, err := json.Marshal(query)
	fmt.Println(string(selectorBytes))
    if err != nil {
        return nil, fmt.Errorf("failed to marshal selector: %v", err)
    }

    resultsIterator, err := ctx.GetStub().GetQueryResult(string(selectorBytes))
    if err != nil {
        return nil, fmt.Errorf("failed to get query result: %v", err)
    }
    defer resultsIterator.Close()


    var users []*User
    for resultsIterator.HasNext() {
        queryResponse, err := resultsIterator.Next()
        if err != nil {
            return nil, fmt.Errorf("failed to iterate query results: %v", err)
        }

        var user User
        if err := json.Unmarshal(queryResponse.Value, &user); err != nil {
            return nil, fmt.Errorf("failed to unmarshal user: %v", err)
        }

		if(s.compareIds(ctx,user.BankID)){
	        users = append(users, &user)
		}
    }

    return users, nil
}


func (s *SmartContract) QueryBankFull(ctx contractapi.TransactionContextInterface, id, headquarters, pib string, yearFounded int, userId string) ([]*Bank, error) {
    selector := map[string]interface{}{}
    if id != "" {
        selector["id"] = id
    }
    if headquarters != "" {
        selector["headquarters"] = headquarters
    }
    if pib != "" {
        selector["pib"] = pib
    }
    if yearFounded != 0 {
        selector["yearFounded"] = yearFounded
    }
    if userId != "" {
        selector["userIds"] = map[string]interface{}{
            "$elemMatch": map[string]interface{}{
                "$eq": userId,
            },
        }
    }

    query := map[string]interface{}{
        "selector": selector,
    }

    queryBytes, err := json.Marshal(query)
    if err != nil {
        return nil, fmt.Errorf("failed to marshal query: %v", err)
    }

    resultsIterator, err := ctx.GetStub().GetQueryResult(string(queryBytes))
    if err != nil {
        return nil, fmt.Errorf("failed to get query result: %v", err)
    }
    defer resultsIterator.Close()

    var banks []*Bank
    for resultsIterator.HasNext() {
        queryResponse, err := resultsIterator.Next()
        if err != nil {
            return nil, fmt.Errorf("failed to iterate query results: %v", err)
        }

        var bank Bank
        if err := json.Unmarshal(queryResponse.Value, &bank); err != nil {
            return nil, fmt.Errorf("failed to unmarshal bank: %v", err)
        }


		if(s.compareIds(ctx,bank.ID)){
        	banks = append(banks, &bank)
		}
    }

    return banks, nil
}

func (s *SmartContract) QueryAccountFull(ctx contractapi.TransactionContextInterface, id string, amount float64, currency, userId string) ([]*Account, error) {
    selector := map[string]interface{}{}
    if id != "" {
        selector["id"] = id
    }
    if amount != 0 {
        selector["amount"] = amount
    }
    if currency != "" {
        selector["currency"] = currency
    }
    if userId != "" {
        selector["userId"] = userId
    }

    query := map[string]interface{}{
        "selector": selector,
    }

    queryBytes, err := json.Marshal(query)
    if err != nil {
        return nil, fmt.Errorf("failed to marshal query: %v", err)
    }

    resultsIterator, err := ctx.GetStub().GetQueryResult(string(queryBytes))
    if err != nil {
        return nil, fmt.Errorf("failed to get query result: %v", err)
    }
    defer resultsIterator.Close()

    var accounts []*Account
    for resultsIterator.HasNext() {
        queryResponse, err := resultsIterator.Next()
        if err != nil {
            return nil, fmt.Errorf("failed to iterate query results: %v", err)
        }

        var account Account
        if err := json.Unmarshal(queryResponse.Value, &account); err != nil {
            return nil, fmt.Errorf("failed to unmarshal account: %v", err)
        }

		bankId, err := s.GetBankIdFromAccount(ctx, account.UserID)
		if err != nil {
			continue
		}

		if(s.compareIds(ctx,bankId)){
        	accounts = append(accounts, &account)
		}
    }

    return accounts, nil
}

func ExtractOrgNumber(input string) (string, error) {
    pattern := `Org(\d+)MSP`

    regex := regexp.MustCompile(pattern)

    matches := regex.FindStringSubmatch(input)

    if len(matches) < 2 {
        return "", errors.New("no OrgMSP found")
    }

	return matches[1], nil
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
