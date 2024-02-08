package main

import (
	"fabcar/src"
	"fmt"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/hyperledger/fabric-sdk-go/pkg/gateway"
)

func main() {
	os.Setenv("DISCOVERY_AS_LOCALHOST", "true")
	wallet, err := gateway.NewFileSystemWallet("wallet")
	if err != nil {
		fmt.Printf("Failed to create wallet: %s\n", err)
		os.Exit(1)
	}

	src.StartServer(&src.Handler{Validator: validator.New(), Wallet: wallet})
}
