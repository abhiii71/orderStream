package middleware

import (
	"context"

	"github.com/gin-gonic/gin"
)

// Key to use when setting the gin context - using string key for compatibility
const GinContextKey = "GinContextKey"

func GinContextToContextMiddlware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Put gin.Context into the request context so gqlgen can retrieve it
		ctx := context.WithValue(c.Request.Context(), GinContextKey, c)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}
