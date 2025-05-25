package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"auth-service/internal/repository"
	"auth-service/internal/utils"
	"auth-service/internal/model"
)

type AuthHandler struct {
	repo repository.UserRepository
}

func Register(repo repository.UserRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req model.User
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		hashed, _ := utils.HashPassword(req.Password)
		req.Password = hashed
		if err := repo.CreateUser(&req); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, gin.H{"message": "Registered"})
	}
}

func Login(repo repository.UserRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req model.User
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		dbUser, err := repo.FindByEmail(req.Email)
		if err != nil || !utils.CheckPasswordHash(req.Password, dbUser.Password) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			return
		}
		accessToken, refreshToken, _ := utils.GenerateTokens(int(dbUser.ID))
		c.JSON(http.StatusOK, gin.H{"token": accessToken, "refresh_token": refreshToken})
	}
}

func Me(c *gin.Context) {
	claims := c.MustGet("claims").(map[string]interface{})
	c.JSON(http.StatusOK, gin.H{"claims": claims})
}

func VerifyToken(c *gin.Context) {
	claims := c.MustGet("claims").(map[string]interface{})
	c.JSON(http.StatusOK, gin.H{"valid": true, "claims": claims})
}
