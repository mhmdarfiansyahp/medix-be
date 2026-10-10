package entities

import "time"

type Return struct {
	IDReturn        uint         `gorm:"primaryKey;column:id_return" json:"id_return"`
	IDTransaksi     uint         `gorm:"column:id_transaksi;not null" json:"id_transaksi"`
	Alasan          string       `gorm:"column:alasan" json:"alasan"`
	AlasanPenolakan *string      `gorm:"column:alasan_penolakan" json:"alasan_penolakan,omitempty"`
	TanggalRetur    time.Time    `gorm:"column:tanggal_retur;autoCreateTime" json:"tanggal_retur"`
	DiajukanOleh    uint         `gorm:"column:diajukan_oleh;not null" json:"diajukan_oleh"`
	Status          string       `gorm:"column:status;default:'pending'" json:"status"` // pending, approved, rejected
	Items           []ReturnItem `gorm:"foreignKey:IDReturn;references:IDReturn" json:"items"`
}

func (Return) TableName() string {
	return "transaksi_retur"
}

type Pengaturan struct {
	Key   string `gorm:"primaryKey;column:key" json:"key"`
	Value string `gorm:"column:value" json:"value"`
}

func (Pengaturan) TableName() string {
	return "pengaturan"
}

type ReturnItem struct {
	IDItem       uint    `gorm:"primaryKey;column:id_item" json:"id_item"`
	IDReturn     uint    `gorm:"column:id_return;not null" json:"id_return"`
	IDObat       uint    `gorm:"column:id_obat;not null" json:"id_obat"`
	Jumlah       int     `gorm:"column:jumlah;not null" json:"jumlah"`
	AlasanItem   string  `gorm:"column:alasan_item" json:"alasan_item"`
	KondisiLayak bool    `gorm:"column:kondisi_layak" json:"kondisi_layak"` // true if layak jual, false if tidak layak
	StokKembali  int     `gorm:"column:stok_kembali" json:"stok_kembali"` // jumlah yang actually dikembalikan ke stok (jika layak)
}

func (ReturnItem) TableName() string {
	return "transaksi_retur_item"
}