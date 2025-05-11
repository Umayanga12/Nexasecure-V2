package handler

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

// VerifySignature validates a message signature against a public address
func VerifySignature(address string, message string, signature string) (bool, error) {
	// 1. Validate address format
	if valid, err := ValidatePublicAddress(address); !valid {
		return false, fmt.Errorf("invalid address: %v", err)
	}

	// 2. Retrieve public key from blockchain
	blockchain.mutex.Lock()
	defer blockchain.mutex.Unlock()
	
	wallet, exists := blockchain.Wallets[address]
	if !exists {
		return false, errors.New("wallet not found")
	}

	// 3. Decode signature
	sigBytes, err := hex.DecodeString(signature)
	if err != nil {
		return false, fmt.Errorf("invalid signature format: %v", err)
	}

	// 4. Hash the message
	hash := sha256.Sum256([]byte(message))

	// 5. Verify the signature
	valid := ecdsa.VerifyASN1(wallet, hash[:], sigBytes)
	return valid, nil
}

// SignatureValidationRequest represents the validation request
type SignatureValidationRequest struct {
	Address   string `json:"address"`
	Message   string `json:"message"`
	Signature string `json:"signature"`
}

// SignatureValidationHandler HTTP endpoint
func SignatureValidationHandler(w http.ResponseWriter, r *http.Request) {
	var req SignatureValidationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	valid, err := VerifySignature(req.Address, req.Message, req.Signature)
	response := map[string]interface{}{
		"valid":  valid,
		"error":  nil,
	}

	if err != nil {
		response["error"] = err.Error()
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// // Add to main function
// func main() {
// 	// Previous endpoints
// 	http.HandleFunc("/verify/signature", SignatureValidationHandler)
// 	// ... rest of main function
// }

// ValidatePublicAddress checks if a public address is properly formatted and exists in the system
func ValidatePublicAddress(address string) (bool, error) {
	fmt.Println("Starting ValidatePublicAddress with address:", address)

	// Normalize the address (e.g., trim spaces and convert to lowercase)
	normalizedAddress := strings.TrimSpace(strings.ToLower(address))
	fmt.Println("Normalized address:", normalizedAddress)

	// Check basic format requirements
	if len(normalizedAddress) < 130 { // Length for uncompressed ECDSA P-256 public key in hex
		fmt.Println("Invalid address length:", len(normalizedAddress))
		return false, errors.New("invalid address length")
	}

	// Validate hexadecimal format
	if _, err := hex.DecodeString(normalizedAddress); err != nil {
		fmt.Println("Invalid hexadecimal format for address:", normalizedAddress, "Error:", err)
		return false, errors.New("invalid hexadecimal format")
	}

	blockchain.mutex.Lock()
	defer blockchain.mutex.Unlock()
	fmt.Println("Acquired blockchain mutex lock")

	// Check if address exists in wallets
	_, exists := blockchain.Wallets[normalizedAddress]
	if !exists {
		fmt.Println("Address not found in registry:", normalizedAddress)
		return false, errors.New("address not found in registry")
	}

	fmt.Println("Address is valid and exists in registry:", normalizedAddress)
	return true, nil
}

// HTTP handler for address validation
func ValidateAddressHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("ValidateAddressHandler called")

	var req struct {
		Address string `json:"address"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		fmt.Println("Error decoding JSON payload:", err)
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	fmt.Println("Request payload:", req)

	if req.Address == "" {
		fmt.Println("Missing address field in request body")
		http.Error(w, "Missing address field in request body", http.StatusBadRequest)
		return
	}

	fmt.Println("Validating address:", req.Address)
	valid, err := ValidatePublicAddress(req.Address)
	fmt.Println("Validation result - valid:", valid, "error:", err)

	response := map[string]interface{}{
		"valid": valid,
	}

	if err != nil {
		response["error"] = err.Error()
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		fmt.Println("Error encoding response:", err)
	}
}

// // Add to main function
// func main() {
// 	// Previous endpoints
// 	http.HandleFunc("/wallet/validate", ValidateAddressHandler)
// 	// ... rest of main function
// }