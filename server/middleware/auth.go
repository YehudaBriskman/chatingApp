package middleware

import (
	"net/http"
	"strings"
	"chatingApp/services"
	"github.com/gin-gonic/gin"
)

// AuthMiddleware validates JWT token and adds user data to the request context
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		claims, err := services.ValidateToken(strings.TrimPrefix(token, "Bearer "))
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		c.Set("userID", claims["id"])
		c.Set("role", claims["role"])
		c.Next()
	}
}

// RoleHierarchy defines the order of roles
var RoleHierarchy = map[string]int{
	"user":        1,
	"admin":       2,
	"super-admin": 3,
}

// HasRequiredRole checks if the user has the required role or higher
func HasRequiredRole(userRole, requiredRole string) bool {
	userRank, userExists := RoleHierarchy[userRole]
	requiredRank, requiredExists := RoleHierarchy[requiredRole]

	if !userExists || !requiredExists {
		return false 
	}

	return userRank >= requiredRank
}

// AdminMiddleware ensures only users with the required role can access the route
func AdminMiddleware(requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists || !HasRequiredRole(role.(string), requiredRole) {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
			c.Abort()
			return
		}
		c.Next()
	}
}
