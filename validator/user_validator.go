package validator

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/ryunosuke121/muscle-SNS/model"
)

type IUserValidator interface {
	UserValidator(user model.User) error
}

type userValidator struct{}

func NewUserValidator() IUserValidator {
	return &userValidator{}
}

func (uv *userValidator) UserValidator(user model.User) error {
	return validation.ValidateStruct(&user,
		validation.Field(&user.Name,
			validation.Required.Error("名前を入力してください"),
			validation.Length(1, 20).Error("名前は20文字以内で入力してください"),
		),
		validation.Field(&user.Email,
			validation.Required.Error("メールアドレスを入力してください"),
			validation.Length(1, 30).Error("メールアドレスは30文字以内で入力してください"),
			is.Email.Error("メールアドレスの形式が正しくありません"),
		),
		validation.Field(&user.Password,
			validation.Required.Error("パスワードを入力してください"),
			validation.Length(8, 32).Error("パスワードは8文字以上32文字以内で入力してください"),
		),
	)
}
