package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"server/logger"
	"server/util"
	//"strings"
)

/*****************
Auth Blockchain NFT
*******************/

//store nft handler
func StoreAuthNFTHandler(w http.ResponseWriter, r *http.Request) {

	config := logger.NewConfigFromEnv()

    // Initialize logger
    log, err := logger.NewLogger(config)
    if err != nil {
        fmt.Printf("Failed to initialize logger: %v\n", err)
        os.Exit(1)
    }
    defer log.Sync()

	var Auth struct{
		AuthWallPubAddr string `json:"authwallPubAddr"`
		Nft string `json:"nft_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&Auth); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	fmt.Println("AuthWallPubAddr:", Auth.AuthWallPubAddr)
	fmt.Println("NFT ID:", Auth.Nft)
	conAddr, err := util.GetConAddrByAuthPubKey(Auth.AuthWallPubAddr)
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
	responce, error := SendMessageToClient(conAddr, "setauthnft "+ Auth.Nft)
	if error != nil {
		http.Error(w, "Error sending message to client: "+error.Error(), http.StatusInternalServerError)
		log.Error("Error sending message to client: " + error.Error())
		return
	}
	if responce == "NFT_AUTH stored" {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]bool{
			"nftstored": true,
		})
	} else {
		http.Error(w, "Failed to store NFT", http.StatusInternalServerError)
	}
	log.Info("Auth NFT stored in " + Auth.AuthWallPubAddr)
}
//get nft handler
func GetAuthNFTHandler(w http.ResponseWriter, r *http.Request) {
	//fmt.Println("GetAuthNFTHandler called")
	config := logger.NewConfigFromEnv()
	log, err := logger.NewLogger(config)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer log.Sync()

	var Auth struct {
		AuthWallPubAddr string `json:"authwallPubAddr"`
	}
	//fmt.Println(Auth.AuthWallPubAddr)
	if err := json.NewDecoder(r.Body).Decode(&Auth); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	fmt.Println("AuthWallPubAddr:", Auth.AuthWallPubAddr)
	conAddr, err := util.GetConAddrByAuthPubKey(Auth.AuthWallPubAddr)
	if err != nil {
		http.Error(w, "Error retrieving wallet address", http.StatusInternalServerError)
		return
	}
	fmt.Println("ConAddr:", conAddr)
	response, err := SendMessageToClient(conAddr, "getauthnft")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	//print("Response from client:", response)
	// fmt.Println("Response from client:", response)
	// // Preprocess the response to extract the NFT ID
	// var nftID string
	// parts := strings.Fields(response)
	// if len(parts) != 2 || parts[0] != "NFT_AUTH " {
	// 	http.Error(w, "Failed to parse NFT ID", http.StatusInternalServerError)
	// 	return
	// }
	// nftID = parts[1]
	//fmt.Println("NFT ID:", response)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"nft": response,
	})
}

//remove nft handler
func RemoveAuthNFTHandler(w http.ResponseWriter, r *http.Request) {
	config := logger.NewConfigFromEnv()

	// Initialize logger
	log, err := logger.NewLogger(config)
	if err != nil {
		fmt.Printf("Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer log.Sync()

	var Auth struct {
		AuthWallPubAddr string `json:"authwallPubAddr"`
	}
	if err := json.NewDecoder(r.Body).Decode(&Auth); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	conAddr, err := util.GetConAddrByAuthPubKey(Auth.AuthWallPubAddr)
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
	var response string
	response, err = SendMessageToClient(conAddr, "removeauthnft")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if response == "NFT_REMOVED" {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]bool{
			"removereqnft": true,
		})
	} else {
		http.Error(w, "Failed to remove NFT", http.StatusInternalServerError)
	}
}

