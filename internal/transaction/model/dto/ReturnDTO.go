package dto

type ReturnItemRequest struct {
	IDObat       uint   `json:"id_obat" binding:"required"`
	Jumlah       int    `json:"jumlah" binding:"required,gt=0"`
	AlasanItem   string `json:"alasan_item,omitempty"`
	KondisiLayak bool   `json:"kondisi_layak" binding:"required"` // true if layak jual
}

type CreateReturnRequest struct {
	IDTransaksi uint              `json:"id_transaksi" binding:"required"`
	Alasan      string            `json:"alasan" binding:"required"`
	Items       []ReturnItemRequest `json:"items" binding:"required,gt=0,dive"`
}

type ReturnItemResponse struct {
	IDItem       uint    `json:"id_item"`
	IDObat       uint    `json:"id_obat"`
	NamaObat     string  `json:"nama_obat"`
	Jumlah       int     `json:"jumlah"`
	HargaSatuan  float64 `json:"harga_satuan"`
	Subtotal     float64 `json:"subtotal"`
	AlasanItem   string  `json:"alasan_item,omitempty"`
	KondisiLayak bool    `json:"kondisi_layak"`
	StokKembali  int     `json:"stok_kembali"`
}

type CreateReturnResponse struct {
	IDReturn        uint                 `json:"id_return"`
	IDTransaksi     uint                 `json:"id_transaksi"`
	Alasan          string               `json:"alasan"`
	AlasanPenolakan *string              `json:"alasan_penolakan,omitempty"`
	TanggalRetur    string               `json:"tanggal_retur"`
	DiajukanOleh    uint                 `json:"diajukan_oleh"`
	Status          string               `json:"status"`
	TotalNilai      float64              `json:"total_nilai"`
	Items           []ReturnItemResponse `json:"items"`
}

// Approve return response
type ApproveReturnResponse struct {
	IDReturn uint    `json:"id_return"`
	Status   string  `json:"status"`
	Message  string  `json:"message"`
}