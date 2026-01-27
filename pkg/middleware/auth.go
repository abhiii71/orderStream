package middleware

import (
	"context"
	"strings"

	"github.com/abhiii71/orderStream/pkg/auth"
	"github.com/abhiii71/orderStream/pkg/contextkeys"
	"github.com/gin-gonic/gin"
)

func AuthorizeJWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		var tokenString string

		// First try to get token from Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
			tokenString = strings.TrimPrefix(authHeader, "Bearer ")
		}

		// If not in header, try cookie
		if tokenString == "" {
			authCookie, err := c.Cookie("token")
			if err == nil && authCookie != "" {
				tokenString = authCookie
			}
		}

		if tokenString == "" {
			c.Set("userID", "")
			c.Next()
			return
		}

		token, err := auth.ValidateToken(tokenString)
		if err != nil {
			c.Set("userID", "")
			c.Next()
			return
		}

		if claims, ok := token.Claims.(*auth.JWTCustomClaims); ok && token.Valid {
			c.Set("userID", claims.UserID)
			ctxWithVal := context.WithValue(c.Request.Context(), contextkeys.UserIDKey, claims.UserID)
			c.Request = c.Request.WithContext(ctxWithVal)
		} else {
			c.Set("userID", "")
		}
		c.Next()
	}
}
