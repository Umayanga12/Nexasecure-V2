package main

import (
	"blockchain-api/handler"
	"blockchain-api/logger"
	"fmt"
	"net/http"
	"os"
)

func startServer(port string, log logger.Logger) {
	// Generate host wallet
	_, hostWalletAddress, err := handler.GenerateWallet()
	if err != nil {
		log.Error("Failed to generate host wallet", err)
		os.Exit(1)
	}

	// Assign host wallet to blockchain
	blockchain := handler.NewBlockchain()
	blockchain.HostWallet = hostWalletAddress

	log.Info("Host wallet created", "address", hostWalletAddress)

	mux := http.NewServeMux()
	mux.HandleFunc("/", handler.RootHandler)
	mux.HandleFunc("/authnft/create", handler.CreateNFTHandler)
	mux.HandleFunc("/authnft/transfer", handler.TransferNFTHandler)
	mux.HandleFunc("/authnft/owner", handler.GetNFTOwnerHandler)
	mux.HandleFunc("/authnft/burn", handler.BurnNFTHandler)
	mux.HandleFunc("/authnft/validate", handler.ValidateNFTOwnerHandler)
	mux.HandleFunc("/authwallet/new", handler.GenerateWalletHandler)
	mux.HandleFunc("/authwallet/pubaddrval", handler.ValidateAddressHandler)
	mux.HandleFunc("/verify/signature", handler.SignatureValidationHandler)

	go func() {
		fmt.Printf("Server running on port %s\n", port)
		log.Info(fmt.Sprintf("auth blockchain Server running on port %s", port))
		if err := http.ListenAndServe(":"+port, mux); err != nil {
			log.Error(fmt.Sprintf("Failed to start server on port %s: %v", port, err))
			os.Exit(1)
		}
	}()
}

func main() {
	config := logger.NewConfigFromEnv()
	log, err := logger.NewLogger(config)
	if err != nil {
		fmt.Printf("Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer log.Sync()

	// Start two server instances on different ports
	startServer("18080", log)

	// Block the main goroutine to keep the servers running
	select {}
}
