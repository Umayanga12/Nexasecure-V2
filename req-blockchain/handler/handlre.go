package handler

import (
	"blockchain-api/logger"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"
)

// Block structure
type Block struct {
	Index        int
	Timestamp    string
	Transactions []Transaction
	Proof        int
	PreviousHash string
}

// Transaction structure for NFT operations
type Transaction struct {
	Sender    string
	Recipient string
	NFTId     string
	Signature string
	Type      string `json:"type,omitempty"` // "transfer" or "burn"
}

// NFT structure
type NFT struct {
	Id    string
	Owner string
}

// Blockchain structure
type Blockchain struct {
	mutex               sync.Mutex
	NFTs                map[string]string // NFT ID to Owner
	CurrentTransactions []Transaction
	HostWallet          string // Address of the host wallet
	Wallets             map[string]*ecdsa.PublicKey
	Chain               []Block
}

var blockchain *Blockchain

func init() {
	blockchain = NewBlockchain()
	blockchain.CreateGenesisBlock()
}

// Initialize new blockchain
func NewBlockchain() *Blockchain {
	return &Blockchain{
		Chain:               []Block{},
		CurrentTransactions: []Transaction{},
		NFTs:                make(map[string]string),
		Wallets:             make(map[string]*ecdsa.PublicKey),
	}
}

// Create genesis block
func (bc *Blockchain) CreateGenesisBlock() {
	bc.mutex.Lock()
	defer bc.mutex.Unlock()

	genesisBlock := Block{
		Index:        1,
		Timestamp:    time.Now().String(),
		Transactions: []Transaction{},
		Proof:        100,
		PreviousHash: "1",
	}

	bc.Chain = append(bc.Chain, genesisBlock)
}

