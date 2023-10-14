package repository

import (
	"context"
	"errors"
	"fmt"
	"mime/multipart"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
	"github.com/ryunosuke121/muscle-SNS/src/model"
	"github.com/ryunosuke121/muscle-SNS/utils"
	"gorm.io/gorm"
)

type IUserRepository interface {
	CreateUser(user *model.User) error
	GetUserById(user *model.User, userId string) error
	UpdateUserName(user *model.User, userId string, userName string) error
	UpdateUserTrainingGroup(user *model.User, userId string, groupId uint) error
	GetUserByEmail(user *model.User, email string) error
	GetUserImageUrlById(userId string) (string, error)
	SetUserImage(user *model.User, userId string, file *multipart.FileHeader) error
}

type userRepository struct {
	db              *gorm.DB
	s3Client        *s3.Client
	s3PresignClient *s3.PresignClient
}

func NewUserRepository(db *gorm.DB, s3Client *s3.Client, s3PresignClient *s3.PresignClient) IUserRepository {
	return &userRepository{db, s3Client, s3PresignClient}
}

func (ur *userRepository) CreateUser(user *model.User) error {
	result := ur.db.Create(&user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (ur *userRepository) GetUserById(user *model.User, userId string) error {
	result := ur.db.First(&user, userId)
	if result.Error != nil {
		return result.Error
	}
	url, err := ur.GetUserImageUrlById(userId)
	if err != nil {
		return err
	}
	user.ImageUrl = url
	return nil
}

func (ur *userRepository) GetUserByEmail(user *model.User, email string) error {
	result := ur.db.Where("email = ?", email).First(&user)
	if result.Error != nil {
		return result.Error
	}
	url, err := ur.GetUserImageUrlById(user.ID)
	if err != nil {
		return err
	}
	user.ImageUrl = url
	return nil
}

func (ur *userRepository) UpdateUserName(user *model.User, userId string, userName string) error {
	result := ur.db.Model(user).Where("id = ?", userId).Update("name", userName)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("user not found or name is same")
	}

	ur.db.First(&user, userId)
	url, err := ur.GetUserImageUrlById(user.ID)
	if err != nil {
		return err
	}
	user.ImageUrl = url
	return nil
}

func (ur *userRepository) UpdateUserTrainingGroup(user *model.User, userId string, groupId uint) error {
	result := ur.db.Model(user).Where("id = ?", userId).Update("training_group_id", groupId)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("user not found or group is same")
	}

	ur.db.First(&user, userId)
	url, err := ur.GetUserImageUrlById(user.ID)
	if err != nil {
		return err
	}
	user.ImageUrl = url
	return nil
}

func (ur *userRepository) GetUserImageUrlById(userId string) (string, error) {
	user := model.User{}
	result := ur.db.First(&user, userId)
	if result.Error != nil {
		return "", result.Error
	}

	if user.ImageUrl == "" {
		return "", nil
	}

	// s3から画像を取得
	param := s3.GetObjectInput{
		Bucket: aws.String(os.Getenv("BUCKET_NAME")),
		Key:    aws.String(user.ImageUrl),
	}
	res, err := ur.s3PresignClient.PresignGetObject(context.Background(), &param)
	if err != nil {
		return "", err
	}
	return res.URL, nil
}

func (ur *userRepository) SetUserImage(user *model.User, userId string, file *multipart.FileHeader) error {
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

	// DBにファイル名を保存
	user.ImageUrl = fileName
	result := ur.db.Model(user).Where("id = ?", userId).Update("image_url", fileName)
	if result.Error != nil {
		return result.Error
	}

	param := s3.PutObjectInput{
		Bucket: aws.String(os.Getenv("BUCKET_NAME")),
		Key:    aws.String(fileName),
		Body:   src,
	}
	_, err = ur.s3Client.PutObject(context.TODO(), &param)
	if err != nil {
		// アップロードに失敗したらDBのファイル名を空にする
		result = ur.db.Model(user).Where("id = ?", userId).Update("image_url", "")
		if result.Error != nil {
			return errors.New("file upload failed and failed to fix DB")
		}
		return err
	}

	ur.db.First(&user, userId)
	url, err := ur.GetUserImageUrlById(user.ID)
	if err != nil {
		return err
	}
	user.ImageUrl = url

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
