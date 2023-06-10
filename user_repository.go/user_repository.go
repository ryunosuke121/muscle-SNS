package repository

import (
	"github.com/ryunosuke121/muscle-SNS/model"
)

type IUserRepsitory interface {
	GetUserByEmail(user *model.User, email string) error
	CreateUser(user *model.User) error
}
