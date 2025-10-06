package auth

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/abhiii71/orderStream/pkg/contextkeys"
	"github.com/gin-gonic/gin"
)

func GetUserId(ctx context.Context, abort bool) string {
	userId, err := GetUserIdInt(ctx, abort)
	if err != nil {
		return ""
	}
	return strconv.Itoa(userId)
}

func GetUserIdInt(ctx context.Context, abort bool) (int, error) {
	accountId, ok := ctx.Value(contextkeys.UserIdKey).(uint64)
	if !ok {
		if abort {
			ginContext, _ := ctx.Value(contextkeys.UserIdKey).(*gin.Context)
			ginContext.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		}
		return 0, errors.New("UserId not found in  context")
	}

	return int(accountId), nil
}
