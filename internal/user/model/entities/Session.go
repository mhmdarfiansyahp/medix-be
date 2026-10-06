package model

import "time"

type Session struct {
	IDSession    string    `gorm:"primaryKey;column:id_session;type:uuid;default:uuid_generate_v4()" json:"id_session"`
	IDUser       uint      `gorm:"column:id_user;not null" json:"id_user"`
	RefreshToken string    `gorm:"column:refresh_token;size:255;not null;uniqueIndex" json:"refresh_token"`
	ExpiresAt    time.Time `gorm:"column:expires_at;not null" json:"expires_at"`
	Revoked      bool      `gorm:"column:revoked;default:false" json:"revoked"`
	CreatedAt    time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	User         User      `gorm:"foreignKey:IDUser;references:IDUser" json:"-"`
}

func (Session) TableName() string {
	return "sessions"
}
