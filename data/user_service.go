package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"task_manager/data"
)

		auth := c.GetHeader("Authorization")
		if !strings.HasPrefix(auth, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing bearer token"})
			return
		}
		token := strings.TrimSpace(strings.TrimPrefix(auth, "Bearer "))

		claims, err := data.ParseToken(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		
		c.Set("claims", claims)
		if v, ok := claims["email"].(string); ok {
			c.Set("email", v)
		}
		if v, ok := claims["role"].(string); ok {
			c.Set("role", v)
		}
		if v, ok := claims["user_id"]; ok {
			c.Set("user_id", v)
		}

		c.Next()
	}
}

func AdminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		roleVal, exists := c.Get("role")
		role, ok := roleVal.(string)
		if !exists || !ok || role != "admin" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "admin only"})
			return
		}
		c.Next()
	}
}

