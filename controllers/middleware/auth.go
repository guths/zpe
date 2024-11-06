package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/guths/zpe/config"
	"github.com/guths/zpe/constants"
	"github.com/guths/zpe/datatransfers"
)

//simple jwt logic, receive a token by a bearer, parse and verify

func AuthMiddleware(c *gin.Context) {
	token := strings.TrimPrefix(c.GetHeader("Authorization"), "Bearer ")
	if token == "" {
		c.Set(constants.IsAuthenticatedKey, false)
		c.Next()
		return
	}
	claims, err := parseToken(token, config.AppConfig.JWTSecret)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, datatransfers.Response{Error: err.Error()})
		return
	}

	c.Set(constants.IsAuthenticatedKey, true)
	c.Set(constants.UserIDKey, claims.ID)
	c.Set(constants.UserRoles, claims.Roles)

	c.Next()
}

func parseToken(tokenString, secret string) (*datatransfers.JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &datatransfers.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil || !token.Valid {
		if err == jwt.ErrTokenExpired {
			return nil, fmt.Errorf("token expired")
		}

		return nil, fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(*datatransfers.JWTClaims)

	if !ok || !token.Valid {
		return nil, jwt.ErrTokenInvalidClaims
	}

	return claims, nil

}
