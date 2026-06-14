package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ob-sync-server/internal/repository"
	"github.com/ob-sync-server/internal/util"
)

type UserHandler struct {
	userRepo *repository.UserRepository
	logger   *util.Logger
}

func NewUserHandler(userRepo *repository.UserRepository, logger *util.Logger) *UserHandler {
	return &UserHandler{userRepo: userRepo, logger: logger}
}

func (h *UserHandler) GenerateUserID(c *gin.Context) {
	user, err := h.userRepo.CreateUser()
	if err != nil {
		h.logger.Error("Failed to generate user ID:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate user ID"})
		return
	}

	h.logger.Info("Generated new user ID:", user.UserID)
	c.JSON(http.StatusOK, gin.H{"user_id": user.UserID})
}

func (h *UserHandler) ValidateUserID(c *gin.Context) {
	var request struct {
		UserID string `json:"user_id"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		h.logger.Error("Invalid request body:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	exists, err := h.userRepo.UserExists(request.UserID)
	if err != nil {
		h.logger.Error("Failed to validate user ID:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"valid": exists})
}
