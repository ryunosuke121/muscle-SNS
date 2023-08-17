package repository

import "github.com/ryunosuke121/muscle-SNS/model"

type IUserRepository interface {
	CreateUser(user *model.User) error
}
