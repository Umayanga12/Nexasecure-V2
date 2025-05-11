package util

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
)

var (
	ctx         = context.Background()
	redisClient *redis.Client
)

func init() {
	if err := godotenv.Load(); err != nil {
        fmt.Println("Error loading .env file")
    }
    // Parse REDIS_DB (default to 0 on error)
    dbNum := 0
    if dbStr := os.Getenv("REDIS_DB"); dbStr != "" {
        if n, err := strconv.Atoi(dbStr); err == nil {
            dbNum = n
        }
    }

    // Print environment variables for debugging
    fmt.Printf("REDIS_HOST: %s\n", os.Getenv("REDIS_HOST"))
    fmt.Printf("REDIS_PORT: %s\n", os.Getenv("REDIS_PORT"))
    fmt.Printf("REDIS_PASSWORD: %s\n", os.Getenv("REDIS_PASSWORD"))
    fmt.Printf("REDIS_DB: %d\n", dbNum)

    // Build Redis client
    redisClient = redis.NewClient(&redis.Options{
        Addr:     os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT"),
        Password: os.Getenv("REDIS_PASSWORD"),
        DB:       dbNum,
    })
}

// AddWallet stores conAddr, reqPubKey, and authPubKey in Redis.
func AddWallet(conAddr, reqPubKey, authPubKey string) error {
	// Store the wallet data in a hash
	walletKey := "wallet:" + conAddr
	err := redisClient.HSet(ctx, walletKey, map[string]interface{}{
		"reqpubkey":  reqPubKey,
		"authpubkey": authPubKey,
		"conaddr":    conAddr,
	}).Err()

	if err != nil {
		// Log error if storing wallet data fails
		fmt.Printf("Error storing wallet data for conAddr %s: %v\n", conAddr, err)
		return err
	}
	fmt.Printf("Successfully stored wallet data for conAddr %s\n", conAddr)

	// Create reverse mappings for quick lookup
	if err := redisClient.Set(ctx, "reqpubkey:"+reqPubKey, conAddr, 0).Err(); err != nil {
		// Log error if setting reqPubKey mapping fails
		fmt.Printf("Error setting reqPubKey mapping for reqPubKey %s: %v\n", reqPubKey, err)
		return err
	}
	fmt.Printf("Successfully set reqPubKey mapping for reqPubKey %s\n", reqPubKey)

	if err := redisClient.Set(ctx, "authpubkey:"+authPubKey, conAddr, 0).Err(); err != nil {
		// Log error if setting authPubKey mapping fails
		fmt.Printf("Error setting authPubKey mapping for authPubKey %s: %v\n", authPubKey, err)
		return err
	}
	fmt.Printf("Successfully set authPubKey mapping for authPubKey %s\n", authPubKey)

	return nil
}

// DeleteWallet removes the wallet data and associated reverse mappings.
func DeleteWallet(conAddr string) error {
	// Retrieve the wallet data to delete reverse mappings
	walletKey := "wallet:" + conAddr
	reqPubKey, _ := redisClient.HGet(ctx, walletKey, "reqpubkey").Result()
	authPubKey, _ := redisClient.HGet(ctx, walletKey, "authpubkey").Result()

	// Log retrieved keys for debugging
	fmt.Printf("Deleting wallet for conAddr: %s\n", conAddr)
	fmt.Printf("Retrieved reqPubKey: %s, authPubKey: %s for conAddr: %s\n", reqPubKey, authPubKey, conAddr)

	// Delete the wallet hash
	if err := redisClient.Del(ctx, walletKey).Err(); err != nil {
		fmt.Printf("Error deleting wallet hash for conAddr %s: %v\n", conAddr, err)
		return err
	}
	fmt.Printf("Successfully deleted wallet hash for conAddr: %s\n", conAddr)

	// Delete reverse mappings
	if reqPubKey != "" {
		if err := redisClient.Del(ctx, "reqpubkey:"+reqPubKey).Err(); err != nil {
			fmt.Printf("Error deleting reqPubKey mapping for reqPubKey %s: %v\n", reqPubKey, err)
			return err
		}
		fmt.Printf("Successfully deleted reqPubKey mapping for reqPubKey: %s\n", reqPubKey)
	}
	if authPubKey != "" {
		if err := redisClient.Del(ctx, "authpubkey:"+authPubKey).Err(); err != nil {
			fmt.Printf("Error deleting authPubKey mapping for authPubKey %s: %v\n", authPubKey, err)
			return err
		}
		fmt.Printf("Successfully deleted authPubKey mapping for authPubKey: %s\n", authPubKey)
	}

	return nil
}

// GetConAddrByReqPubKey retrieves the conAddr associated with the given reqPubKey.
func GetConAddrByReqPubKey(reqPubKey string) (string, error) {
	return redisClient.Get(ctx, "reqpubkey:"+reqPubKey).Result()
}

// GetConAddrByAuthPubKey retrieves the conAddr associated with the given authPubKey.
func GetConAddrByAuthPubKey(authPubKey string) (string, error) {
	return redisClient.Get(ctx, "authpubkey:"+authPubKey).Result()
}
