package handler

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"auth-service/internal/model"
	"auth-service/internal/repository"
	"auth-service/internal/utils"
)

type AuthHandler struct {
	repo repository.UserRepository
}

func NewAuthHandler(repo repository.UserRepository) *AuthHandler {
	return &AuthHandler{repo: repo}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req model.User
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	hashed, _ := utils.HashPassword(req.Password)
	req.Password = hashed

	if err := h.repo.CreateUser(&req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Registered"})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req model.User
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	dbUser, err := h.repo.FindByEmail(req.Email)
	if err != nil || !utils.CheckPasswordHash(req.Password, dbUser.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}
	accessToken, refreshToken, _ := utils.GenerateTokens(int(dbUser.ID))
	c.JSON(http.StatusOK, gin.H{
		"token":         accessToken,
		"refresh_token": refreshToken,
	})
}

func (h *AuthHandler) RefreshToken(c *gin.Context) {
	claims := c.MustGet("claims").(map[string]interface{})
	userId := int(claims["user_id"].(float64))
	user, err := h.repo.FindByID(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User not found"})
		return
	}
	newAccessToken, newRefreshToken, _ := utils.GenerateTokens(int(user.ID))
	c.JSON(http.StatusOK, gin.H{
		"token":         newAccessToken,
		"refresh_token": newRefreshToken,
	})
}

func (h *AuthHandler) Me(c *gin.Context) {
	claims := c.MustGet("claims").(map[string]interface{})
	c.JSON(http.StatusOK, gin.H{"claims": claims})
}

func (h *AuthHandler) VerifyToken(c *gin.Context) {
	claims := c.MustGet("claims").(map[string]interface{})
	c.JSON(http.StatusOK, gin.H{"valid": true, "claims": claims})
}
