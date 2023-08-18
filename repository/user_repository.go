package repository

import (
	"github.com/ryunosuke121/muscle-SNS/model"
	"gorm.io/gorm"
)

type IUserRepository interface {
	CreateUser(user *model.User) error
	GetUserById(user *model.User, userId uint) error
	UpdateUser(user *model.User, userId uint) error
	GetUserByEmail(user *model.User, email string) error
	GetUserImageUrlById(userId uint) (string, error)
	SetUserImage(user *model.User, file []byte) (string, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) IUserRepository {
	return &userRepository{db}
}

func (ur *userRepository) CreateUser(user *model.User) error {
	return nil
}

func (ur *userRepository) GetUserById(user *model.User, userId uint) error {
	return nil
}

func (ur *userRepository) GetUserByEmail(user *model.User, email string) error {
	return nil
}

func (ur *userRepository) UpdateUser(user *model.User, id uint) error {
	return nil
}

func (ur *userRepository) GetUserImageUrlById(userId uint) (string, error) {
	return "", nil
}

func (ur *userRepository) SetUserImage(user *model.User, file []byte) (string, error) {

	return "", nil
}