// Create new NFT
func (bc *Blockchain) CreateNFT(owner string, nftId string) error {
	config := logger.NewConfigFromEnv()
	log, err := logger.NewLogger(config)
	if err != nil {
		fmt.Printf("Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer log.Sync()

	log.Info("Starting CreateNFT", "owner", owner, "nft_id", nftId)

	bc.mutex.Lock()
	defer bc.mutex.Unlock()

	if _, exists := bc.NFTs[nftId]; exists {
		log.Error("NFT already exists", "nft_id", nftId)
		return errors.New("NFT already exists")
	}

	bc.NFTs[nftId] = owner
	log.Info("NFT created successfully", "nft_id", nftId, "owner", owner)

	bc.CurrentTransactions = append(bc.CurrentTransactions, Transaction{
		Recipient: owner,
		NFTId:     nftId,
	})
	log.Info("Transaction added for NFT creation", "nft_id", nftId, "owner", owner)

	return nil
}

// Transfer NFT to host wallet
func (bc *Blockchain) TransferNFT(sender string, nftId string, signature string) error {
	config := logger.NewConfigFromEnv()
	log, err := logger.NewLogger(config)
	if err != nil {
		fmt.Printf("Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer log.Sync()

	log.Info("Starting TransferNFT", "sender", sender, "nft_id", nftId)

	bc.mutex.Lock()
	defer bc.mutex.Unlock()

	currentOwner, exists := bc.NFTs[nftId]
	if !exists {
		log.Error("NFT does not exist", "nft_id", nftId)
		return errors.New("NFT does not exist")
	}

	if currentOwner != sender {
		log.Error("Sender not authorized", "sender", sender, "current_owner", currentOwner)
		return errors.New("sender not authorized")
	}

	if !ValidateSignature(sender, bc.HostWallet+nftId, signature) {
		log.Error("Invalid signature", "sender", sender, "nft_id", nftId)
		return errors.New("invalid signature")
	}

	bc.NFTs[nftId] = bc.HostWallet
	log.Info("NFT transferred successfully to host wallet", "nft_id", nftId, "new_owner", bc.HostWallet)

	bc.CurrentTransactions = append(bc.CurrentTransactions, Transaction{
		Sender:    sender,
		Recipient: bc.HostWallet,
		NFTId:     nftId,
		Signature: signature,
	})
	log.Info("Transaction added for NFT transfer to host wallet", "nft_id", nftId, "sender", sender)

	return nil
}

// Proof of Work algorithm
func (bc *Blockchain) ProofOfWork(lastProof int) int {
	config := logger.NewConfigFromEnv()
	log, err := logger.NewLogger(config)
	if err != nil {
		fmt.Printf("Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer log.Sync()

	log.Info("Starting ProofOfWork", "last_proof", lastProof)

	proof := 0
	for !bc.ValidProof(lastProof, proof) {
		proof++
	}

	log.Info("ProofOfWork completed", "proof", proof)
	return proof
}

// Validate proof
func (bc *Blockchain) ValidProof(lastProof, proof int) bool {
	config := logger.NewConfigFromEnv()
	log, err := logger.NewLogger(config)
	if err != nil {
		fmt.Printf("Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer log.Sync()

	guess := fmt.Sprintf("%d%d", lastProof, proof)
	hash := sha256.Sum256([]byte(guess))
	isValid := hex.EncodeToString(hash[:])[:4] == "0000"

	log.Info("ValidProof check", "last_proof", lastProof, "proof", proof, "is_valid", isValid)
	return isValid
}

// Mine new block
func (bc *Blockchain) MineBlock() Block {
	config := logger.NewConfigFromEnv()
	log, err := logger.NewLogger(config)
	if err != nil {
		fmt.Printf("Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer log.Sync()

	log.Info("Starting MineBlock")

	bc.mutex.Lock()
	defer bc.mutex.Unlock()

	lastBlock := bc.Chain[len(bc.Chain)-1]
	log.Info("Last block retrieved", "last_block_index", lastBlock.Index)

	proof := bc.ProofOfWork(lastBlock.Proof)
	log.Info("Proof of work completed", "proof", proof)

	previousHash := Hash(lastBlock)
	log.Info("Previous hash calculated", "previous_hash", previousHash)

	block := Block{
		Index:        len(bc.Chain) + 1,
		Timestamp:    time.Now().String(),
		Transactions: bc.CurrentTransactions,
		Proof:        proof,
		PreviousHash: previousHash,
	}

	bc.CurrentTransactions = nil
	bc.Chain = append(bc.Chain, block)

	log.Info("New block mined", "block_index", block.Index)
	return block
}

// Generate wallet
func GenerateWallet() (*ecdsa.PrivateKey, string, error) {
	config := logger.NewConfigFromEnv()
	log, err := logger.NewLogger(config)
	if err != nil {
		fmt.Printf("Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer log.Sync()

	log.Info("Starting GenerateWallet")

	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		log.Error("Failed to generate private key", err)
		return nil, "", err
	}

	publicKeyBytes, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		log.Error("Failed to marshal public key", err)
		return nil, "", err
	}

	address := hex.EncodeToString(publicKeyBytes)
	log.Info("Wallet generated", "address", address)
	return privateKey, address, nil
}

// Validate transaction signature
func ValidateSignature(address string, data string, signature string) bool {
	config := logger.NewConfigFromEnv()
	log, err := logger.NewLogger(config)
	if err != nil {
		fmt.Printf("Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer log.Sync()

	log.Info("Starting ValidateSignature", "address", address, "data", data)

	publicKey := blockchain.Wallets[address]
	if publicKey == nil {
		log.Error("Public key not found for address", "address", address)
		return false
	}

	sig, err := hex.DecodeString(signature)
	if err != nil {
		log.Error("Failed to decode signature", err)
		return false
	}

	hash := sha256.Sum256([]byte(data))
	isValid := ecdsa.VerifyASN1(publicKey, hash[:], sig)

	log.Info("Signature validation result", "is_valid", isValid)
	return isValid
}

// Hash block
func Hash(block Block) string {
	config := logger.NewConfigFromEnv()
	log, err := logger.NewLogger(config)
	if err != nil {
		fmt.Printf("Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer log.Sync()

	log.Info("Starting Hash", "block_index", block.Index)

	blockBytes, _ := json.Marshal(block)
	hash := sha256.Sum256(blockBytes)
	hashString := hex.EncodeToString(hash[:])

	log.Info("Block hashed", "hash", hashString)
	return hashString
}

func (bc *Blockchain) BurnNFT(sender string, nftId string, signature string) error {
	bc.mutex.Lock()
	defer bc.mutex.Unlock()

	config := logger.NewConfigFromEnv()
	log, err := logger.NewLogger(config)
	if err != nil {
		fmt.Printf("Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer log.Sync()

	currentOwner, exists := bc.NFTs[nftId]
	if !exists {
		log.Error("NFT does not exist")
		return errors.New("NFT does not exist")
	}

	if currentOwner != sender {
		log.Error("Sender not authorized")
		return errors.New("sender not authorized")
	}

	if !ValidateSignature(sender, nftId+"burn", signature) {
		log.Error("Invalid signature")
		return errors.New("invalid signature")
	}

	bc.NFTs[nftId] = bc.HostWallet
	log.Info("NFT ownership transferred to host wallet", "nft_id", nftId, "new_owner", bc.HostWallet)

	bc.CurrentTransactions = append(bc.CurrentTransactions, Transaction{
		Sender:    sender,
		Recipient: bc.HostWallet,
		NFTId:     nftId,
		Signature: signature,
	})
	return nil
}
func TransferNFTHandler(w http.ResponseWriter, r *http.Request) {
	type Request struct {
		Sender       string `json:"sender"`
		NFTId        string `json:"nft_id"`
		SignedNFTToken string `json:"signed_nfttoken"`
	}
	config := logger.NewConfigFromEnv()
	log, err := logger.NewLogger(config)
	if err != nil {
		fmt.Printf("Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer log.Sync()

	var req Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Error("Failed to decode request body", err)
		return
	}

	recipient := blockchain.HostWallet // All NFTs are transferred to the host wallet

	// Validate the signed NFT token
	if !ValidateSignature(req.Sender, blockchain.HostWallet+req.NFTId, req.SignedNFTToken) {
		http.Error(w, "Invalid signed NFT token", http.StatusBadRequest)
		log.Error("Invalid signed NFT token", "sender", req.Sender, "nft_id", req.NFTId)
		return
	}

	if err := blockchain.TransferNFT(req.Sender, req.NFTId, req.SignedNFTToken); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Error("Failed to transfer NFT", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	log.Info("NFT transferred to host wallet", "nft_id", req.NFTId, "host_wallet", recipient)
	//burn the NFT
}

func GetNFTOwnerHandler(w http.ResponseWriter, r *http.Request) {
	type Request struct {
		NFTId string `json:"nft_id"`
	}

	var req Request
	config := logger.NewConfigFromEnv()
	log, err := logger.NewLogger(config)
	if err != nil {
		fmt.Printf("Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer log.Sync()

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Error("Failed to decode request body", err)
		return
	}

	owner, exists := blockchain.NFTs[req.NFTId]
	if !exists {
		http.Error(w, "NFT not found", http.StatusNotFound)
		log.Error("NFT not found", "nft_id", req.NFTId)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"owner": owner})
}

// Burn NFT handler
func BurnNFTHandler(w http.ResponseWriter, r *http.Request) {
	type Request struct {
		Sender    string `json:"sender"`
		NFTId     string `json:"nft_id"`
		Signature string `json:"signature"`
	}

	config := logger.NewConfigFromEnv()
	log, err := logger.NewLogger(config)
	if err != nil {
		fmt.Printf("Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer log.Sync()

	var req Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Error("Failed to decode request body", err)
		return
	}

	if err := blockchain.BurnNFT(req.Sender, req.NFTId, req.Signature); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Error("Failed to burn NFT", err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func GenerateWalletHandler(w http.ResponseWriter, r *http.Request) {
	config := logger.NewConfigFromEnv()
	log, err := logger.NewLogger(config)
	if err != nil {
		fmt.Printf("Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	privateKey, address, err := GenerateWallet()
	if err != nil {
		log.Error("Failed to generate wallet", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	blockchain.mutex.Lock()
	blockchain.Wallets[address] = &privateKey.PublicKey
	blockchain.mutex.Unlock()

	privateKeyBytes, err := x509.MarshalECPrivateKey(privateKey)
	if err != nil {
		log.Error("Failed to marshal private key", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"address":    address,
		"privateKey": hex.EncodeToString(privateKeyBytes),
	})
}

// HTTP Handlers
func CreateNFTHandler(w http.ResponseWriter, r *http.Request) {

	type Request struct {
		Owner string `json:"owner"`
	}

	config := logger.NewConfigFromEnv()
	log, err := logger.NewLogger(config)
	if err != nil {
		fmt.Printf("Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer log.Sync()

	var req Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Error("Failed to decode request body", err)
		return
	}

	// Generate a unique NFT ID
	nftId := fmt.Sprintf("nft-%d", time.Now().UnixNano())

	if err := blockchain.CreateNFT(req.Owner, nftId); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Error("Failed to create NFT", err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	log.Info("NFT created")
	json.NewEncoder(w).Encode(map[string]string{"nft_id": nftId})
}

// Add validation handler
func ValidateNFTOwnerHandler(w http.ResponseWriter, r *http.Request) {

	config := logger.NewConfigFromEnv()
	log, err := logger.NewLogger(config)
	if err != nil {
		fmt.Printf("Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer log.Sync()

	log.Info("Starting ValidateNFTOwnerHandler")

	type Request struct {
		NFTId   string `json:"nft_id"`
		Address string `json:"address"`
	}

	var req Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Error("Failed to decode request body", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.NFTId == "" || req.Address == "" {
		log.Error("Missing nft_id or address in request body")
		http.Error(w, "Missing nft_id or address in request body", http.StatusBadRequest)
		return
	}

	log.Info("Received request", "nft_id", req.NFTId, "address", req.Address)

	blockchain.mutex.Lock()
	owner, exists := blockchain.NFTs[req.NFTId]
	blockchain.mutex.Unlock()

	if !exists {
		log.Error("NFT not found", "nft_id", req.NFTId)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"valid": false,
			"error": "NFT not found",
		})
		return
	}

	log.Info("NFT found", "nft_id", req.NFTId, "owner", owner)

	// Check if the provided address matches the owner
	isValid := owner == req.Address
	if !isValid {
		log.Error("Address does not match the owner", "nft_id", req.NFTId, "address", req.Address)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"valid": false,
			"error": "Address does not match the owner",
		})
		return
	}

	log.Info("Ownership validated successfully", "nft_id", req.NFTId, "address", req.Address)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"valid": true,
	})
}

func RootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(w, "Welcome to the Blockchain API!")
}
