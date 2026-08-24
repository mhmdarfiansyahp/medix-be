package dto

type UserResponse struct {
	IDUser   uint   `json:"id_user"`
	NamaUser string `json:"nama_user"`
	NoTelp   string `json:"no_telp"`
	Role     string `json:"role"`
	Username string `json:"username"`
	Status   string `json:"status"`
	Foto     string `json:"foto"`
}

type CreateUserRequest struct {
	NamaUser string `json:"nama_user" binding:"required,max=100"`
	NoTelp   string `json:"no_telp" binding:"max=13"`
	Role     string `json:"role" binding:"required,oneof=admin kasir owner"`
	Username string `json:"username" binding:"required,max=50"`
	Password string `json:"password" binding:"required,min=6"`
	Status   string `json:"status" binding:"omitempty,oneof=aktif nonaktif"`
	Foto     string `json:"foto"`
}

type UpdateUserRequest struct {
	NamaUser string `json:"nama_user" binding:"omitempty,max=100"`
	NoTelp   string `json:"no_telp" binding:"omitempty,max=13"`
	Role     string `json:"role" binding:"omitempty,oneof=admin kasir owner"`
	Username string `json:"username" binding:"omitempty,max=50"`
	Password string `json:"password" binding:"omitempty,min=6"`
	Status   string `json:"status" binding:"omitempty,oneof=aktif nonaktif"`
	Foto     string `json:"foto"`
}

type UserFilterRequest struct {
	Page   int    `form:"page"`
	Limit  int    `form:"limit"`
	Search string `form:"search"`
	Role   string `form:"role"`
	Status string `form:"status"`
}

type PaginationResponse struct {
	CurrentPage  int   `json:"current_page"`
	TotalPages   int   `json:"total_pages"`
	TotalItems   int64 `json:"total_items"`
	ItemsPerPage int   `json:"items_per_page"`
}

type UserListResponse struct {
	Data       []*UserResponse    `json:"data"`
	Pagination PaginationResponse `json:"pagination"`
}
