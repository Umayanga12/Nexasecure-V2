// internal/models/models.go
package models

import (
	"time"
)

type User struct {
	ID          string `json:"id"`
	Username    string `json:"username" validate:"required"`
	Email       string `json:"email" validate:"required,email"`
	Password    string `json:"password" validate:"required,min=8"`
	ReqPubAddr  string `json:"reqpubaddr"`
	AuthPubAddr string `json:"authpubaddr"`
}

type OTP struct {
	UserID    string    `json:"user_id"`
	Code      string    `json:"code"`
	ExpiresAt time.Time `json:"expires_at"`
}