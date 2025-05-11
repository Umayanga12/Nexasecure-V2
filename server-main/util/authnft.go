package util

import (
	"encoding/json"
	"net/http"
	"os"
	"fmt"
	"nexasecure/logger"
	
)


func MintNewAuthNFT(userAddr string) string {
	config := logger.NewConfigFromEnv()

	// Initialize logger
	log, err := logger.NewLogger(config)
	if err != nil {
		fmt.Printf("Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer log.Sync()
	baseURL := "http://localhost:18080"//os.Getenv("BASE_URL")
	if baseURL == "" {
		log.Fatal("BASE_URL not set in .env file")
	}

	apiEndpoint := "/authnft/create"
	url := baseURL + apiEndpoint

	payload := map[string]string{"owner": userAddr}
	jsonData, err := json.Marshal(payload)
	if err != nil {
		log.Error(fmt.Sprintf("Error marshalling JSON: %v", err))
		return ""
	}

	headers := map[string]string{
		"Content-Type": "application/json",
	}
	statusCode, responseBody, err := MakeAPICall("POST", url, headers, jsonData)
	if err != nil {
		log.Error(fmt.Sprintf("Error making API call: %v", err))
		return ""
	}
	if statusCode != http.StatusOK && statusCode != http.StatusCreated {
		log.Error(fmt.Sprintf("API call failed with status: %d, response: %s", statusCode, responseBody))
		return ""
	}

	var response map[string]string
	err = json.Unmarshal(responseBody, &response)
	if err != nil {
		log.Error(fmt.Sprintf("Error unmarshalling JSON: %v", err))
		return ""
	}
	// fmt.Println("response", response)
	newNFT, exists := response["nft_id"]
	if !exists {
		log.Error("newnft not found in response")
		return ""
	}
	log.Info("auth NFT successfully minted")
	return newNFT

}


func StoreAuthNFT(AuthPubAddr string, newNFT string) bool {
	SocketServerUrl := "http://localhost:10081" //os.Getenv("SOCKET_SERVER_URL")
	endpoint := "/setauthnft"
	url := SocketServerUrl + endpoint
	fmt.Println("StoreAuthNFT - URL:", url) // Debugging statement

	payload := map[string]string{"authwallPubAddr": AuthPubAddr, "nft_id": newNFT}
	fmt.Println("StoreAuthNFT - Payload:", payload) // Debugging statement

	jsonData, err := json.Marshal(payload)
	if err != nil {
		log.Error(fmt.Sprintf("Error marshalling JSON: %v", err))
		return false
	}
	fmt.Println("StoreAuthNFT - JSON Data:", string(jsonData)) // Debugging statement

	headers := map[string]string{
		"Content-Type": "application/json",
	}
	statusCode, responseBody, err := MakeAPICall("POST", url, headers, jsonData)
	if err != nil {
		log.Error(fmt.Sprintf("Error making API call: %v", err))
		return false
	}
	fmt.Println("StoreAuthNFT - Response Body:", string(responseBody)) // Debugging statement
	fmt.Println("StoreAuthNFT - Status Code:", statusCode) // Debugging statement

	if statusCode != http.StatusOK {
		log.Error(fmt.Sprintf("API call failed with status: %d, response: %s", statusCode, responseBody))
		return false
	}
	var response map[string]bool
	err = json.Unmarshal(responseBody, &response)
	if err != nil {
		log.Error(fmt.Sprintf("Error unmarshalling JSON: %v", err))
		return false
	}
	fmt.Println("StoreAuthNFT - Parsed Response:", response) // Debugging statement

	nftStored, exists := response["nftstored"]
	if !exists || !nftStored {
		log.Error("nftstored not found or false in response")
		return false
	}
	log.Info("Auth NFT successfully stored")
	return true
}
