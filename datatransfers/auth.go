package datatransfers

import (
	"github.com/golang-jwt/jwt/v5"
)

type JWTClaims struct {
	ID        uint     `json:"sub,omitempty"`
	Roles     []string `json:"roles"`
	ExpiresAt int64    `json:"exp,omitempty"`
	IssuedAt  int64    `json:"iat,omitempty"`
	jwt.RegisteredClaims
}
