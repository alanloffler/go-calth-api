package middleware

import (
	"net/http"

	"github.com/alanloffler/go-calth-api/internal/auth"
	"github.com/alanloffler/go-calth-api/internal/common/response"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware(service *auth.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr, err := c.Cookie("access_token")
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, response.Error(http.StatusUnauthorized, "Token requerido"))
			return
		}

		claims, err := service.ValidateAccessToken(tokenStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, response.Error(http.StatusUnauthorized, "Token inválido o expirado", err))
			return
		}

		c.Set("userID", claims.UserID)
		c.Set("businessID", claims.BusinessID)
		c.Set("roleID", claims.RoleID)

		c.Next()
	}
}
