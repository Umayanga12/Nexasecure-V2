// internal/utils/mailer.go
package utils

import (
	"auth-server/logger"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type Mail struct {
	From    EmailUser   `json:"from"`
	To      []EmailUser `json:"to"`
	Subject string      `json:"subject"`
	Text    string      `json:"text"`
}

type EmailUser struct {
	Email string `json:"email"`
	Name  string `json:"name,omitempty"`
}

func SendOTPEmail(email, otp string) error {
	fmt.Println("Sending OTP email...")

	config := logger.NewConfigFromEnv()

	// Initialize logger
	log, err := logger.NewLogger(config)
	if err != nil {
		fmt.Printf("Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer log.Sync()

	mail := Mail{
		From: EmailUser{
			Email: "hello@demomailtrap.co", 
			Name:  "Nexasecure 2FA",
		},
		To: []EmailUser{
			{Email: email},
		},
		Subject: "Your OTP Code",
		Text:    fmt.Sprintf("Your OTP code is: %s\nIt will expire in 5 minutes.", otp),
	}
	fmt.Println("Mail object created:", mail)
	body, err := json.Marshal(mail)
	if err != nil {
		return err
	}
	fmt.Println("Mail object marshalled to JSON:", string(body))
	req, err := http.NewRequest("POST", "https://send.api.mailtrap.io/api/send", bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	fmt.Println("HTTP request created:", req)
	//fmt.Println(os.Getenv("MAILTRAP_API_KEY"))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+os.Getenv("MAILTRAP_API_KEY")) // Make sure this is set correctly

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	fmt.Println("HTTP response received:", resp)
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("email sending failed with status: %d", resp.StatusCode)
	}

	return nil
}
