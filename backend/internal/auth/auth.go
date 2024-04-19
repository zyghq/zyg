package auth

import "github.com/golang-jwt/jwt/v5"

// AuthJWTClaims taken from Supabase JWT encoding
type AuthJWTClaims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}

type CustomerJWTClaims struct {
	WorkspaceId string `json:"workspaceId"`
	ExternalId  string `json:"externalId"`
	Email       string `json:"email"`
	Phone       string `json:"phone"`
	jwt.RegisteredClaims
}
