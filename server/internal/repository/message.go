package repository

import (
	"time"

	"github.com/jinzhu/gorm"
	"github.com/ob-sync-server/internal/model"
	"github.com/ob-sync-server/internal/util"
)

type MessageRepository struct {
	db     *gorm.DB
	logger *util.Logger
}

func NewMessageRepository(db *gorm.DB, logger *util.Logger) *MessageRepository {
	return &MessageRepository{db: db, logger: logger}
}

func (r *MessageRepository) CreateMessage(message *model.Message) error {
	if err := r.db.Create(message).Error; err != nil {
		r.logger.Error("Failed to create message:", err)
		return err
	}
	r.logger.Info("Created message for user:", message.UserID)
	return nil
}

func (r *MessageRepository) GetMessagesSince(userID string, since time.Time) ([]*model.Message, error) {
	var messages []*model.Message
	r.logger.Info("Getting messages since:", since)
	
	if err := r.db.Where("user_id = ? AND created_at > ?", userID, since).Order("created_at ASC").Find(&messages).Error; err != nil {
		r.logger.Error("Failed to get messages:", err)
		return nil, err
	}
	r.logger.Info("Found", len(messages), "messages")
	return messages, nil
}

func (r *MessageRepository) GetMessageByID(id uint) (*model.Message, error) {
	var message model.Message
	if err := r.db.First(&message, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		r.logger.Error("Failed to get message:", err)
		return nil, err
	}
	return &message, nil
}

func (r *MessageRepository) UpdateMessage(message *model.Message) error {
	if err := r.db.Save(message).Error; err != nil {
		r.logger.Error("Failed to update message:", err)
		return err
	}
	return nil
}

func (r *MessageRepository) DeleteMessage(id uint) error {
	if err := r.db.Delete(&model.Message{}, id).Error; err != nil {
		r.logger.Error("Failed to delete message:", err)
		return err
	}
	return nil
}

func (r *MessageRepository) GetUnsyncedMessages(userID string) ([]*model.Message, error) {
	var messages []*model.Message
	if err := r.db.Where("user_id = ? AND sync_completed = ?", userID, false).Order("created_at ASC").Find(&messages).Error; err != nil {
		r.logger.Error("Failed to get unsynced messages:", err)
		return nil, err
	}
	return messages, nil
}

func (r *MessageRepository) GetMessageByStringID(id string) (*model.Message, error) {
	var message model.Message
	if err := r.db.Where("id = ?", id).First(&message).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		r.logger.Error("Failed to get message:", err)
		return nil, err
	}
	return &message, nil
}
