package src

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/gateway"
)

func enrollUser(wallet *gateway.Wallet, user RegisterUser) error {
	if wallet.Exists(user.UserID) {
		return fmt.Errorf("User %s already enrolled", user.UserID)
	}

	err := populateWallet(wallet, user.BankID, user.UserID)
	if err != nil {
		return fmt.Errorf("Failed to populate wallet contents: %s\n", err)
	}

	contract, err := getContract(wallet, user.BankID, user.UserID)

	if err != nil {
		wallet.Remove(user.UserID)
		return fmt.Errorf("Failed to get contract")
	}

	_, err = contract.SubmitTransaction(
		"CreateUser",
		user.BankID,
		user.UserID,
		user.FirstName,
		user.LastName,
		user.Email,
	)

	if err != nil {
		wallet.Remove(user.UserID)
		return fmt.Errorf("Could not create user %s in bank %s", user.UserID, user.BankID)
	}

	return nil
}

func getContract(wallet *gateway.Wallet, bank string, userId string) (*gateway.Contract, error) {
	ccpPath := filepath.Join(
		"..",
		"..",
		"test-network",
		"organizations",
		"peerOrganizations",
		fmt.Sprintf("org%s.example.com", bank),
		fmt.Sprintf("connection-org%s.yaml", userId),
	)

	gw, err := gateway.Connect(
		gateway.WithConfig(config.FromFile(filepath.Clean(ccpPath))),
		gateway.WithIdentity(wallet, userId),
	)
	if err != nil {
		return nil, fmt.Errorf("Failed to connect to gateway: %s\n", err)
	}
	defer gw.Close()

	network, err := gw.GetNetwork("mychannel")
	if err != nil {
		return nil, fmt.Errorf("Failed to get network: %s\n", err)
	}

	return network.GetContract("basic"), nil
}

func populateWallet(wallet *gateway.Wallet, bank string, userId string) error {
	credPath := filepath.Join(
		"..",
		"..",
		"test-network",
		"organizations",
		"peerOrganizations",
		fmt.Sprintf("org%s.example.com", bank),
		"users",
		fmt.Sprintf("User1@org%s.example.com", bank),
		"msp",
	)

	certPath := filepath.Join(credPath, "signcerts", "cert.pem")
	cert, err := os.ReadFile(filepath.Clean(certPath))
	if err != nil {
		return err
	}

	keyDir := filepath.Join(credPath, "keystore")
	files, err := os.ReadDir(keyDir)
	if err != nil {
		return err
	}
	if len(files) != 1 {
		return errors.New("keystore folder should have contain one file")
	}
	keyPath := filepath.Join(keyDir, files[0].Name())
	key, err := os.ReadFile(filepath.Clean(keyPath))
	if err != nil {
		return err
	}

	identity := gateway.NewX509Identity(fmt.Sprintf("Org%sMSP", bank), string(cert), string(key))

	err = wallet.Put(userId, identity)
	if err != nil {
		return err
	}
	return nil
}
