package dto

import "time"

type DetailItemRequest struct {
	IDObat      uint    `json:"id_obat" binding:"required"`
	Jumlah      int     `json:"jumlah" binding:"required,gt=0"`
	HargaSatuan float64 `json:"harga_satuan" binding:"required,gt=0"`
}

type CreateTransactionRequest struct {
	IDUser  uint                `json:"id_user" binding:"required"`
	Details []DetailItemRequest `json:"details" binding:"required,gt=0,dive"`
}

type DetailItemResponse struct {
	IDDetail    uint    `json:"id_detail"`
	IDObat      uint    `json:"id_obat"`
	Jumlah      int     `json:"jumlah"`
	HargaSatuan float64 `json:"harga_satuan"`
	Subtotal    float64 `json:"subtotal"`
}

type TransactionResponse struct {
	IDTransaksi  uint                 `json:"id_transaksi"`
	IDUser       uint                 `json:"id_user"`
	TglTransaksi time.Time            `json:"tgl_transaksi"`
	TotalHarga   float64              `json:"total_harga"`
	Status       string               `json:"status"`
	Details      []DetailItemResponse `json:"details,omitempty"`
}