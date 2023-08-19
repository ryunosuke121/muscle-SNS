package repository

import (
	"context"
	"fmt"
	"mime/multipart"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
	"github.com/ryunosuke121/muscle-SNS/model"
	"github.com/ryunosuke121/muscle-SNS/utils"
	"gorm.io/gorm"
)

type IUserRepository interface {
	CreateUser(user *model.User) error
	GetUserById(user *model.User, userId uint) error
	UpdateUserName(user *model.User, userId uint) error
	GetUserByEmail(user *model.User, email string) error
	GetUserImageUrlById(userId uint) (string, error)
	SetUserImage(user *model.User, file *multipart.FileHeader) error
}

type userRepository struct {
	db       *gorm.DB
	s3Client *s3.Client
}

func NewUserRepository(db *gorm.DB, s3Client *s3.Client) IUserRepository {
	return &userRepository{db, s3Client}
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

func (ur *userRepository) UpdateUserName(user *model.User, userId uint) error {

	return nil
}

func (ur *userRepository) UpdateUserTrainingGroup(user *model.User, groupId uint) error {
	return nil
}

func (ur *userRepository) GetUserImageUrlById(userId uint) (string, error) {
	return "", nil
}

func (ur *userRepository) SetUserImage(user *model.User, file *multipart.FileHeader) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	// s3にアップロード
	u, err := uuid.NewRandom()
	if err != nil {
		return err
	}
	contentType, err := utils.InspectFileMimeType(file)
	if err != nil {
		return err
	}
	extension := getFileExtension(contentType)

	fileName := fmt.Sprintf("user_image/%s%s", u.String(), extension)

	param := s3.PutObjectInput{
		Bucket: aws.String(os.Getenv("BUCKET_NAME")),
		Key:    aws.String(fileName),
		Body:   src,
	}
	_, err = ur.s3Client.PutObject(context.TODO(), &param)
	if err != nil {
		return err
	}

	// DBに保存
	user.ImageUrl = fileName
	result := ur.db.Model(user).Update("image_url", fileName)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func getFileExtension(mime_type string) string {
	switch mime_type {
	case "image/jpeg":
		return ".jpg"
	case "image/png":
		return ".png"
	default:
		return ""
	}
}
