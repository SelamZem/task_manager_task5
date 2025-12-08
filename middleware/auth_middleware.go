package middleware

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

// getSecret reads JWT secret from env or returns a default (change in prod)
func getSecret() []byte {
	if s := os.Getenv("JWT_SECRET"); s != "" {
		return []byte(s)
	}
	return []byte("replace_this_secret")
}

// AuthMiddleware validates Bearer JWT and exposes claims in context:
// - "claims" => jwt.MapClaims
// - "email", "role", "user_id" if present
func AuthMiddleware() gin.HandlerFunc {
	secret := getSecret()
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if !strings.HasPrefix(auth, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing bearer token"})
			return
		}
		tokenStr := strings.TrimSpace(strings.TrimPrefix(auth, "Bearer "))

		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrTokenUnverifiable
			}
			return secret, nil
		})
		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token claims"})
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

// AdminOnly aborts with 403 if role != "admin"
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
