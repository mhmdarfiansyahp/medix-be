package model

import "time"

type User struct {
	IDUser    uint      `gorm:"primaryKey;column:id_user" json:"id_user"`
	NamaUser  string    `gorm:"column:nama_user;size:100;not null" json:"nama_user"`
	NoTelp    string    `gorm:"column:no_telp;size:20" json:"no_telp"`
	Role      string    `gorm:"column:role;size:20;not null" json:"role"`
	Username  string    `gorm:"column:username;size:50;not null;unique" json:"username"`
	Password  string    `gorm:"column:password;size:255;not null" json:"-"` // Hidden dari response JSON
	Status    string    `gorm:"column:status;size:20;default:aktif" json:"status"`
	Foto      string    `gorm:"column:foto;size:255" json:"foto"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
}

func (User) TableName() string {
	return "users"
}