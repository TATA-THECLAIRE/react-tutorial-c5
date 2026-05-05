package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"piggy.com/internal/auth"
)

func RequireAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		header := ctx.GetHeader("Authorization")
		if !strings.HasPrefix(header, "Bearer ") {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing token"})
			return
		}
		tokenStr := strings.TrimPrefix(header, "Bearer ")
		claims, err := auth.ValidateToken(tokenStr)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}
		// Store as pgtype.UUID by converting from [16]byte
		ctx.Set("userID", claims.UserID) // [16]byte
		ctx.Set("username", claims.Username)
		ctx.Next()
	}
}