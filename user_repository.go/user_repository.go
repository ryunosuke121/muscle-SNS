package repository

import (
	"github.com/ryunosuke121/go-rest/model"
)

type IUserRepsitory interface {
	GetUserByEmail(user *model.User, email string) error
	CreateUser(user *model.User) error
}
