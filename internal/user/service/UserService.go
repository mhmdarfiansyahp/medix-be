package service

import (
	"errors"
	"math"
	"medix-be/internal/user/model/dto"
	model "medix-be/internal/user/model/entities"
	"medix-be/internal/user/repository"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	CreateUser(req dto.CreateUserRequest) (*dto.UserResponse, error)
	Login(req *dto.LoginRequest) (*dto.LoginResponse, error)

	GetAllUsers(page int, limit int, search string, role string, status string) (*dto.UserListResponse, error)
	GetUserByID(id uint) (*dto.UserResponse, error)
	UpdateUser(id uint, req dto.UpdateUserRequest) (*dto.UserResponse, error)
	DeleteUser(id uint) error
	GetProfile(id uint) (*dto.UserResponse, error)
	UpdateProfile(id uint, req *dto.UpdateProfileRequest) (*dto.UserResponse, error)
	UpdateProfilePhoto(id uint, fotoPath string) (*dto.UserResponse, error)
}

type userService struct {
	repo      repository.UserRepository
	jwtSecret string
}

func NewUserService(repo repository.UserRepository, jwtSecret string) UserService {
	return &userService{repo: repo, jwtSecret: jwtSecret}
}

func (s *userService) CreateUser(req dto.CreateUserRequest) (*dto.UserResponse, error) {

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)

	if err != nil {
		return nil, errors.New("failed to process password")
	}

	status := req.Status
	if status == "" {
		status = "aktif"
	}

	user := model.User{
		NamaUser: req.NamaUser,
		NoTelp:   req.NoTelp,
		Role:     req.Role,
		Username: req.Username,
		Password: string(hashedPassword),
		Status:   statusToInt(string(status)),
		Foto:     req.Foto,
	}

	if err := s.repo.Create(&user); err != nil {
		return nil, err
	}

	res := toUserResponse(user)
	return &res, nil
}

func (s *userService) GetAllUsers(page int, limit int, search string, role string, status string) (*dto.UserListResponse, error) {

	users, total, err := s.repo.FindAll(
		page,
		limit,
		search,
		role,
		normalizeStatusFilter(status),
	)

	if err != nil {
		return nil, err
	}

	responses := make([]*dto.UserResponse, 0, len(users))

	for _, user := range users {
		res := toUserResponse(*user)
		responses = append(responses, &res)
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	return &dto.UserListResponse{
		Data: responses,
		Pagination: dto.PaginationResponse{
			CurrentPage:  page,
			TotalPages:   totalPages,
			TotalItems:   total,
			ItemsPerPage: limit,
		},
	}, nil
}

func (s *userService) GetUserByID(id uint) (*dto.UserResponse, error) {
	user, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("user not found")
	}

	res := toUserResponse(*user)
	return &res, nil
}

func (s *userService) UpdateUser(id uint, req dto.UpdateUserRequest) (*dto.UserResponse, error) {

	user, err := s.repo.FindByID(id)

	if err != nil {
		return nil, errors.New("user not found")
	}

	if req.NamaUser != "" {
		user.NamaUser = req.NamaUser
	}
	if req.NoTelp != "" {
		user.NoTelp = req.NoTelp
	}
	if req.Role != "" {
		user.Role = req.Role
	}
	if req.Username != "" {
		user.Username = req.Username
	}
	if req.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, errors.New("failed to process new password")
		}
		user.Password = string(hashedPassword)
	}
	if req.Status != "" {
		user.Status = statusToInt(string(req.Status))
	}
	if req.Foto != "" {
		user.Foto = req.Foto
	}

	if err := s.repo.Update(user); err != nil {
		return nil, err
	}

	res := toUserResponse(*user)
	return &res, nil
}

func (s *userService) DeleteUser(id uint) error {
	_, err := s.repo.FindByID(id)
	if err != nil {
		return errors.New("user tidak ditemukan")
	}

	return s.repo.Delete(id)
}

func (s *userService) GetProfile(id uint) (*dto.UserResponse, error) {

	user, err := s.repo.FindByID(id)

	if err != nil {
		return nil, errors.New("profile not found")
	}

	res := toUserResponse(*user)
	return &res, nil
}

func (s *userService) Login(req *dto.LoginRequest) (*dto.LoginResponse, error) {

	user, err := s.repo.FindByUsername(*req.Username)
	if err != nil {
		return nil, errors.New("username or password incorrect")
	}

	if user.Status == 0 {
		return nil, errors.New("user is inactive")
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(*req.Password),
	)
	if err != nil {
		return nil, errors.New("username or password incorrect")
	}

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims{
			"user_id":  user.IDUser,
			"username": user.Username,
			"role":     user.Role,
			"exp":      time.Now().Add(8 * time.Hour).Unix(),
			"iat":      time.Now().Unix(),
		},
	)

	tokenString, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return nil, errors.New("failed to create token")
	}

	userResponse := toUserResponse(*user)

	return &dto.LoginResponse{
		Token: &tokenString,
		User:  &userResponse,
	}, nil
}

func (s *userService) UpdateProfile(id uint, req *dto.UpdateProfileRequest) (*dto.UserResponse, error) {

	user, err := s.repo.FindByID(id)

	if err != nil {
		return nil, errors.New("profile not found")
	}

	if req.NamaUser != "" {
		user.NamaUser = req.NamaUser
	}

	if req.NoTelp != "" {
		user.NoTelp = req.NoTelp
	}

	if req.Username != "" {
		user.Username = req.Username
	}

	if err := s.repo.Update(user); err != nil {
		return nil, err
	}

	res := toUserResponse(*user)

	return &res, nil
}

func (s *userService) UpdateProfilePhoto(id uint, fotoPath string) (*dto.UserResponse, error) {

	user, err := s.repo.FindByID(id)

	if err != nil {
		return nil, errors.New("profile not found")
	}

	oldPath := user.Foto
	user.Foto = fotoPath

	if err := s.repo.Update(user); err != nil {
		return nil, err
	}

	if oldPath != "" {
		os.Remove(oldPath)
	}

	res := toUserResponse(*user)
	return &res, nil
}

func toUserResponse(user model.User) dto.UserResponse {
	return dto.UserResponse{
		IDUser:   user.IDUser,
		NamaUser: user.NamaUser,
		NoTelp:   user.NoTelp,
		Role:     user.Role,
		Username: user.Username,
		Status:   intToStatus(user.Status),
		Foto:     user.Foto,
	}
}

func statusToInt(status string) int {
	if status == "aktif" {
		return 1
	}
	return 0
}

func intToStatus(status int) string {
	if status == 1 {
		return "aktif"
	}
	return "nonaktif"
}

func normalizeStatusFilter(status string) string {
	switch status {
	case "aktif":
		return "1"
	case "nonaktif":
		return "0"
	default:
		return status
	}
}
