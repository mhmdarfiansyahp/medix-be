package model

import (
	"medix-be/internal/drug/model/entities"
	"time"
)

type Obat struct {
	IDObat        uint           `gorm:"primaryKey;column:id_obat" json:"id_obat"`
	NamaObat      string         `gorm:"column:nama_obat;size:150;not null" json:"nama_obat"`
	MerkObat      string         `gorm:"column:merk_obat;size:100" json:"merk_obat"`
	JenisObatID   uint           `gorm:"column:jenis_obat_id" json:"jenis_obat_id"`
	JenisObat     model.TypeDrug `gorm:"foreignKey:JenisObatID;references:IDJenis;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"jenis_obat,omitempty"`
	Barcode      *string         `gorm:"column:barcode;size:50;unique" json:"barcode"`
	TglKadaluarsa time.Time      `gorm:"column:tgl_kadaluarsa;type:date;not null" json:"tgl_kadaluarsa"`
	Harga         float64        `gorm:"column:harga;type:numeric(12,2);not null" json:"harga"`
	Stok          int            `gorm:"column:stok;not null;default:0" json:"stok"`
	StokMinimum   int            `gorm:"column:stok_minimum;default:10" json:"stok_minimum"`
	Keterangan    string         `gorm:"column:keterangan" json:"keterangan"`
	Status        int            `gorm:"column:status;default:1" json:"status"`
	Gambar        string         `gorm:"column:gambar;size:255" json:"gambar"`
	CreatedAt     time.Time      `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time      `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
}

func (Obat) TableName() string {
	return "obat"
}
