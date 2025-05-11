package handlers

import (
	"auth-server/utils"
	"encoding/json"
	"net/http"
)


func GetAuthPubAddr(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var userDetails struct {
		Username string `json:"username"`
	}
	if err := json.NewDecoder(r.Body).Decode(&userDetails); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	authPubAddr, err := utils.GetAuthPubAddrFromUsername(userDetails.Username)
	if err != nil {
		http.Error(w, "Failed to fetch authpubaddr", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := map[string]string{
		"authpubaddr": authPubAddr,
	}
	json.NewEncoder(w).Encode(response)
}

func GetReqPubAddr(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var userDetails struct {
		Username string `json:"username"`
	}
	if err := json.NewDecoder(r.Body).Decode(&userDetails); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	reqPubAddr, err := utils.GetReqPubAddrFromUsername(userDetails.Username)
	if err != nil {
		http.Error(w, "Failed to fetch reqpubaddr", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := map[string]string{
		"reqpubaddr": reqPubAddr,
	}
	json.NewEncoder(w).Encode(response)
}