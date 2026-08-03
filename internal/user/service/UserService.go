package service

import (
	"errors"
	"medix-be/internal/user/model/dto"
	"medix-be/internal/user/model/entities"
	"medix-be/internal/user/repository"
)

type UserService interface {
	CreateUser(req dto.CreateUserRequest) (*dto.UserResponse, error)
	GetAllUsers() ([]dto.UserResponse, error)
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
	user := model.User{
		NamaUser: req.NamaUser,
		NoTelp:   req.NoTelp,
		Role:     req.Role,
		Username: req.Username,
		Password: req.Password,
		Status:   req.Status,
		Foto:     req.Foto,
	}

	if err := s.repo.Create(&user); err != nil {
		return nil, err
	}

	res := toUserResponse(user)
	return &res, nil
}

func (s *userService) GetAllUsers() ([]dto.UserResponse, error) {
	users, err := s.repo.FindAll()
	if err != nil {
		return nil, err
	}

	var responses []dto.UserResponse
	for _, u := range users {
		responses = append(responses, toUserResponse(u))
	}
	return responses, nil
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
		user.Password = req.Password
	}
	if req.Status != "" {
		user.Status = req.Status
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