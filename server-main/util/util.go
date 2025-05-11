package util

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"nexasecure/logger"

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
  //  fmt.Println(body)
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



func ManageNFTforLogin(Username string) bool {
    if Username == "" {
        fmt.Println("Username cannot be empty")
        return false
    }

    config := logger.NewConfigFromEnv()

    // Initialize logger
    log, err := logger.NewLogger(config)
    if err != nil {
        fmt.Printf("Failed to initialize logger: %v\n", err)
        return false
    }
    defer log.Sync()

    AuthPubAddr := GetUserAuthPubAddr(Username)
    if AuthPubAddr == "" {
        log.Error(fmt.Sprintf("Failed to retrieve AuthPubAddr for user: %s", Username))
        return false
    }

    AuthNft := GetAuthNFT(AuthPubAddr)
    if AuthNft == "" {
        log.Error(fmt.Sprintf("Failed to retrieve AuthNFT for user: %s", Username))
        return false
    }

    // AuthNft := "nft-1746515757065257410"
    // AuthPubAddr := "3059301306072a8648ce3d020106082a8648ce3d03010703420004fb78b4a65dde4c6f6aff0cf6f9db210fcac1e8d3aaba1181b4dc4ab8c4065c533ea69023479bb8b1e9daecb817738c3e8368081e5c6c364abdf584a49770d068"

    isVerified := ValidateAuthNFT(AuthPubAddr, AuthNft)
    if !isVerified {
        log.Error(fmt.Sprintf("Auth NFT verification failed for user: %s", Username))
        return false
    }
    log.Info(fmt.Sprintf("Auth NFT verified successfully for user: %s", Username))

    isRemovedAuthNFT := RemoveAuthNFTfromUser(AuthPubAddr)
    if !isRemovedAuthNFT {
        log.Error(fmt.Sprintf("Failed to remove Auth NFT from user: %s", Username))
        return false
    }
    log.Info(fmt.Sprintf("Auth NFT removed successfully for user: %s", Username))

    // // Revoke Auth NFT

    // // Transfer Auth NFT to server
    // signToken := SignAuthNFT(AuthPubAddr, AuthNft)
    // if signToken == "" {
    //     log.Error(fmt.Sprintf("Failed to sign Auth NFT for user: %s", Username))
    //     return false
    // }
    // log.Info(fmt.Sprintf("Auth NFT signed successfully for user: %s", Username))
    // transferSuccess := TransferAuthNFT(AuthPubAddr, AuthNft, signToken)
    // if !transferSuccess {
    //     log.Error(fmt.Sprintf("Failed to transfer Auth NFT for user: %s", Username))
    //     return false
    // }

    reqpubaddr := GetUserReqPubAddr(Username)
    if reqpubaddr == "" {
        log.Error(fmt.Sprintf("Failed to retrieve ReqPubAddr for user: %s", Username))
        return false
    }
    fmt.Println("reqpubaddr", reqpubaddr)

    NewReqNFT := MintReqNFT(reqpubaddr)
    if NewReqNFT == "" {
        log.Error(fmt.Sprintf("Failed to mint ReqNFT for user: %s", Username))
        return false
    }
    fmt.Println("NewReqNFT", NewReqNFT)

    isStored := StoreReqNFT(reqpubaddr, NewReqNFT)
    if !isStored {
        log.Error(fmt.Sprintf("Error storing REQ NFT for user: %s", Username))
        return false
    }
    fmt.Println("isStored", isStored)

    log.Info(fmt.Sprintf("Initial Login successful for user: %s.", Username))
    return true
}
