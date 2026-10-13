package model

import (
	"time"
)

type PasswordResetToken struct {
	ID        uint      `gorm:"primaryKey;column:id" json:"id"`
	IDUser    uint      `gorm:"column:id_user;not null" json:"id_user"`
	Token     string    `gorm:"column:token;size:255;not null;unique" json:"token"`
	ExpiresAt time.Time `gorm:"column:expires_at;not null" json:"expires_at"`
	Used      bool      `gorm:"column:used;default:false" json:"used"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	User      User      `gorm:"foreignKey:IDUser;references:IDUser" json:"-"`
}

func (PasswordResetToken) TableName() string {
	return "password_reset_tokens"
}
