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
	"github.com/ryunosuke121/muscle-SNS/src/domain"
	"github.com/ryunosuke121/muscle-SNS/utils"
	"gorm.io/gorm"
)

type userRepository struct {
	db       *gorm.DB
	s3Client *s3.Client
}

func NewUserRepository(db *gorm.DB, s3Client *s3.Client) domain.IUserRepository {
	return &userRepository{db, s3Client}
}

// ユーザーを作成する
func (ur *userRepository) CreateUser(ctx context.Context, user *domain.User) error {
	newUser := User{
		ID:    user.ID.String(),
		Name:  user.Name.String(),
		Email: user.Email,
	}

	result := ur.db.WithContext(ctx).Create(&newUser)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			return domain.ErrUserAlreadyExists
		}
		return result.Error
	}
	return nil
}

// IDのリストからユーザーを取得する
func (ur *userRepository) GetUsersByIds(ctx context.Context, userIds []domain.UserID) ([]*domain.User, error) {
	var users []*User
	result := ur.db.WithContext(ctx).Where(userIds).Joins("UserGroup").Find(&users)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, domain.ErrUserNotFound
		}
		return nil, result.Error
	}

	var domainUsers []*domain.User
	for _, user := range users {
		avatarUrl := ur.getImageUrlByFileName(user.ImageUrl)
		domainUser := domain.User{
			ID:        domain.UserID(user.ID),
			Name:      domain.UserName(user.Name),
			Email:     user.Email,
			AvatarUrl: avatarUrl,
			UserGroup: &domain.UserGroup{
				ID:       domain.UserGroupID(user.UserGroup.ID),
				Name:     user.UserGroup.Name,
				ImageUrl: user.UserGroup.ImageUrl,
			},
		}
		domainUsers = append(domainUsers, &domainUser)
	}
	return domainUsers, nil
}

// メールアドレスからユーザーを取得する
func (ur *userRepository) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user User
	result := ur.db.WithContext(ctx).Where("email = ?", email).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, domain.ErrUserNotFound
		}
		return nil, result.Error
	}

	domainUser := domain.User{
		ID:        domain.UserID(user.ID),
		Name:      domain.UserName(user.Name),
		Email:     user.Email,
		AvatarUrl: "",
		UserGroup: &domain.UserGroup{
			ID:       domain.UserGroupID(user.UserGroup.ID),
			Name:     user.UserGroup.Name,
			ImageUrl: user.UserGroup.ImageUrl,
		},
	}

	return &domainUser, nil
}

// ユーザーの名前を更新する
func (ur *userRepository) ChangeUserName(ctx context.Context, userId domain.UserID, userName domain.UserName) error {
	var user User
	result := ur.db.WithContext(ctx).Model(&user).Where("id = ?", userId).Update("name", userName)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return domain.ErrUserNotFound
		}
		return result.Error
	}
	if result.RowsAffected == 0 {
		return domain.ErrUserInfoNotChanged
	}

	return nil
}

// ユーザーの所属するグループを変更する
func (ur *userRepository) ChangeUserGroup(ctx context.Context, userId domain.UserID, groupId domain.UserGroupID) error {
	var user User
	result := ur.db.WithContext(ctx).Model(&user).Where("id = ?", userId).Update("user_group_id", groupId)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return domain.ErrUserNotFound
		}
		return result.Error
	}
	if result.RowsAffected == 0 {
		return domain.ErrUserInfoNotChanged
	}

	return nil
}

// ユーザーの画像を変更する
func (ur *userRepository) ChangeUserImage(ctx context.Context, userId domain.UserID, file *multipart.FileHeader) error {
	fileName, err := ur.saveUserImage(ctx, userId, file)
	if err != nil {
		return err
	}

	var user User
	result := ur.db.WithContext(ctx).Model(&user).Where("id = ?", userId).Update("image_url", fileName)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return domain.ErrUserNotFound
		}
		return result.Error
	}
	if result.RowsAffected == 0 {
		return domain.ErrUserInfoNotChanged
	}

	return nil
}

func (ur *userRepository) GetUsersInGroup(ctx context.Context, groupId domain.UserGroupID) ([]*domain.User, error) {
	var users []*User
	result := ur.db.WithContext(ctx).Where("user_group_id = ?", groupId).Find(&users)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, domain.ErrUserNotFound
		}
		return nil, result.Error
	}

	var domainUsers []*domain.User
	for _, user := range users {
		avatarUrl := ur.getImageUrlByFileName(user.ImageUrl)
		domainUser := domain.User{
			ID:        domain.UserID(user.ID),
			Name:      domain.UserName(user.Name),
			Email:     user.Email,
			AvatarUrl: avatarUrl,
			UserGroup: &domain.UserGroup{
				ID:       domain.UserGroupID(user.UserGroup.ID),
				Name:     user.UserGroup.Name,
				ImageUrl: user.UserGroup.ImageUrl,
			},
		}
		domainUsers = append(domainUsers, &domainUser)
	}
	return domainUsers, nil
}

func (ur *userRepository) saveUserImage(ctx context.Context, userId domain.UserID, file *multipart.FileHeader) (fileName string, err error) {
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	// s3にアップロード
	u, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}
	contentType, err := utils.InspectFileMimeType(file)
	if err != nil {
		return "", err
	}
	extension := utils.GetFileExtension(contentType)

	fileName = fmt.Sprintf("user_image/%s%s", u.String(), extension)

	param := s3.PutObjectInput{
		Bucket: aws.String(os.Getenv("BUCKET_NAME")),
		Key:    aws.String(fileName),
		Body:   src,
	}
	_, err = ur.s3Client.PutObject(ctx, &param)
	if err != nil {
		return "", err
	}

	return fileName, nil
}

// ファイル名からs3に保存されている画像のURLを取得する
func (ur *userRepository) getImageUrlByFileName(fileName string) string {
	return fmt.Sprintf("https://%s.s3-%s.amazonaws.com/%s", os.Getenv("BUCKET_NAME"), os.Getenv("AWS_REGION"), fileName)
}
