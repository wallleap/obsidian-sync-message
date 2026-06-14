package handler

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/ob-sync-server/config"
	"github.com/ob-sync-server/internal/model"
	"github.com/ob-sync-server/internal/repository"
	"github.com/ob-sync-server/internal/util"
)

type MessageHandler struct {
	userRepo       *repository.UserRepository
	messageRepo    *repository.MessageRepository
	attachmentRepo *repository.AttachmentRepository
	logger         *util.Logger
	cfg            *config.Config
	snowflake      *util.Snowflake
}

func NewMessageHandler(
	userRepo *repository.UserRepository,
	messageRepo *repository.MessageRepository,
	attachmentRepo *repository.AttachmentRepository,
	logger *util.Logger,
	cfg *config.Config,
) *MessageHandler {
	snowflake, _ := util.NewSnowflake(1)
	return &MessageHandler{
		userRepo:       userRepo,
		messageRepo:    messageRepo,
		attachmentRepo: attachmentRepo,
		logger:         logger,
		cfg:            cfg,
		snowflake:      snowflake,
	}
}

func (h *MessageHandler) SendMessage(c *gin.Context) {
	var request struct {
		UserID      string `json:"user_id"`
		Type        string `json:"type"`
		Content     string `json:"content"`
		OriginalURL string `json:"original_url"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		h.logger.Error("Invalid request body:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	exists, err := h.userRepo.UserExists(request.UserID)
	if err != nil {
		h.logger.Error("Failed to check user:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	content := request.Content
	title := ""
	rawContent := ""
	if request.Type == "url" && request.Content != "" {
		var markdown string
		title, markdown, rawContent, err = util.FetchURLContent(request.Content)
		if err != nil {
			h.logger.Warn("Failed to fetch URL content:", err)
			// Keep original URL as content if fetch fails
		} else {
			content = markdown
		}
	}

	message := &model.Message{
		UUID:        h.snowflake.GenerateString(),
		UserID:      request.UserID,
		Type:        request.Type,
		Title:       title,
		Content:     content,
		RawContent:  rawContent,
		OriginalURL: request.OriginalURL,
	}

	if err := h.messageRepo.CreateMessage(message); err != nil {
		h.logger.Error("Failed to create message:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create message"})
		return
	}

	h.logger.Info("Message sent successfully for user:", request.UserID)
	c.JSON(http.StatusOK, gin.H{"message": "Message sent successfully", "id": message.ID})
}

func (h *MessageHandler) UploadAttachment(c *gin.Context) {
	userID := c.PostForm("user_id")
	if userID == "" {
		h.logger.Error("User ID is required")
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
		return
	}

	exists, err := h.userRepo.UserExists(userID)
	if err != nil {
		h.logger.Error("Failed to check user:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		h.logger.Error("Failed to get file:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file provided"})
		return
	}

	ext := filepath.Ext(file.Filename)
	newFilename := uuid.New().String() + ext
	savePath := filepath.Join(h.cfg.UploadPath, newFilename)

	if err := c.SaveUploadedFile(file, savePath); err != nil {
		h.logger.Error("Failed to save file:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}

	message := &model.Message{
		UUID:     h.snowflake.GenerateString(),
		UserID:   userID,
		Type:     "attachment",
		FilePath: savePath,
	}

	if err := h.messageRepo.CreateMessage(message); err != nil {
		h.logger.Error("Failed to create message:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create message"})
		return
	}

	attachment := &model.Attachment{
		UserID:    userID,
		MessageID: message.ID,
		Filename:  file.Filename,
		FilePath:  savePath,
		FileType:  strings.Split(file.Header.Get("Content-Type"), "/")[0],
	}

	if err := h.attachmentRepo.CreateAttachment(attachment); err != nil {
		h.logger.Error("Failed to create attachment:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create attachment"})
		return
	}

	h.logger.Info("Attachment uploaded successfully for user:", userID)
	c.JSON(http.StatusOK, gin.H{"message": "File uploaded successfully", "message_id": message.ID})
}

func (h *MessageHandler) SyncMessages(c *gin.Context) {
	var request struct {
		UserID         string `json:"user_id"`
		LastSyncTime   string `json:"last_sync_time"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		h.logger.Error("Invalid request body:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	exists, err := h.userRepo.UserExists(request.UserID)
	if err != nil {
		h.logger.Error("Failed to check user:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	var since time.Time
	if request.LastSyncTime != "" {
		since, err = time.Parse(time.RFC3339, request.LastSyncTime)
		if err != nil {
			h.logger.Error("Invalid last sync time format:", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid time format"})
			return
		}
	}

	messages, err := h.messageRepo.GetMessagesSince(request.UserID, since)
	if err != nil {
		h.logger.Error("Failed to get messages:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get messages"})
		return
	}

	response := make([]gin.H, 0)
	for _, msg := range messages {
		item := gin.H{
			"id":            msg.UUID,
			"type":          msg.Type,
			"title":         msg.Title,
			"content":       msg.Content,
			"raw_content":   msg.RawContent,
			"original_url":  msg.OriginalURL,
			"file_path":     msg.FilePath,
			"created_at":    msg.CreatedAt.Format(time.RFC3339),
		}
		
		if msg.Type == "attachment" {
			attachments, err := h.attachmentRepo.GetAttachmentsByMessageID(msg.ID)
			if err == nil && len(attachments) > 0 {
				item["attachment"] = gin.H{
					"filename": attachments[0].Filename,
					"file_type": attachments[0].FileType,
				}
			}
		}
		response = append(response, item)
		
		msg.SyncCompleted = true
		h.messageRepo.UpdateMessage(msg)
	}

	h.logger.Info(fmt.Sprintf("Synced %d messages for user: %s", len(response), request.UserID))
	c.JSON(http.StatusOK, response)
}

func (h *MessageHandler) GetMessageFile(c *gin.Context) {
	messageID := c.Param("id")

	msg, err := h.messageRepo.GetMessageByStringID(messageID)
	if err != nil {
		h.logger.Error("Failed to get message:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	if msg == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Message not found"})
		return
	}

	c.File(msg.FilePath)
}
