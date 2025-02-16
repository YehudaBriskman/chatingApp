package middleware

import (
	"chatingApp/services"
	"errors"
	"net/http"
	"strings"

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

		// שמירת הנתונים ב-Context לשימוש מאוחר יותר
		c.Set("userID", extractUserID(claims))
		c.Set("role", claims["role"])
		c.Set("email", claims["email"])
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

// ExtractTokenData validates the token, checks the role, and returns user info
func ExtractTokenData(c *gin.Context, requiredRole string) (int, string, string, error) {
	token := c.GetHeader("Authorization")

	if token == "" {
		return 0, "", "", errors.New("missing authentication token")
	}

	claims, err := services.ValidateToken(strings.TrimPrefix(token, "Bearer "))
	if err != nil {
		return 0, "", "", errors.New("invalid or expired token")
	}

	role, ok := claims["role"].(string)
	if !ok {
		return 0, "", "", errors.New("invalid token data: missing role")
	}

	if requiredRole != "" && !HasRequiredRole(role, requiredRole) {
		return 0, "", "", errors.New("access denied: insufficient role permissions")
	}

	userID := extractUserID(claims)
	if userID == 0 {
		return 0, "", "", errors.New("invalid token data: user ID missing or invalid")
	}

	email, _ := claims["email"].(string)

	return userID, role, email, nil
}

// RequireTokenAndRole is a middleware that ensures token authentication and a specific role
func RequireTokenAndRole(requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		_, _, _, err := ExtractTokenData(c, requiredRole)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return
		}
		c.Next()
	}
}

// extractUserID extracts user ID from JWT claims
func extractUserID(claims map[string]interface{}) int {
	userIDRaw, exists := claims["user_id"]
	if !exists {
		return 0
	}

	switch v := userIDRaw.(type) {
	case float64:
		return int(v)
	case int:
		return v
	default:
		return 0
	}
}
