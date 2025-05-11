package handler

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"server/util"
	"strings"
	"sync"
	"time"
)

var (
	clients         = make(map[string]net.Conn)
	clientsMu       sync.Mutex
	responseChannels = make(map[string]chan string)
	responsesMu     sync.Mutex
)

func HandleTCPConnection(conn net.Conn) {
	defer func() {
		clientsMu.Lock()
		delete(clients, conn.RemoteAddr().String())
		clientsMu.Unlock()
		conn.Close()
	}()

	// Authentication
	scanner := bufio.NewScanner(conn)
	if !scanner.Scan() {
		return
	}

	authData := strings.Split(scanner.Text(), ",")
	if len(authData) != 2 {
		conn.Write([]byte("Invalid authentication format\n"))
		return
	}

	requestPubAddr := strings.TrimSpace(authData[0])
	authPubAddr := strings.TrimSpace(authData[1])

	if !util.ValidateAuthaddr(authPubAddr) || !util.ValidateReqaddr(requestPubAddr) {
		conn.Write([]byte("Validation failed\n"))
		return
	}

	conn.Write([]byte("VALIDATED\n"))
	util.AddWallet(conn.RemoteAddr().String(), requestPubAddr, authPubAddr)

	clientsMu.Lock()
	clients[conn.RemoteAddr().String()] = conn
	clientsMu.Unlock()

	// Setup client-specific response channel
	clientAddr := conn.RemoteAddr().String()
	responsesMu.Lock()
	responseChannels[clientAddr] = make(chan string, 1)
	responsesMu.Unlock()

	defer func() {
		responsesMu.Lock()
		delete(responseChannels, clientAddr)
		responsesMu.Unlock()
	}()

	// Read messages
	reader := bufio.NewReader(conn)
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				util.DeleteWallet(clientAddr)
				fmt.Printf("Client %s disconnected\n", clientAddr)
			}
			return
		}

		// Forward raw message to response channel
		responsesMu.Lock()
		if ch, exists := responseChannels[clientAddr]; exists {
			ch <- strings.TrimSpace(message)
		}
		responsesMu.Unlock()
	}
}

func SendMessageToClient(clientAddr string, message string) (string, error) {
	clientsMu.Lock()
	conn, exists := clients[clientAddr]
	clientsMu.Unlock()

	if !exists {
		return "", fmt.Errorf("client %s not found", clientAddr)
	}

	// Get client-specific response channel
	responsesMu.Lock()
	responseChan, responseChanExists := responseChannels[clientAddr]
	responsesMu.Unlock()

	if !responseChanExists {
		return "", fmt.Errorf("response channel for client %s not found", clientAddr)
	}

	// Send message
	_, err := conn.Write([]byte(message + "\n"))
	if err != nil {
		return "", fmt.Errorf("failed to send message: %v", err)
	}

	// Wait for response
	select {
	case response := <-responseChan:
		if strings.TrimSpace(response) == "" {
			return "No response from client", nil
		}
		return response, nil
	case <-time.After(10 * time.Second):
		return "No response from client", nil
	}
}