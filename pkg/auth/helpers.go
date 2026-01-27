package auth

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/abhiii71/orderStream/pkg/contextkeys"
	"github.com/gin-gonic/gin"
)

// GinContextKey matches the key used in middleware package
const GinContextKey = "GinContextKey"

func GetUserId(ctx context.Context, abort bool) string {
	userId, err := GetUserIdInt(ctx, abort)
	if err != nil {
		return ""
	}
	return strconv.Itoa(userId)
}

func GetUserIdInt(ctx context.Context, abort bool) (int, error) {
	accountId, ok := ctx.Value(contextkeys.UserIDKey).(uint64)
	if !ok {
		if abort {
			ginContext, ginOk := ctx.Value(GinContextKey).(*gin.Context)
			if ginOk && ginContext != nil {
				ginContext.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			}
		}
		return 0, errors.New("UserId not found in context")
	}

	return int(accountId), nil
}
