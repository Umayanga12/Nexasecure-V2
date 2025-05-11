package main

import (
	"fmt"
	"net"
	"net/http"
	"server/handler"
)



func main() {
	fmt.Println("Socket Server starting...")

	// Start the TCP server
	go startTCPServer()

	// Start the HTTP server
	// http.HandleFunc("/send", handler.SendMessageHandler)
	http.HandleFunc("/setauthnft", handler.StoreAuthNFTHandler)
	http.HandleFunc("/getauthnft", handler.GetAuthNFTHandler)
	http.HandleFunc("/removeauthnft", handler.RemoveAuthNFTHandler)
	http.HandleFunc("/setreqnft", handler.StoreReqNFTHandler)
	http.HandleFunc("/getreqnft", handler.GetReqNFTHandler)
	http.HandleFunc("/removereqnft", handler.RemoveReqNFTHandler)
	http.HandleFunc("/signauthwallet", handler.SignAuthwalletHandler)
	http.HandleFunc("/signreqwallet", handler.SignReqwalletHandler)
	// http.HandleFunc("/login", handler.LoginHandle)
	// http.HandleFunc("/signmsg", handler.SignMsgHandler)
	fmt.Println("HTTP API listening on port 10081")
	if err := http.ListenAndServe(":10081", nil); err != nil {
		fmt.Println("Error starting HTTP server:", err)
	}
}

func startTCPServer() {
	ln, err := net.Listen("tcp", ":10080")
	if err != nil {
		fmt.Println("Error listening:", err)
		return
	}
	defer ln.Close()
	fmt.Println("TCP server listening on port 10080")

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		go handler.HandleTCPConnection(conn)
	}
}



