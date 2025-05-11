package util

import (
	"encoding/json"
	"fmt"
	"net/http"
	"nexasecure/logger"
	"os"
)

var log logger.Logger


//init logger
func init() {
	config := logger.NewConfigFromEnv()
	var logerr error
	log, logerr = logger.NewLogger(config)
	if logerr != nil {
		fmt.Printf("Failed to initialize logger: %v\n", logerr)
		os.Exit(1)
	}
}


//get user auth pub addr
func GetUserAuthPubAddr(username string) string {

	MainServerUrl := "http://127.0.0.1:18909" //os.Getenv("MAIN_SERVER_URL")
	ApiEndpoint := "/getauthpubaddr"

	url := MainServerUrl + ApiEndpoint
	payload := map[string]string{"username": username}
	jsonData, err := json.Marshal(payload)
	if err != nil {
		log.Error(fmt.Sprintf("Error marshalling JSON: %v", err))
		return ""
	}
	headers := map[string]string{
		"Content-Type": "application/json",
	}
	//fmt.Println("url", url)
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
	authPubAddr, exists := response["authpubaddr"]
	if !exists {
		log.Error("authpubaddr not found in response")
		return ""
	}
	log.Info("AuthPubAddr successfully retrieved")
	return authPubAddr
	
}



//send socket server and get Auth NFT
func GetAuthNFT(authPubAddr string) string {
	// SocketServerUrl := "http://localhost:10081"//os.Getenv("SOCKET_SERVER_URL")
	// ApiEndpoint := "/getauthnft"
	url := "http://localhost:10081/getauthnft"//SocketServerUrl + ApiEndpoint
	payload := map[string]string{"authwallPubAddr": authPubAddr}
	jsonData, err := json.Marshal(payload)
	//fmt.Println("jsonData", jsonData)
	if err != nil {
		log.Error(fmt.Sprintf("Error marshalling JSON: %v", err))
		return ""
	}
	headers := map[string]string{
		"Content-Type": "application/json",
	}
	//fmt.Println(payload)
	statusCode, responseBody, err := MakeAPICall("POST", url, headers, jsonData)
	//fmt.Println("responseBody", string(responseBody))
	if err != nil {
		log.Error(fmt.Sprintf("Error making API call: %v", err))
		return ""
	}
	if statusCode != http.StatusOK {
		log.Error(fmt.Sprintf("API call failed with status: %d, response: %s", statusCode, responseBody))
		return ""
	}	
	var response map[string]string
	
	err = json.Unmarshal(responseBody, &response)
	if err != nil {
		log.Error(fmt.Sprintf("Error unmarshalling JSON: %v", err))
		return ""
	}
	authNFT, exists := response["nft"]
	if !exists {
		log.Error("authnft not found in response")
		return ""
	}
	log.Info("AuthNFT successfully retrieved")
	return authNFT
}

//verify Auth NFT
func ValidateAuthNFT(authPubAddr string, authNFT string) bool {
	SocketServerUrl := "http://localhost:18080"//os.Getenv("SOCKET_SERVER_URL")
	ApiEndpoint := "/authnft/validate"
	url := SocketServerUrl + ApiEndpoint
	payload := map[string]string{
		"address": authPubAddr, 
		"nft_id": authNFT,
	}
	jsonData, err := json.Marshal(payload)
	if err != nil {
		log.Error(fmt.Sprintf("Error marshalling JSON: %v", err))
		return false
	}
	// fmt.Println(payload)
	// fmt.Println(jsonData)
	headers := map[string]string{
		"Content-Type": "application/json",
	}
	//fmt.Println("send api call  to validate nft to url", url)
	statusCode, responseBody, err := MakeAPICall("POST", url, headers, jsonData)
	if err != nil {
		log.Error(fmt.Sprintf("Error making API call: %v", err))
		return false
	}
	// fmt.Println(responseBody)
	if statusCode != http.StatusOK {
		log.Error(fmt.Sprintf("API call failed with status: %d, response: %s", statusCode, responseBody))
		return false
	}	
	var response map[string]interface{}
	err = json.Unmarshal(responseBody, &response)
	if err != nil {
		log.Error(fmt.Sprintf("Error unmarshalling JSON: %v", err))
		return false
	}
	fmt.Println("response", response)
	if errorObj, exists := response["error"]; exists {
		log.Error(fmt.Sprintf("API response contains error: %v", errorObj))
		return false
	}

	valid, exists := response["valid"]
	if !exists {
		log.Error("valid not found in response")
		return false
	}

	isValid, ok := valid.(bool)
	if !ok {
		log.Error("valid field is not of type bool")
		return false
	}

	log.Info("AuthNFT successfully verified")
	return isValid
}

