package main

import (
	"context"
	"gas-price-notifier/config"
	"gas-price-notifier/notifier"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"
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

	ctx, cancel := context.WithCancel(context.TODO())

	// Attempt to catch interrupt event(s)
	// so that graceful shutdown can be performed
	interruptChan := make(chan os.Signal, 1)
	signal.Notify(interruptChan, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL)

	go func() {

		// To be invoked when returning from this
		// go rountine's execution scope
		defer func() {

			// Stopping process
			log.Printf("\n[✅] Gracefully shut down `gasz`\n")
			os.Exit(0)

		}()

		<-interruptChan
		// Once signal is received, it'll cancel context
		// so that worker(s) can stop working on what they're doing now
		cancel()

		// Giving worker(s) 3 seconds, before forcing shutdown
		<-time.After(time.Second * time.Duration(3))
		return

	}()

	notifier.Start(ctx)

}
