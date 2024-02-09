package src

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/gateway"
)

func getContract(wallet *gateway.Wallet, bank string) (*gateway.Contract, error) {
	ccpPath := filepath.Join(
		"..",
		"..",
		"test-network",
		"organizations",
		"peerOrganizations",
		fmt.Sprintf("org%s.example.com", bank),
		fmt.Sprintf("connection-org%s.yaml", bank),
	)

	gw, err := gateway.Connect(
		gateway.WithConfig(config.FromFile(filepath.Clean(ccpPath))),
		gateway.WithIdentity(wallet, fmt.Sprintf("org%s", bank)),
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

func populateWallet(wallet *gateway.Wallet, bank string) error {
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

	err = wallet.Put(fmt.Sprintf("org%s", bank), identity)
	if err != nil {
		return err
	}
	return nil
}
