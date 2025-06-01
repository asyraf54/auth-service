package middleware

import (
	"auth-service/internal/response"
	"auth-service/internal/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func JWTAuth(expectedType string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			response.AbortWithStatusJSON(c, http.StatusUnauthorized, nil, &response.ErrorDetail{
				ErrorCode:    response.ErrInvalidToken,
				ErrorMessage: "Invalid or missing token.", 
				ErrorDebugMessage: "Invalid or missing token.",
			})
			return
		}
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := utils.ValidateToken(tokenString, expectedType)
		if err != nil {
			response.AbortWithStatusJSON(c, http.StatusUnauthorized, nil, &response.ErrorDetail{
				ErrorCode:    response.ErrInvalidToken,
				ErrorMessage: err.Error(), 
				ErrorDebugMessage: err.Error(),
			})
			return
		}
		c.Set("claims", claims)
		c.Next()
	}
}
