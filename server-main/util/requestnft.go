package util

import (
	"encoding/json"
	"fmt"
	"net/http"
)

//Mint new Req nft
func MintReqNFT(ReqPubAddr string) string {
	ReqBlockchainUrl := "http://localhost:18085"//os.Getenv("REQ_BLOCKCHAIN_URL")
	endpoint := "/reqnft/create"
	url := ReqBlockchainUrl + endpoint
	payload := map[string]string{"owner": ReqPubAddr}
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
	if statusCode != http.StatusOK && statusCode != http.StatusCreated  {
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
	log.Info("Req NFT successfully minted")
	return newNFT
}


//Store in client wallet 
func StoreReqNFT(ReqPubAddr string, newNFT string) bool {
	SocketServerUrl := "http://localhost:10081" //os.Getenv("SOCKET_SERVER_URL")
	endpoint := "/setreqnft"
	url := SocketServerUrl + endpoint
	payload := map[string]string{"requestwallPubAddr": ReqPubAddr, "nft_id": newNFT}
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
	fmt.Println("responseBody", string(responseBody))
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
	nftStored, exists := response["nftstored"]
	if !exists || !nftStored {
		log.Error("nftstored not found or false in response")
		return false
	}
	log.Info("Req NFT successfully stored")
	return true
}

func GetUserReqPubAddr(username string) string {

	MainServerUrl := "http://127.0.0.1:18909"//os.Getenv("MAIN_SERVER_URL")
	ApiEndpoint := "/getreqpubaddr"

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
	reqPubAddr, exists := response["reqpubaddr"]
	if !exists {
		log.Error("req_pub_addr not found in response")
		return ""
	}
	log.Info("ReqPubAddr successfully retrieved")
	return reqPubAddr
}

func GetReqNFT(authPubAddr string) string {
	url := "http://localhost:10081/getreqnft"//SocketServerUrl + ApiEndpoint
	payload := map[string]string{"requestwallPubAddr": authPubAddr}
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
	//fmt.Println(responseBody)
	if statusCode != http.StatusOK {
		log.Error(fmt.Sprintf("API call failed with status: %d, response: %s", statusCode, responseBody))
		return ""
	}
	if len(responseBody) == 0 {
		log.Error("API response body is empty")
		return ""
	}
	var response map[string]string
	fmt.Println("responseBody", string(responseBody))
	err = json.Unmarshal(responseBody, &response)
	// fmt.Println("response", response)
	if err != nil && response == nil {
		log.Error(fmt.Sprintf("Error unmarshalling JSON or response is nil: %v", err))
		return ""
	}
	
	authNFT, exists := response["nft"]
	if !exists {
		log.Error("reqnft not found in response")
		return ""
	}
	log.Info("ReqNFT successfully retrieved")
	return authNFT
}

func RemoveReqNFTfromUser(reqPubAddr string) bool {
	url := "http://localhost:10081/removereqnft" // SocketServerUrl + ApiEndpoint
	payload := map[string]string{"requestwallPubAddr": reqPubAddr}
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
	if !exists || !removed {
		log.Error("removereqnft not found or false in response")
		return false
	}
	log.Info("AuthNFT successfully removed from user")
	return true
}

func ValidateReqNFT(authPubAddr string, authNFT string) bool {
	SocketServerUrl := "http://localhost:18085"//os.Getenv("SOCKET_SERVER_URL")
	ApiEndpoint := "/reqnft/validate"
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
	fmt.Println("response", response)
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