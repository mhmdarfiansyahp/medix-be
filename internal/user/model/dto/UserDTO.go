package dto

type UserResponse struct {
	IDUser    uint   `json:"id_user"`
	NamaUser  string `json:"nama_user"`
	NoTelp    string `json:"no_telp"`
	Role      string `json:"role"`
	Username  string `json:"username"`
	Status    string `json:"status"`
	Foto      string `json:"foto"`
}

type CreateUserRequest struct {
	NamaUser string `json:"nama_user" binding:"required,max=100"`
	NoTelp   string `json:"no_telp" binding:"max=13"`
	Role     string `json:"role" binding:"required,oneof=admin kasir apoteker"`
	Username string `json:"username" binding:"required,max=50"`
	Password string `json:"password" binding:"required,min=6"`
	Status string `json:"status" binding:"omitempty,oneof=aktif nonaktif"`
	Foto   string `json:"foto"`
}

type UpdateUserRequest struct {
	NamaUser string `json:"nama_user" binding:"omitempty,max=100"`
	NoTelp   string `json:"no_telp" binding:"omitempty,max=13"`
	Role     string `json:"role" binding:"omitempty,oneof=admin kasir apoteker"`
	Username string `json:"username" binding:"omitempty,max=50"`
	Password string `json:"password" binding:"omitempty,min=6"`
	Status string `json:"status" binding:"omitempty,oneof=aktif nonaktif"`
	Foto   string `json:"foto"`
}