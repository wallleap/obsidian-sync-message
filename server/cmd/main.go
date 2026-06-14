package main

import (
	"fmt"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/ob-sync-server/config"
	"github.com/ob-sync-server/internal/handler"
	"github.com/ob-sync-server/internal/model"
	"github.com/ob-sync-server/internal/repository"
	"github.com/ob-sync-server/internal/util"
)

func main() {
	cfg := config.LoadConfig()
	logger := util.NewLogger(cfg.LogPath)

	db, err := gorm.Open("sqlite3", cfg.DBPath)
	if err != nil {
		logger.Error("Failed to connect to database:", err)
		return
	}
	defer db.Close()

	db.AutoMigrate(&model.User{}, &model.Message{}, &model.Attachment{})

	logger.Info("Database connected and migrated successfully")

	userRepo := repository.NewUserRepository(db, logger)
	messageRepo := repository.NewMessageRepository(db, logger)
	attachmentRepo := repository.NewAttachmentRepository(db, logger)

	userHandler := handler.NewUserHandler(userRepo, logger)
	messageHandler := handler.NewMessageHandler(userRepo, messageRepo, attachmentRepo, logger, cfg)

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"POST", "GET", "OPTIONS"},
		AllowHeaders: []string{"Content-Type", "Authorization"},
	}))

	r.POST("/api/user/generate", userHandler.GenerateUserID)
	r.POST("/api/user/validate", userHandler.ValidateUserID)

	r.POST("/api/message/send", messageHandler.SendMessage)
	r.POST("/api/message/upload", messageHandler.UploadAttachment)
	r.POST("/api/message/sync", messageHandler.SyncMessages)
	r.GET("/api/message/file/:id", messageHandler.GetMessageFile)

	logger.Info(fmt.Sprintf("Server starting on port %s...", cfg.ServerPort))
	if err := r.Run(":" + cfg.ServerPort); err != nil {
		logger.Error("Failed to start server:", err)
	}
}
