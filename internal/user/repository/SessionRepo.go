package repository

import (
	model "medix-be/internal/user/model/entities"

	"gorm.io/gorm"
)

type SessionRepository interface {
	Create(session *model.Session) error
	FindByRefreshToken(token string) (*model.Session, error)
	RevokeByID(id string) error
	RevokeAllByUserID(userID uint) error
	DeleteExpired() error
}

type sessionRepository struct {
	db *gorm.DB
}

func NewSessionRepository(db *gorm.DB) SessionRepository {
	return &sessionRepository{db: db}
}

func (r *sessionRepository) Create(session *model.Session) error {
	return r.db.Create(session).Error
}

func (r *sessionRepository) FindByRefreshToken(token string) (*model.Session, error) {
	var session model.Session
	err := r.db.Where("refresh_token = ? AND revoked = false AND expires_at > NOW()", token).First(&session).Error
	if err != nil {
		return nil, err
	}
	return &session, nil
}

func (r *sessionRepository) RevokeByID(id string) error {
	return r.db.Model(&model.Session{}).Where("id_session = ?", id).Update("revoked", true).Error
}

func (r *sessionRepository) RevokeAllByUserID(userID uint) error {
	return r.db.Model(&model.Session{}).Where("id_user = ? AND revoked = false", userID).Update("revoked", true).Error
}

func (r *sessionRepository) DeleteExpired() error {
	return r.db.Where("expires_at < NOW() AND revoked = true").Delete(&model.Session{}).Error
}
