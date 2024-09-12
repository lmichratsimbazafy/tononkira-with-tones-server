package http

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"lmich.com/tononkira/domain"
)

type AuthHandler struct {
	UserService domain.UserService
	AuthService domain.AuthService
}

type AuthPayload struct {
	UserName string `json:"userName" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Login handles user login and token generation
func (h *AuthHandler) Login(c *gin.Context) {
	var credentials AuthPayload
	if err := c.ShouldBindJSON(&credentials); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": fmt.Sprintf("%v", err)})
	}

	user, err := h.UserService.GetUser(domain.UserFilter{UserName: credentials.UserName})

	if err != nil || user == nil {
		c.JSON(http.StatusNotFound, gin.H{"errors": fmt.Sprintf("%v", err)})
		return
	}
	err = checkPassword(credentials.Password, user.Password)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"errors": fmt.Sprintf("%v", err)})
		return
	}
	token, err := h.AuthService.GenerateToken(user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": fmt.Sprintf("%v", err)})
		fmt.Println("error while validating token", gin.H{"errors": fmt.Sprintf("%v", err)})
		return
	}
	c.IndentedJSON(http.StatusOK, token)
}

func checkPassword(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
