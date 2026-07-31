package model

type TypeDrug struct {
	IDJenis   uint   `gorm:"primaryKey;column:id_jenis" json:"id_jenis"`
	NamaJenis string `gorm:"column:nama_jenis;size:50;not null;unique" json:"nama_jenis"`
}

func (TypeDrug) TableName() string {
	return "jenis_obat"
}