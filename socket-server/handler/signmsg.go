package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"server/logger"
	"server/util"
)

func SignAuthwalletHandler(w http.ResponseWriter, r *http.Request) {
	config := logger.NewConfigFromEnv()

	// Initialize logger
	log, err := logger.NewLogger(config)
	if err != nil {
		fmt.Printf("Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer log.Sync()
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	var SignNFT struct {
		AuthWallPubAddr string `json:"authwallPubAddr"`
		Nft            string `json:"nft_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&SignNFT); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	conAddr, err := util.GetConAddrByAuthPubKey(SignNFT.AuthWallPubAddr)
	if err != nil {
		http.Error(w, "Error retrieving connection address: "+err.Error(), http.StatusInternalServerError)
		log.Error("Error retrieving connection address: " + err.Error())
		return
	}
	if conAddr == "" {
		http.Error(w, "No wallet found for the given request address", http.StatusNotFound)
		log.Error("No wallet found for the given request address")
		return
	}
	response, err := SendMessageToClient(conAddr, "signauthmsg "+SignNFT.Nft)
	if err != nil {
		http.Error(w, "Error sending message to client: "+err.Error(), http.StatusInternalServerError)
		log.Error("Error sending message to client: " + err.Error())
		return
	}
	log.Info("NFT sent to client for signing")

	// Send the response back to the HTTP client as a JSON object
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	jsonResponse := map[string]string{
		"response": response,
	}
	json.NewEncoder(w).Encode(jsonResponse)
}

//remove nft handler
func SignReqwalletHandler(w http.ResponseWriter, r *http.Request) {
	config := logger.NewConfigFromEnv()

	// Initialize logger
	log, err := logger.NewLogger(config)
	if err != nil {
		fmt.Printf("Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer log.Sync()
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	var SignNFT struct {
		RequestWallPubAddr string `json:"requestwallPubAddr"`
		Nft               string `json:"nft_id"`	
	}
	if err := json.NewDecoder(r.Body).Decode(&SignNFT); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	conAddr, err := util.GetConAddrByReqPubKey(SignNFT.RequestWallPubAddr)
	if err != nil {
		http.Error(w, "Error retrieving wallet address: "+err.Error(), http.StatusInternalServerError)
		log.Error("Error retrieving wallet address: " + err.Error())
		return
	}
	if conAddr == "" {
		http.Error(w, "No wallet found for the given request address", http.StatusNotFound)
		log.Error("No wallet found for the given request address")
		return
	}
	response, err := SendMessageToClient(conAddr, "signreqmsg "+SignNFT.Nft)
	if err != nil {
		http.Error(w, "Error sending message to client: "+err.Error(), http.StatusInternalServerError)
		log.Error("Error sending message to client: " + err.Error())
		return
	}
	log.Info("NFT sent to client for signing")
	// Send the response back to the HTTP client
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(response))
	
}
