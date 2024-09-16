package middlewares

import (
	"net/http"

	"github.com/Atheer-Ganayem/Go-social-media-api/utils"
	"github.com/gin-gonic/gin"
)

func IsAuth(ctx *gin.Context) {
	token := ctx.Request.Header.Get("Authorization")
	if token == "" {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Not authenticated."})
		return
	}

	userId, err := utils.VerifyToken(token)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Not authenticated."})
		return
	}

	ctx.Set("userId", userId)

	ctx.Next()
}