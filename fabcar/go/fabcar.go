/*
Copyright 2020 IBM All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	// "encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/gateway"
)

func main() {
	os.Setenv("DISCOVERY_AS_LOCALHOST", "true")
	wallet, err := gateway.NewFileSystemWallet("wallet")
	if err != nil {
		fmt.Printf("Failed to create wallet: %s\n", err)
		os.Exit(1)
	}

	if !wallet.Exists("appUser") {
		err = populateWallet(wallet)
		if err != nil {
			fmt.Printf("Failed to populate wallet contents: %s\n", err)
			os.Exit(1)
		}
	}

	ccpPath := filepath.Join(
		"..",
		"..",
		"test-network",
		"organizations",
		"peerOrganizations",
		"org1.example.com",
		"connection-org1.yaml",
	)

	gw, err := gateway.Connect(
		gateway.WithConfig(config.FromFile(filepath.Clean(ccpPath))),
		gateway.WithIdentity(wallet, "appUser"),
	)
	if err != nil {
		fmt.Printf("Failed to connect to gateway: %s\n", err)
		os.Exit(1)
	}
	defer gw.Close()

	network, err := gw.GetNetwork("mychannel")
	if err != nil {
		fmt.Printf("Failed to get network: %s\n", err)
		os.Exit(1)
	}

	contract := network.GetContract("basic")

	_, err = contract.SubmitTransaction("InitLedger")
	if err != nil {
		fmt.Printf("Failed to submit transaction: %s\n", err)
		os.Exit(1)
	}
	fmt.Println("--------------------------------")
	fmt.Println("------------------131231-------------")

	// cert, err := contract.EvaluateTransaction("QueryAllBanks")
	// if err != nil {
	// 	fmt.Printf("Failed to evaluate transaction: %s\n", err)
	// 	os.Exit(1)
	// }
	// fmt.Println(string(cert))
	// fmt.Println("--------------------------------")
	// fmt.Println("--------------------------------")
	// fmt.Println("--------------------------------")

	// result, err = contract.EvaluateTransaction("QueryAllUsers")
	// if err != nil {
	// 	fmt.Printf("Failed to evaluate transaction: %s\n", err)
	// 	os.Exit(1)
	// }
	// fmt.Println(string(result))
	// fmt.Println("--------------------------------")
	// fmt.Println("--------------------------------")
	// result, err = contract.EvaluateTransaction("QueryAllAccounts")
	// if err != nil {
	// 	fmt.Printf("Failed to evaluate transaction: %s\n", err)
	// 	os.Exit(1)
	// }
	// fmt.Println(string(result))




	// _, err = contract.SubmitTransaction("CreateUser", "BANK1", "USER123", "Pera", "Peric", "pera@gmail.com")
	// if err != nil {
	// 	fmt.Printf("Failed to submit transaction: %s\n", err)
	// 	os.Exit(1)
	// }

	fmt.Println("--------------------------------")
	fmt.Println("--------------------------------")
	fmt.Println("--------------------------------")
	fmt.Println("--------------------------------")

	// result, err := contract.EvaluateTransaction("QueryAccount", "1_0")
	// if err != nil {
	// 	fmt.Printf("Failed to evaluate transaction: %s\n", err)
	// 	os.Exit(1)
	// }
	// fmt.Println(string(result))

	// result, err = contract.SubmitTransaction("MakePayment", "ACCOUNT123", "31","RSD")
	// if err != nil {
	// 	fmt.Printf("Failed to evaluate transaction: %s\n", err)
	// 	os.Exit(1)
	// }
	// fmt.Println(string(result))


	// result, err := contract.EvaluateTransaction("CheckCurrencyMatch", "ACCOUNT1", "ACCOUNT123")
	// if err != nil {
	// 	fmt.Printf("Failed to evaluate transaction: %s\n", err)
	// 	os.Exit(1)
	// }
	// fmt.Println(string(result))


	// result, err := contract.EvaluateTransaction("QueryAccount", "1")
	// if err != nil {
	// 	fmt.Printf("Failed to evaluate transaction: %s\n", err)
	// 	os.Exit(1)
	// }
	// fmt.Println(string(result))

	// result, err = contract.EvaluateTransaction("QueryAccount", "2")
	// if err != nil {
	// 	fmt.Printf("Failed to evaluate transaction: %s\n", err)
	// 	os.Exit(1)
	// }
	// fmt.Println(string(result))


	// result, err = contract.SubmitTransaction("TransferBetweenAccounts", "ACCOUNT123", "ACCOUNT2","100","RSD")
	// if err != nil {
	// 	fmt.Printf("Failed to evaluate transaction: %s\n", err)
	// 	os.Exit(1)
	// }
	// fmt.Println(string(result))


	// result, err := contract.EvaluateTransaction("QueryAllUsers")
	// if err != nil {
	// 	fmt.Printf("Failed to evaluate transaction: %s\n", err)
	// 	os.Exit(1)
	// }
	// fmt.Println(string(result))
	// fmt.Println("--------------------------------")
	// fmt.Println("--------------------------------")
	// fmt.Println("---------------000000-----------------")
	// id := "1_1"
	// firstName:= ""
	// lastName:= ""
	// email:= ""
	// bankId:= ""
	// accountId:= ""

	// selector := map[string]interface{}{}
    // if id != "" {
    //     selector["id"] = id
    // }
    // if firstName != "" {
    //     selector["firstName"] = firstName
    // }
    // if lastName != "" {
    //     selector["lastName"] = lastName
    // }
    // if email != "" {
    //     selector["email"] = email
    // }
    // if bankId != "" {
    //     selector["bankId"] = bankId
    // }
    // if accountId != "" {
	// 	selector["accountIds"] = map[string]interface{}{
	// 		"$elemMatch": map[string]interface{}{
	// 			"$eq": accountId,
	// 		},
	// 	}
	// }
    // query := map[string]interface{}{
    //     "selector": selector,
    // }

	// fmt.Println(query)

	// result, err := contract.EvaluateTransaction("QueryUserFull", "1_1", "", "", "", "", "1")
	// if err != nil {
	// 	fmt.Printf("Failed to evaluate transaction: %s\n", err)
	// 	os.Exit(1)
	// }
	// fmt.Println(string(result))
	// fmt.Println("-------------00120000-------------------")
	// fmt.Println("--------------------------------")
	// fmt.Println("--------------------------------")

	result, err := contract.EvaluateTransaction("QueryAccount", "2_3")
	if err != nil {
		fmt.Printf("Failed to evaluate transaction: %s\n", err)
		os.Exit(1)
	}
	fmt.Println(string(result))

	// cardIds := []string{}

	// // Marshal the slice into a JSON string
	// cardIdsJSON, err := json.Marshal(cardIds)
	// if err != nil {
	// 	// Handle error
	// }

	// result, err = contract.SubmitTransaction("CreateAccount", "USER123", "ACCOUNT123", "5", "RSD", string(cardIdsJSON))


	// result, err = contract.SubmitTransaction("createCar", "CAR10", "VW", "Polo", "Grey", "Mary")
	// if err != nil {
	// 	fmt.Printf("Failed to submit transaction: %s\n", err)
	// 	os.Exit(1)
	// }
	// fmt.Println(string(result))

	// result, err = contract.EvaluateTransaction("queryCar", "CAR10")
	// if err != nil {
	// 	fmt.Printf("Failed to evaluate transaction: %s\n", err)
	// 	os.Exit(1)
	// }
	// fmt.Println(string(result))

	// _, err = contract.SubmitTransaction("changeCarOwner", "CAR10", "Archie")
	// if err != nil {
	// 	fmt.Printf("Failed to submit transaction: %s\n", err)
	// 	os.Exit(1)
	// }

	// result, err = contract.EvaluateTransaction("queryCar", "CAR10")
	// if err != nil {
	// 	fmt.Printf("Failed to evaluate transaction: %s\n", err)
	// 	os.Exit(1)
	// }
	// fmt.Println(string(result))
}

func populateWallet(wallet *gateway.Wallet) error {
	credPath := filepath.Join(
		"..",
		"..",
		"test-network",
		"organizations",
		"peerOrganizations",
		"org1.example.com",
		"users",
		"User1@org1.example.com",
		"msp",
	)

	certPath := filepath.Join(credPath, "signcerts", "cert.pem")
	// read the certificate pem
	cert, err := ioutil.ReadFile(filepath.Clean(certPath))
	if err != nil {
		return err
	}

	keyDir := filepath.Join(credPath, "keystore")
	// there's a single file in this dir containing the private key
	files, err := ioutil.ReadDir(keyDir)
	if err != nil {
		return err
	}
	if len(files) != 1 {
		return errors.New("keystore folder should have contain one file")
	}
	keyPath := filepath.Join(keyDir, files[0].Name())
	key, err := ioutil.ReadFile(filepath.Clean(keyPath))
	if err != nil {
		return err
	}

	identity := gateway.NewX509Identity("Org1MSP", string(cert), string(key))
	err = wallet.Put("appUser", identity)
	if err != nil {
		return err
	}
	return nil
}
