package util

// import (
// 	"fmt"
// 	"os"
// 	"server/logger"
// )

// //validate signature
// func ValSignNFT(command string, client string){
// 	config := logger.NewConfigFromEnv()
// 	logger, err := logger.NewLogger(config)
// 	if err != nil {
// 		fmt.Printf("Failed to initialize logger: %v\n", err)
// 		os.Exit(1)
// 	}
// 	defer logger.Sync()

// 	apiUrl := os.Getenv("AUTH_BLOCKCHAIN_SERVER")
// 	APIurlEndpoint := apiUrl + "/reqnft/create"
// }

