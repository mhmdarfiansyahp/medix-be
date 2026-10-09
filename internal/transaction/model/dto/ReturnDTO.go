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
	Jumlah       int     `json:"jumlah"`
	AlasanItem   string  `json:"alasan_item,omitempty"`
	KondisiLayak bool    `json:"kondisi_layak"`
	StokKembali  int     `json:"stok_kembali"`
}

type CreateReturnResponse struct {
	IDReturn     uint              `json:"id_return"`
	IDTransaksi  uint              `json:"id_transaksi"`
	Alasan       string            `json:"alasan"`
	TanggalRetur string            `json:"tanggal_retur"`
	DiajukanOleh uint              `json:"diajukan_oleh"`
	Status       string            `json:"status"`
	Items        []ReturnItemResponse `json:"items"`
}