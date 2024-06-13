package models

import (
	"github.com/golang-jwt/jwt/v4"
)

type ResetPasswordParams struct {
	Token       string `json:"token"`
	NewPassword string `json:"new_password"`
}

type ResetPasswordResponse struct {
	Message string `json:"message"`
}

type LoginParams struct {
	Email    string
	Password string
}

type LoginResponse struct {
	AccessToken string `json:"access_token"`
	User        User   `json:"user"`
	Message     string `json:"message"`
}

type Claims struct {
	Sub   int    `json:"sub"`
	Fresh bool   `json:"fresh"`
	Type  string `json:"type"`
	jwt.RegisteredClaims
}