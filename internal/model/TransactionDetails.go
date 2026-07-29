package model

type TransactionDetail struct {
	IDDetail    uint    `gorm:"primaryKey;column:id_detail" json:"id_detail"`
	IDTransaksi uint    `gorm:"column:id_transaksi;not null" json:"id_transaksi"`
	IDObat      uint    `gorm:"column:id_obat;not null" json:"id_obat"`
	Obat        Obat    `gorm:"foreignKey:IDObat" json:"obat,omitempty"`
	Jumlah      int     `gorm:"column:jumlah;not null" json:"jumlah"`
	HargaSatuan float64 `gorm:"column:harga_satuan;type:numeric(12,2);not null" json:"harga_satuan"`
	Subtotal    float64 `gorm:"column:subtotal;type:numeric(14,2);not null" json:"subtotal"`
}

func (TransactionDetail) TableName() string {
	return "detail_pembelian"
}
