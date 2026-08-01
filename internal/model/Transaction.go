package model

import (
	"time"
	"medix-be/internal/user/model/entities"
)

type Transaction struct {
	IDTransaksi  uint                `gorm:"primaryKey;column:id_transaksi" json:"id_transaksi"`
	IDUser       uint                `gorm:"column:id_user;not null" json:"id_user"`
	User         model.User       `gorm:"foreignKey:IDUser" json:"user,omitempty"`
	TglTransaksi time.Time           `gorm:"column:tgl_transaksi;autoCreateTime" json:"tgl_transaksi"`
	TotalHarga   float64             `gorm:"column:total_harga;type:numeric(14,2);not null" json:"total_harga"`
	Status       string              `gorm:"column:status;size:20;default:selesai" json:"status"`
	CreatedAt    time.Time           `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	Details      []TransactionDetail `gorm:"foreignKey:IDTransaksi" json:"details,omitempty"`
}

func (Transaction) TableName() string {
	return "transaksi"
}
