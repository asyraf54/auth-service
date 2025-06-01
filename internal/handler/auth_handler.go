package handler

import (
	"auth-service/internal/model"
	"auth-service/internal/repository"
	"auth-service/internal/response"
	"auth-service/internal/utils"
	"net/http"
	"github.com/gin-gonic/gin"
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
		response.JSON(c, http.StatusBadRequest, nil, &response.ErrorDetail{
			ErrorCode:         response.ErrInvalidJSON,
			ErrorMessage:      "Format data tidak valid.",
			ErrorDebugMessage: err.Error(),
		})
		return
	}
	hashed, err := utils.HashPassword(req.Password)
	if err != nil {
		response.JSON(c, http.StatusInternalServerError, nil, &response.ErrorDetail{
			ErrorCode:         response.ErrHash,
			ErrorMessage:      "Terjadi kesalahan pada server.",
			ErrorDebugMessage: err.Error(),
		})
		return
	}
	req.Password = hashed

	if err := h.repo.CreateUser(&req); err != nil {
		response.JSON(c, http.StatusInternalServerError, nil, &response.ErrorDetail{
			ErrorCode:         response.ErrCreateUserError,
			ErrorMessage:      "Terjadi kesalahan pada server.",
			ErrorDebugMessage: err.Error(),
		})
		return
	}
	response.JSON(c, http.StatusCreated, nil, nil)
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req model.User
	if err := c.ShouldBindJSON(&req); err != nil {
		response.JSON(c, http.StatusBadRequest, nil, &response.ErrorDetail{
			ErrorCode:         response.ErrInvalidJSON,
			ErrorMessage:      "Format data tidak valid.",
			ErrorDebugMessage: err.Error(),
		})
		return
	}
	dbUser, err := h.repo.FindByEmail(req.Email)
	if err != nil {
		response.JSON(c, http.StatusUnauthorized, nil, &response.ErrorDetail{
			ErrorCode:         response.ErrUserNotFound,
			ErrorMessage:      "Email atau password salah.",
			ErrorDebugMessage: err.Error(),
		})
		return
	}

	if !utils.CheckPasswordHash(req.Password, dbUser.Password) {
		response.JSON(c, http.StatusUnauthorized, nil, &response.ErrorDetail{
			ErrorCode:         response.ErrInvalidLogin,
			ErrorMessage:      "Email atau password salah.",
			ErrorDebugMessage: "Password tidak cocok",
		})
		return
	}

	accessToken, refreshToken, err := utils.GenerateTokens(int(dbUser.ID))
	if err != nil {
		response.JSON(c, http.StatusInternalServerError, nil, &response.ErrorDetail{
			ErrorCode:         response.ErrGenerateToken,
			ErrorMessage:      "Terjadi kesalahan pada server.",
			ErrorDebugMessage: err.Error(),
		})
		return
	}

	response.JSON(c, http.StatusOK, gin.H{
		"token":         accessToken,
		"refresh_token": refreshToken,
	}, nil)
}

func (h *AuthHandler) RefreshToken(c *gin.Context) {
	claims := c.MustGet("claims").(map[string]interface{})
	userId := int(claims["user_id"].(float64))
	user, err := h.repo.FindByID(userId)
	if err != nil {
		response.JSON(c, http.StatusUnauthorized, nil, &response.ErrorDetail{
			ErrorCode:         response.ErrUserNotFound,
			ErrorMessage:      "Pengguna tidak ditemukan.",
			ErrorDebugMessage: err.Error(),
		})
	}
	newAccessToken, newRefreshToken, err := utils.GenerateTokens(int(user.ID))
	if err != nil {
		response.JSON(c, http.StatusInternalServerError, nil, &response.ErrorDetail{
			ErrorCode:         response.ErrGenerateToken,
			ErrorMessage:      "Terjadi kesalahan pada server.",
			ErrorDebugMessage: err.Error(),
		})
		return
	}

	response.JSON(c, http.StatusOK, gin.H{
		"token":         newAccessToken,
		"refresh_token": newRefreshToken,
	}, nil)
}

func (h *AuthHandler) Me(c *gin.Context) {
	claims := c.MustGet("claims").(map[string]interface{})
	response.JSON(c, http.StatusOK, gin.H{"claims": claims}, nil)
}

func (h *AuthHandler) VerifyToken(c *gin.Context) {
	claims := c.MustGet("claims").(map[string]interface{})
	response.JSON(c, http.StatusOK, gin.H{"valid": true, "claims": claims}, nil)
}
