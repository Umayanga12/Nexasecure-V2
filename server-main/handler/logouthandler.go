package handler

import (
	"fmt"
	"nexasecure/logger"
	"os"
	"nexasecure/util"
	"github.com/gin-gonic/gin"
)


func LogoutHandler(c *gin.Context) {

	config := logger.NewConfigFromEnv()

	// Initialize logger
	log, err := logger.NewLogger(config)
	if err != nil {
		fmt.Printf("Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer log.Sync()

	var UserData struct {
		Username string `json:"username"`
	}

	if err := c.ShouldBindJSON(&UserData); err != nil {
		c.JSON(400, gin.H{"error": "Invalid JSON payload"})
		log.Error("Invalid JSON payload: %v", err)
		return
	}
	if UserData.Username == "" {
		c.JSON(400, gin.H{"error": "Username is required"})
		log.Error("Username is required")
		return
	}

	//get user authpub addr
	authpubaddr := util.GetUserAuthPubAddr(UserData.Username)
	if authpubaddr == "" {
		log.Error(fmt.Sprintf("Failed to retrieve AuthPubAddr for user: %s", UserData.Username))
		c.JSON(500, gin.H{"error": "Failed to retrieve AuthPubAddr"})
		return
	}
	log.Info(fmt.Sprintf("AuthPubAddr retrieved successfully for user: %s", UserData.Username))
	//get user reqpub addr
	reqpubaddr := util.GetUserReqPubAddr(UserData.Username)
	if reqpubaddr == "" {
		log.Error(fmt.Sprintf("Failed to retrieve ReqPubAddr for user: %s", UserData.Username))
		c.JSON(500, gin.H{"error": "Failed to retrieve ReqPubAddr"})
		return
	}
	log.Info(fmt.Sprintf("ReqPubAddr retrieved successfully for user: %s", UserData.Username))
	//get reqnft
	reqnft := util.GetReqNFT(reqpubaddr)
	if reqnft == "" {
		log.Error(fmt.Sprintf("Failed to retrieve ReqNFT for user: %s", UserData.Username))
		c.JSON(500, gin.H{"error": "Failed to retrieve ReqNFT"})
		return
	}
	fmt.Println(reqnft)
	log.Info(fmt.Sprintf("Req NFT retrieved successfully for user: %s", UserData.Username))
	//validate reqnft
	isVerified := util.ValidateReqNFT(reqpubaddr, reqnft)
	if !isVerified {
		log.Error(fmt.Sprintf("Req NFT verification failed for user: %s", UserData.Username))
		c.JSON(400, gin.H{"error": "Req NFT verification failed"})
		return
	}
	log.Info(fmt.Sprintf("Req NFT verified successfully for user: %s", UserData.Username))
	//sign reqnft
	//transfer reqnft 
	//remove reqnft from user
	isRemoved := util.RemoveReqNFTfromUser(reqpubaddr)
	if !isRemoved {
		log.Error(fmt.Sprintf("Failed to remove Req NFT from user: %s", UserData.Username))
		c.JSON(500, gin.H{"error": "Failed to remove Req NFT"})
		return
	}
	log.Info(fmt.Sprintf("Req NFT removed successfully for user: %s", UserData.Username))
	//mint new authnft	
	newAuthNFT := util.MintNewAuthNFT(authpubaddr)
	if newAuthNFT == "" {
		log.Error(fmt.Sprintf("Failed to mint new Auth NFT for user: %s", UserData.Username))
		c.JSON(500, gin.H{"error": "Failed to mint new Auth NFT"})
		return
	}
	log.Info(fmt.Sprintf("New Auth NFT minted successfully for user: %s", UserData.Username))
	//store int wallet
	isStored := util.StoreAuthNFT(authpubaddr,newAuthNFT)
	if !isStored {
		log.Error(fmt.Sprintf("Failed to store new Auth NFT for user: %s", UserData.Username))
		return
	}
	log.Info(fmt.Sprintf("New Auth NFT stored successfully for user: %s", UserData.Username))

	c.JSON(200, gin.H{"message": "User logged out successfully"})
}