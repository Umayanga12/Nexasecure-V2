package handler

import (
	"net"
	//"net/http"
	"sync"
	//"time"

	//"github.com/google/uuid"
)

type TCPConnManager struct {
    connections map[string]net.Conn
    responses   map[string]chan string
    mu          sync.Mutex
}

func NewTCPConnManager() *TCPConnManager {
    return &TCPConnManager{
        connections: make(map[string]net.Conn),
        responses:   make(map[string]chan string),
    }
}

func (m *TCPConnManager) Register(requestID string, ch chan string) {
    m.mu.Lock()
    defer m.mu.Unlock()
    m.responses[requestID] = ch
}

func (m *TCPConnManager) NotifyResponse(requestID, response string) {
    m.mu.Lock()
    defer m.mu.Unlock()
    if ch, exists := m.responses[requestID]; exists {
        ch <- response
        delete(m.responses, requestID)
    }
}

