package dto

import "time"

type MedicineFilterParams struct {
	Search      string `form:"search"`
	JenisObatID uint   `form:"jenis_obat_id"`
	StatusStok  string `form:"status_stok"`
	Status      string `form:"status"`
	Barcode     string `form:"barcode"`
}
type CreateMedicineRequest struct {
	NamaObat      string  `json:"nama_obat" binding:"required,max=150"`
	MerkObat      string  `json:"merk_obat" binding:"max=100"`
	JenisObatID   uint    `json:"jenis_obat_id" binding:"required"`
	Barcode       *string `json:"barcode" binding:"max=50"`
	TglKadaluarsa string  `json:"tgl_kadaluarsa" binding:"required"`
	Harga         float64 `json:"harga" binding:"required,gte=0"`
	Stok          int     `json:"stok" binding:"gte=0"`
	StokMinimum   int     `json:"stok_minimum"`
	Keterangan    string  `json:"keterangan"`
	Status        string  `json:"status" binding:"omitempty,oneof=aktif nonaktif"`
	Gambar        string  `json:"gambar"`
}

type UpdateMedicineRequest struct {
	NamaObat      string  `json:"nama_obat" binding:"omitempty,max=150"`
	MerkObat      string  `json:"merk_obat" binding:"max=100"`
	JenisObatID   uint    `json:"jenis_obat_id"`
	Barcode       *string `json:"barcode" binding:"max=50"`
	TglKadaluarsa string  `json:"tgl_kadaluarsa"`
	Harga         float64 `json:"harga" binding:"omitempty,gte=0"`
	Stok          int     `json:"stok" binding:"omitempty,gte=0"`
	StokMinimum   int     `json:"stok_minimum"`
	Keterangan    string  `json:"keterangan"`
	Status        string  `json:"status" binding:"omitempty,oneof=aktif nonaktif"`
	Gambar        string  `json:"gambar"`
}

type MedicineResponse struct {
	IDObat        uint        `json:"id_obat"`
	NamaObat      string      `json:"nama_obat"`
	MerkObat      string      `json:"merk_obat"`
	JenisObatID   uint        `json:"jenis_obat_id"`
	JenisObat     interface{} `json:"jenis_obat,omitempty"`
	Barcode       *string     `json:"barcode"`
	TglKadaluarsa string      `json:"tgl_kadaluarsa"`
	Harga         float64     `json:"harga"`
	Stok          int         `json:"stok"`
	StokMinimum   int         `json:"stok_minimum"`
	Keterangan    string      `json:"keterangan"`
	Status        string      `json:"status"`
	Gambar        string      `json:"gambar"`
	CreatedAt     time.Time   `json:"created_at"`
	UpdatedAt     time.Time   `json:"updated_at"`
}
