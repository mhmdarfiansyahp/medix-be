package repository

import (
	model "medix-be/internal/user/model/entities"

	"gorm.io/gorm"
)

type PasswordResetTokenRepository interface {
	Create(token *model.PasswordResetToken) error
	FindByToken(token string) (*model.PasswordResetToken, error)
	MarkAsUsed(id uint) error
	InvalidateAllByUserID(userID uint) error
}

type passwordResetTokenRepository struct {
	db *gorm.DB
}

func NewPasswordResetTokenRepository(db *gorm.DB) PasswordResetTokenRepository {
	return &passwordResetTokenRepository{db: db}
}

func (r *passwordResetTokenRepository) Create(token *model.PasswordResetToken) error {
	return r.db.Create(token).Error
}

func (r *passwordResetTokenRepository) FindByToken(token string) (*model.PasswordResetToken, error) {
	var record model.PasswordResetToken
	err := r.db.Where("token = ? AND used = false AND expires_at > NOW()", token).First(&record).Error
	if err != nil {
		return nil, err
	}
	return &record, nil
}

func (r *passwordResetTokenRepository) MarkAsUsed(id uint) error {
	return r.db.Model(&model.PasswordResetToken{}).Where("id = ?", id).Update("used", true).Error
}

func (r *passwordResetTokenRepository) InvalidateAllByUserID(userID uint) error {
	return r.db.Model(&model.PasswordResetToken{}).Where("id_user = ? AND used = false", userID).Update("used", true).Error
}
