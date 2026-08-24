package repository

import (
	"gorm.io/gorm"
	"medix-be/internal/user/model/entities"
)

type UserRepository interface {
	Create(user *model.User) error
	FindAll(page int, limit int, search string, role string, status string) ([]*model.User, int64, error)
	FindByID(id uint) (*model.User, error)
	Update(user *model.User) error
	Delete(id uint) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(user *model.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) FindAll(page int, limit int, search string, role string, status string) ([]*model.User, int64, error) {

	var users []*model.User
	var total int64
	query := r.db.Model(&model.User{})

	if search != "" {
		searchValue := "%" + search + "%"

		query = query.Where(
			"nama_user ILIKE ? OR username ILIKE ?",
			searchValue,
			searchValue,
		)
	}

	if role != "" {
		query = query.Where("role = ?", role)
	}

	if status != "" {
		query = query.Where("status = ?", status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit

	err := query.
		Offset(offset).
		Limit(limit).
		Find(&users).Error

	return users, total, err
}

func (r *userRepository) FindByID(id uint) (*model.User, error) {
	var user model.User
	err := r.db.First(&user, id).Error
	return &user, err
}

func (r *userRepository) Update(user *model.User) error {
	return r.db.Save(user).Error
}

func (r *userRepository) Delete(id uint) error {
	return r.db.Delete(&model.User{}, id).Error
}
