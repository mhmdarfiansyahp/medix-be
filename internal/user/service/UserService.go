package service

import (
	"crypto/rand"
	"encoding/hex"
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
	RefreshToken(req *dto.RefreshRequest) (*dto.RefreshResponse, error)
	Logout(req *dto.LogoutRequest) error

	GetAllUsers(page int, limit int, search string, role string, status string) (*dto.UserListResponse, error)
	GetUserByID(id uint) (*dto.UserResponse, error)
	UpdateUser(id uint, req dto.UpdateUserRequest) (*dto.UserResponse, error)
	DeleteUser(id uint) error
	GetProfile(id uint) (*dto.UserResponse, error)
	UpdateProfile(id uint, req *dto.UpdateProfileRequest) (*dto.UserResponse, error)
	UpdateProfilePhoto(id uint, fotoPath string) (*dto.UserResponse, error)
	GetProfileByUsername(username string) (*dto.UserResponse, error)
}

type userService struct {
	repo        repository.UserRepository
	sessionRepo repository.SessionRepository
	jwtSecret   string
}

func NewUserService(repo repository.UserRepository, sessionRepo repository.SessionRepository, jwtSecret string) UserService {
	return &userService{repo: repo, sessionRepo: sessionRepo, jwtSecret: jwtSecret}
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

	tokenString, err := s.generateAccessToken(user)
	if err != nil {
		return nil, err
	}

	refreshTokenString, err := generateRefreshToken()
	if err != nil {
		return nil, errors.New("failed to create refresh token")
	}

	session := &model.Session{
		IDUser:       user.IDUser,
		RefreshToken: refreshTokenString,
		ExpiresAt:    time.Now().Add(30 * 24 * time.Hour),
	}
	if err := s.sessionRepo.Create(session); err != nil {
		return nil, errors.New("failed to create session")
	}

	userResponse := toUserResponse(*user)

	return &dto.LoginResponse{
		Token:        &tokenString,
		RefreshToken: &refreshTokenString,
		User:         &userResponse,
	}, nil
}

func (s *userService) RefreshToken(req *dto.RefreshRequest) (*dto.RefreshResponse, error) {
	session, err := s.sessionRepo.FindByRefreshToken(req.RefreshToken)
	if err != nil {
		return nil, errors.New("invalid or expired refresh token")
	}

	if err := s.sessionRepo.RevokeByID(session.IDSession); err != nil {
		return nil, errors.New("failed to revoke old session")
	}

	user, err := s.repo.FindByID(session.IDUser)
	if err != nil {
		return nil, errors.New("user not found")
	}

	if user.Status == 0 {
		return nil, errors.New("user is inactive")
	}

	tokenString, err := s.generateAccessToken(user)
	if err != nil {
		return nil, err
	}

	newRefreshToken, err := generateRefreshToken()
	if err != nil {
		return nil, errors.New("failed to create refresh token")
	}

	newSession := &model.Session{
		IDUser:       user.IDUser,
		RefreshToken: newRefreshToken,
		ExpiresAt:    time.Now().Add(30 * 24 * time.Hour),
	}
	if err := s.sessionRepo.Create(newSession); err != nil {
		return nil, errors.New("failed to create session")
	}

	return &dto.RefreshResponse{
		Token:        &tokenString,
		RefreshToken: &newRefreshToken,
	}, nil
}

func (s *userService) Logout(req *dto.LogoutRequest) error {
	session, err := s.sessionRepo.FindByRefreshToken(req.RefreshToken)
	if err != nil {
		return errors.New("session not found or already revoked")
	}
	return s.sessionRepo.RevokeByID(session.IDSession)
}

func (s *userService) generateAccessToken(user *model.User) (string, error) {
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims{
			"user_id":  user.IDUser,
			"username": user.Username,
			"role":     user.Role,
			"exp":      time.Now().Add(6 * time.Hour).Unix(),
			"iat":      time.Now().Unix(),
		},
	)

	tokenString, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", errors.New("failed to create token")
	}
	return tokenString, nil
}

func generateRefreshToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
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

func (s *userService) GetProfileByUsername(username string) (*dto.UserResponse, error) {

	user, err := s.repo.FindByUsername(username)

	if err != nil {
		return nil, errors.New("profile not found")
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
