package dto

import (
	"encoding/json"
	"regexp"
)

// UserStatus accepts both string ("aktif"/"nonaktif") and numeric (1/0)
// representations of the user status field, so callers sending either shape
// do not trigger an unmarshalling error.
type UserStatus string

func (s *UserStatus) UnmarshalJSON(data []byte) error {
	// Try numeric representation first (1 = active, 0 = inactive).
	var num int
	if err := json.Unmarshal(data, &num); err == nil {
		if num == 1 {
			*s = "aktif"
		} else {
			*s = "nonaktif"
		}
		return nil
	}

	// Fall back to the string representation.
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}
	*s = UserStatus(str)
	return nil
}

type UserResponse struct {
	IDUser             uint   `json:"id_user"`
	NamaUser           string `json:"nama_user"`
	NoTelp             string `json:"no_telp"`
	Role               string `json:"role"`
	Username           string `json:"username"`
	Status             string `json:"status"`
	Foto               string `json:"foto"`
	MustChangePassword bool   `json:"must_change_password"`
}

type CreateUserRequest struct {
	NamaUser string     `json:"nama_user" binding:"required,max=100"`
	NoTelp   string     `json:"no_telp" binding:"max=15"`
	Role     string     `json:"role" binding:"required,oneof=admin kasir owner"`
	Username string     `json:"username" binding:"required,max=50"`
	Password string     `json:"password" binding:"required,min=6"`
	Status   UserStatus `json:"status" binding:"omitempty,oneof=aktif nonaktif"`
	Foto     string     `json:"foto"`
}

type UpdateUserRequest struct {
	NamaUser string     `json:"nama_user" binding:"omitempty,max=100"`
	NoTelp   string     `json:"no_telp" binding:"omitempty,max=15"`
	Role     string     `json:"role" binding:"omitempty,oneof=admin kasir owner"`
	Username string     `json:"username" binding:"omitempty,max=50"`
	Password string     `json:"password" binding:"omitempty,min=6"`
	Status   UserStatus `json:"status" binding:"omitempty,oneof=aktif nonaktif"`
	Foto     string     `json:"foto"`
}

type UpdateProfileRequest struct {
	NamaUser string `json:"nama_user" form:"nama_user" binding:"omitempty,max=100"`
	NoTelp   string `json:"no_telp" form:"no_telp" binding:"omitempty,max=15"`
	Username string `json:"username" form:"username" binding:"omitempty,max=50"`
	Foto     string `json:"foto" form:"foto"`
}

func (r *UpdateProfileRequest) ValidateNoTelp() bool {
	if r.NoTelp == "" {
		return true
	}
	pattern := `^(\+62|0)[0-9]{9,12}$`
	match, _ := regexp.MatchString(pattern, r.NoTelp)
	return match
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

type LoginRequest struct {
	Username *string `json:"username" binding:"required"`
	Password *string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token        *string       `json:"token"`
	RefreshToken *string       `json:"refresh_token"`
	User         *UserResponse `json:"user"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type RefreshResponse struct {
	Token        *string `json:"token"`
	RefreshToken *string `json:"refresh_token"`
}

type LogoutRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type ResetPasswordResponse struct {
	Password string `json:"password"`
	Message  string `json:"message"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}