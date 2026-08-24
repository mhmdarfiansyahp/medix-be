package service

import (
	"errors"
	"math"
	"medix-be/internal/user/model/dto"
	model "medix-be/internal/user/model/entities"
	"medix-be/internal/user/repository"

	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	CreateUser(req dto.CreateUserRequest) (*dto.UserResponse, error)
	GetAllUsers(page int, limit int, search string, role string, status string) (*dto.UserListResponse, error)
	GetUserByID(id uint) (*dto.UserResponse, error)
	UpdateUser(id uint, req dto.UpdateUserRequest) (*dto.UserResponse, error)
	DeleteUser(id uint) error
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) CreateUser(req dto.CreateUserRequest) (*dto.UserResponse, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("gagal memproses password")
	}

	user := model.User{
		NamaUser: req.NamaUser,
		NoTelp:   req.NoTelp,
		Role:     req.Role,
		Username: req.Username,
		Password: string(hashedPassword),
		Status:   req.Status,
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
		status,
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
		return nil, errors.New("user tidak ditemukan")
	}

	res := toUserResponse(*user)
	return &res, nil
}

func (s *userService) UpdateUser(id uint, req dto.UpdateUserRequest) (*dto.UserResponse, error) {

	user, err := s.repo.FindByID(id)

	if err != nil {
		return nil, errors.New("user tidak ditemukan")
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
			return nil, errors.New("gagal memproses password baru")
		}
		user.Password = string(hashedPassword)
	}
	if req.Status != "" {
		user.Status = req.Status
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

func toUserResponse(user model.User) dto.UserResponse {
	return dto.UserResponse{
		IDUser:   user.IDUser,
		NamaUser: user.NamaUser,
		NoTelp:   user.NoTelp,
		Role:     user.Role,
		Username: user.Username,
		Status:   user.Status,
		Foto:     user.Foto,
	}
}
