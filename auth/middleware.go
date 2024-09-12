package auth

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"lmich.com/tononkira/domain"
)

type AuthMiddleware struct {
	AuthService domain.AuthService
}

// NewAuthMiddleware creates a new AuthMiddleware
func NewAuthMiddleware(authService domain.AuthService) *AuthMiddleware {
	return &AuthMiddleware{AuthService: authService}
}

// Middleware function to enforce authentication
func (am *AuthMiddleware) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Récupérer le header Authorization
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		// Extraire le token du header
		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == authHeader {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Bearer token required"})
			c.Abort()
			return
		}

		// Valider le token via le service d'authentification
		userId, err := am.AuthService.ValidateToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		// Stocker l'ID de l'utilisateur dans le contexte pour un usage ultérieur
		c.Set("userId", userId)
		c.Next() // Passer à la route suivante
	}
}
