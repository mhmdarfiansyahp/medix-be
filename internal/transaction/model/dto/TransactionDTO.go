package dto

import "time"

type DetailItemRequest struct {
	IDObat uint `json:"id_obat" binding:"required"`
	Jumlah int  `json:"jumlah" binding:"required,gt=0"`
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

type CancelTransactionResponse struct {
	IDTransaksi uint    `json:"id_transaksi"`
	Status      int     `json:"status"`
	TotalHarga  float64 `json:"total_harga"`
	Message     string  `json:"message"`
}

type TransactionSummaryResponse struct {
	TotalTransaksi int     `json:"total_transaksi"`
	TotalPenjualan float64 `json:"total_penjualan"`
}

type TodayTransactionResponse struct {
	Transactions []TransactionResponse      `json:"transactions"`
	Summary      TransactionSummaryResponse `json:"summary"`
}

type TransactionResponse struct {
	IDTransaksi  uint                 `json:"id_transaksi"`
	IDUser       uint                 `json:"id_user"`
	TglTransaksi time.Time            `json:"tgl_transaksi"`
	TotalHarga   float64              `json:"total_harga"`
	Status       int                  `json:"status"`
	Details      []DetailItemResponse `json:"details,omitempty"`
}

type ReceiptResponse struct {
	IDTransaksi  uint                 `json:"id_transaksi"`
	TglTransaksi time.Time            `json:"tgl_transaksi"`
	IDUser       uint                 `json:"id_user"`
	Details      []DetailItemResponse `json:"details"`
	TotalHarga   float64              `json:"total_harga"`
}
