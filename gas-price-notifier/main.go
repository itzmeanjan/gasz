package main

import (
	"gas-price-notifier/config"
	"gas-price-notifier/notifier"
	"log"
	"path/filepath"
)

// Entry point of program
func main() {
	absPath, err := filepath.Abs(".env")
	if err != nil {
		log.Fatalf("[!] Failed find `.env` file : %s\n", err.Error())
	}

	if err := config.Load(absPath); err != nil {
		log.Fatalf("[!] Failed load `.env` file : %s\n", err.Error())
	}

	notifier.Start()
}
