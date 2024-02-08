package main

import (
	"log"
	"main/src"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	src.StartServer(&src.Handler{Validator: validator.New()})
}
