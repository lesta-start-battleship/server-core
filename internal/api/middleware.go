package api

import (
	"github.com/gin-gonic/gin"
	"strings"
	"encoding/base64"
	"encoding/json"

	// "fmt"
)

type JWTPayload struct {
	Fresh    bool   `json:"fresh"`
	Iat      int64  `json:"iat"`
	Jti      string `json:"jti"`
	Type     string `json:"type"`
	Sub      string `json:"sub"`
	Nbf      int64  `json:"nbf"`
	Csrf     string `json:"csrf"`
	Exp      int64  `json:"exp"`
	Username string `json:"username"`
	Role     string `json:"role"`
	IsActive bool   `json:"is_active"`
}

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(401, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		parts := strings.Split(tokenString, ".")
		if len(parts) != 3 {
			c.JSON(400, gin.H{"error": "Invalid JWT format"})
			c.Abort()
			return
		}

		payloadBytes, err := base64.RawURLEncoding.DecodeString(parts[1])
		if err != nil {
			c.JSON(400, gin.H{"error": "Failed to decode JWT payload"})
			c.Abort()
			return
		}

		var payload JWTPayload
		if err := json.Unmarshal(payloadBytes, &payload); err != nil {
			c.JSON(400, gin.H{"error": "Failed to parse JWT payload"})
			c.Abort()
			return
		}

		
		c.Set("player_id", payload.Sub)
		c.Set("auth_header", authHeader)
		c.Set("login", payload.Username)


		c.Next()
	}
}