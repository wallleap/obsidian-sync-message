package model

import (
	"time"

	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model
	UserID string `gorm:"unique;not null;size:36"`
}

type Message struct {
	gorm.Model
	UUID          string    `gorm:"size:36;index"` // Snowflake ID for external reference
	UserID        string    `gorm:"not null;index"`
	Type          string    `gorm:"not null;size:20"`
	Title         string    `gorm:"size:500"`
	Content       string    `gorm:"type:text"`
	RawContent    string    `gorm:"type:text"` // 原始网页HTML内容，方便调试
	OriginalURL   string    `gorm:"size:500"`
	FilePath      string    `gorm:"size:500"`
	ProcessedAt   time.Time
	SyncCompleted bool `gorm:"default:false"`
}

type Attachment struct {
	gorm.Model
	UserID    string `gorm:"not null;index"`
	MessageID uint
	Filename  string `gorm:"not null"`
	FilePath  string `gorm:"not null"`
	FileType  string `gorm:"size:50"`
}

func (u *User) TableName() string {
	return "users"
}

func (m *Message) TableName() string {
	return "messages"
}

func (a *Attachment) TableName() string {
	return "attachments"
}
