package repository

import (
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"github.com/ob-sync-server/internal/model"
	"github.com/ob-sync-server/internal/util"
)

type UserRepository struct {
	db     *gorm.DB
	logger *util.Logger
}

func NewUserRepository(db *gorm.DB, logger *util.Logger) *UserRepository {
	return &UserRepository{db: db, logger: logger}
}

func (r *UserRepository) CreateUser() (*model.User, error) {
	userID := uuid.New().String()
	
	user := &model.User{
		UserID: userID,
	}
	
	if err := r.db.Create(user).Error; err != nil {
		r.logger.Error("Failed to create user:", err)
		return nil, err
	}
	
	r.logger.Info("Created new user with ID:", userID)
	return user, nil
}

func (r *UserRepository) GetUserByID(userID string) (*model.User, error) {
	var user model.User
	if err := r.db.Where("user_id = ?", userID).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			r.logger.Debug("User not found:", userID)
			return nil, nil
		}
		r.logger.Error("Failed to get user:", err)
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) UserExists(userID string) (bool, error) {
	var count int
	if err := r.db.Model(&model.User{}).Where("user_id = ?", userID).Count(&count).Error; err != nil {
		r.logger.Error("Failed to check user existence:", err)
		return false, err
	}
	return count > 0, nil
}
