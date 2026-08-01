package dto

type TypeDrugResponse struct {
	IDJenis   uint   `json:"id_jenis"`
	NamaJenis string `json:"nama_jenis"`
}

type CreateTypeDrugRequest struct {
	NamaJenis string `json:"nama_jenis" binding:"required,max=50"`
}

type UpdateTypeDrugRequest struct {
	NamaJenis string `json:"nama_jenis" binding:"required,max=50"`
}