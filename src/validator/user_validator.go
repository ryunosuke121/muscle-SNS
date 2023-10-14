package validator

import (
	"errors"
	"io"
	"mime/multipart"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/ryunosuke121/muscle-SNS/src/model"
	"github.com/ryunosuke121/muscle-SNS/utils"
)

type IUserValidator interface {
	UserValidator(user model.User) error
	UserImageValidator(file *multipart.FileHeader) error
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
	)
}

func (uv *userValidator) UserImageValidator(file *multipart.FileHeader) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	// ファイルサイズの検証
	const maxFileSize = 10 << 20
	fileBytes, err := io.ReadAll(src)
	if err != nil {
		return errors.New("ファイルを読み込めません")
	}
	if len(fileBytes) > maxFileSize {
		return errors.New("ファイルサイズが大きすぎます")
	}

	// MIMEタイプの検証
	contentType, err := utils.InspectFileMimeType(file)
	if err != nil {
		return err
	}
	if contentType != "image/jpeg" && contentType != "image/png" {
		return errors.New("ファイル形式が正しくありません")
	}

	return nil
}