//sign the auth NFT and return the signed token
func SignAuthNFT(authPubAddr string, authNFT string) string {
	SocketServerUrl := "http://localhost:10081"//os.Getenv("SOCKET_SERVER_URL")
	endpoint := "/signauthwallet"
	url := SocketServerUrl + endpoint
	payload := map[string]string{"address": authPubAddr, "nft_id": authNFT}
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
	if statusCode != http.StatusOK {
		log.Error(fmt.Sprintf("API call failed with status: %d, response: %s", statusCode, responseBody))
		return ""
	}
	var response map[string]string
	err = json.Unmarshal(responseBody, &response)
	if err != nil {
		log.Error(fmt.Sprintf("Error unmarshalling JSON: %v", err))
		return ""
	}
	signedToken, exists := response["signed_token"]
	if !exists {
		log.Error("signed_token not found in response")
		return ""
	}
	log.Info("AuthNFT successfully signed and token retrieved")
	return signedToken
}

//trasfer Auth NFT to server
func TransferAuthNFT(authPubAddr string, authNFT string, signedNFTToken string) bool {
	AuthBlockchainUrl := os.Getenv("AUTH_BLOCKCHAIN_URL")
	endpoint := "/reqnft/transfer"
	url := AuthBlockchainUrl + endpoint
	payload := map[string]string{"sender": authPubAddr, "nft_id": authNFT, "signed_nfttoken": signedNFTToken}
	jsonData, err := json.Marshal(payload)
	if err != nil {
		log.Error(fmt.Sprintf("Error marshalling JSON: %v", err))
		return false
	}
	headers := map[string]string{
		"Content-Type": "application/json",
	}
	statusCode, responseBody, err := MakeAPICall("POST", url, headers, jsonData)
	if err != nil {
		log.Error(fmt.Sprintf("Error making API call: %v", err))
		return false
	}
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
	authNFTTransferred, exists := response["transferred"]
	if !exists {
		log.Error("transferred not found in response")
		return false
	}
	log.Info("AuthNFT successfully transferred")
	return authNFTTransferred
}

//verify req pub addr and auth pub addr 
func VerifyReqPubAddr(ReqPubAddr1 string, ReqPubAddr2 string) bool {

	if ReqPubAddr1 == "" || ReqPubAddr2 == "" {
		log.Error("One or both public addresses are empty")
		return false
	}
	if ReqPubAddr1 != ReqPubAddr2 {
		log.Error("Public addresses do not match")
		return false
	}
	log.Info("Public addresses verified successfully")
	return true
}

func RemoveAuthNFTfromUser(authPubAddr string) bool {
	//remove nft from user
	url := "http://localhost:10081/removeauthnft"
	payload := map[string]string{"authwallPubAddr": authPubAddr}
	jsonData, err := json.Marshal(payload)
	if err != nil {
		log.Error(fmt.Sprintf("Error marshalling JSON: %v", err))
		return false
	}
	headers := map[string]string{
		"Content-Type": "application/json",
	}
	statusCode, responseBody, err := MakeAPICall("POST", url, headers, jsonData)
	if err != nil {
		log.Error(fmt.Sprintf("Error making API call: %v", err))
		return false
	}
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

	removed, exists := response["removereqnft"]
	if !exists {
		log.Error("removereqnft not found in response")
		return false
	}

	if removed {
		log.Info("AuthNFT successfully removed from user")
		return true
	}

	log.Error("Failed to remove AuthNFT from user")
	return false
}