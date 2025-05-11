package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"server/logger"
)

// MakeAPICall sends an HTTP request and returns the status code, response body, and error if any.
func MakeAPICall(method, url string, headers map[string]string, body []byte) (int, []byte, error) {
    req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
    if err != nil {
        return 0, nil, err
    }

    for key, value := range headers {
        req.Header.Set(key, value)
    }

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return 0, nil, err
    }
    defer resp.Body.Close()

    respBody, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return 0, nil, err
    }

    return resp.StatusCode, respBody, nil
}

//validate the request blockchain wallet private key
func ValidateReqaddr(reqaddr string) bool {
    config := logger.NewConfigFromEnv()

    // Initialize logger
    log, err := logger.NewLogger(config)
    if err != nil {
        fmt.Printf("Failed to initialize logger: %v\n", err)
        os.Exit(1)
    }
    defer log.Sync()

    requestaddrvfalidaton := "http://localhost:18085/reqwallet/pubaddrval"

    if requestaddrvfalidaton == "" {
        log.Error("REQUEST_BLOCKCHAIN_SERVER is not set in the environment")
        return false
    }

    requestPayload := map[string]string{
        "address": reqaddr,
    }

    requestjsonPayload, _ := json.Marshal(requestPayload)

    requestheaders := map[string]string{
        "Content-Type": "application/json",
    }

    statusCode, body, err := MakeAPICall("POST", requestaddrvfalidaton, requestheaders, requestjsonPayload)
    if err != nil || statusCode != http.StatusOK {
        log.Error(err.Error())
        return false
    }

    var response map[string]any
    if err := json.Unmarshal(body, &response); err != nil {
        log.Error(err.Error())
        return false
    }

    if valid, ok := response["valid"].(bool); ok && valid {
        return true
    }

    return false
}

//validate the auth blockchain wallet public key
func ValidateAuthaddr(authaddr string) bool {
    config := logger.NewConfigFromEnv()

    // Initialize logger
    log, err := logger.NewLogger(config)
    if err != nil {
        fmt.Printf("Failed to initialize logger: %v\n", err)
        os.Exit(1)
    }
    defer log.Sync()

    authpubaddrvalidation := "http://localhost:18080/authwallet/pubaddrval"

    authPayload := map[string]string{
        "address": authaddr,
    }

    authjsonPayload, _ := json.Marshal(authPayload)

    authheaders := map[string]string{
        "Content-Type": "application/json",
    }

    statusCode, body, err := MakeAPICall("POST", authpubaddrvalidation, authheaders, authjsonPayload)
    if err != nil || statusCode != http.StatusOK {
        log.Error(err.Error())
        return false
    }

    var response map[string]any
    if err := json.Unmarshal(body, &response); err != nil {
        log.Error(err.Error())
        return false
    }

    if valid, ok := response["valid"].(bool); ok && valid {
        return true
    }

    return false
}


