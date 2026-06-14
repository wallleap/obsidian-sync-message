package repository

import (
	"github.com/jinzhu/gorm"
	"github.com/ob-sync-server/internal/model"
	"github.com/ob-sync-server/internal/util"
)

type AttachmentRepository struct {
	db     *gorm.DB
	logger *util.Logger
}

func NewAttachmentRepository(db *gorm.DB, logger *util.Logger) *AttachmentRepository {
	return &AttachmentRepository{db: db, logger: logger}
}

func (r *AttachmentRepository) CreateAttachment(attachment *model.Attachment) error {
	if err := r.db.Create(attachment).Error; err != nil {
		r.logger.Error("Failed to create attachment:", err)
		return err
	}
	r.logger.Info("Created attachment for message:", attachment.MessageID)
	return nil
}

func (r *AttachmentRepository) GetAttachmentsByMessageID(messageID uint) ([]*model.Attachment, error) {
	var attachments []*model.Attachment
	if err := r.db.Where("message_id = ?", messageID).Find(&attachments).Error; err != nil {
		r.logger.Error("Failed to get attachments:", err)
		return nil, err
	}
	return attachments, nil
}

func (r *AttachmentRepository) GetAttachmentsByUserID(userID string) ([]*model.Attachment, error) {
	var attachments []*model.Attachment
	if err := r.db.Where("user_id = ?", userID).Find(&attachments).Error; err != nil {
		r.logger.Error("Failed to get attachments:", err)
		return nil, err
	}
	return attachments, nil
}

func (r *AttachmentRepository) DeleteAttachment(id uint) error {
	if err := r.db.Delete(&model.Attachment{}, id).Error; err != nil {
		r.logger.Error("Failed to delete attachment:", err)
		return err
	}
	return nil
}
