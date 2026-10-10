package dto

import "time"

type DetailItemRequest struct {
	IDObat uint `json:"id_obat" binding:"required"`
	Jumlah int  `json:"jumlah" binding:"required,gt=0"`
}

type CreateTransactionRequest struct {
	IDUser  uint                `json:"id_user"`
	Details []DetailItemRequest `json:"details" binding:"required,gt=0,dive"`
}

type AddToCartRequest struct {
	IDObat uint `json:"id_obat" binding:"required"`
	Jumlah int  `json:"jumlah" binding:"required,gt=0"`
}

type AddToCartResponse struct {
	Success      bool                `json:"success"`
	Transaction  *TransactionResponse `json:"transaction,omitempty"`
	Message      string              `json:"message,omitempty"`
}

type PaymentRequest struct {
	MetodeBayar    string  `json:"metode_bayar" binding:"required,oneof=tunai QRIS debit kredit transfer"`
	UangDiterima   float64 `json:"uang_diterima" binding:"omitempty,gte=0"`
}

type PaymentResponse struct {
	IDTransaksi  uint      `json:"id_transaksi"`
	MetodeBayar    string  `json:"metode_bayar"`
	UangDiterima   float64 `json:"uang_diterima"`
	Kembalian      float64 `json:"kembalian"`
	Status       int      `json:"status"`
	TotalHarga   float64  `json:"total_harga"`
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
	TotalTunai     float64 `json:"total_tunai"`
	TotalNonTunai  float64 `json:"total_non_tunai"`
}

type TodayTransactionResponse struct {
	Transactions []TransactionResponse      `json:"transactions"`
	Summary      TransactionSummaryResponse `json:"summary"`
	Returns      []ReturnSummaryResponse    `json:"returns"`
}

type ReturnSummaryResponse struct {
	IDReturn     uint      `json:"id_return"`
	IDTransaksi  uint      `json:"id_transaksi"`
	Alasan       string    `json:"alasan"`
	TanggalRetur time.Time `json:"tanggal_retur"`
	Status       string    `json:"status"`
}

type TransactionResponse struct {
	IDTransaksi  uint                 `json:"id_transaksi"`
	IDUser       uint                 `json:"id_user"`
	TglTransaksi time.Time            `json:"tgl_transaksi"`
	TotalHarga   float64              `json:"total_harga"`
	Status       int                  `json:"status"`
	Details      []DetailItemResponse `json:"details,omitempty"`
}

// Receipt response (same structure as TransactionResponse for now)
type ReceiptResponse struct {
	IDTransaksi  uint                 `json:"id_transaksi"`
	TglTransaksi time.Time            `json:"tgl_transaksi"`
	IDUser       uint                 `json:"id_user"`
	Details      []DetailItemResponse `json:"details"`
	TotalHarga   float64              `json:"total_harga"`
	MetodeBayar  string               `json:"metode_bayar"`
	UangDiterima float64              `json:"uang_diterima"`
	Kembalian    float64              `json:"kembalian"`
}

// Receipt generation response
type GenerateReceiptResponse struct {
	ReceiptID uint    `json:"receipt_id"`
	Message  string   `json:"message"`
}

// User activity response for kasir report
type TodayKasirReportResponse struct {
	Transactions []TransactionResponse      `json:"transactions"`
	Returns      []ReturnSummaryResponse    `json:"returns"`
	Summary      TransactionSummaryResponse `json:"summary"`
}

// Reject return request
type RejectReturnRequest struct {
	Alasan string `json:"alasan" binding:"required"`
}

// Reject return response
type RejectReturnResponse struct {
	IDReturn uint   `json:"id_return"`
	Status   string `json:"status"`
	Alasan   string `json:"alasan"`
	Message  string `json:"message"`
}

// Approval threshold request
type ApprovalThresholdRequest struct {
	Threshold float64 `json:"threshold" binding:"required,gte=0"`
}

// Approval threshold response
type ApprovalThresholdResponse struct {
	Threshold float64 `json:"threshold"`
	Message   string  `json:"message"`
}
