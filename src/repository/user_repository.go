package repository

import (
	"context"
	"fmt"
	"log"
	"mime/multipart"
	"os"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
	"github.com/ryunosuke121/muscle-SNS/src/domain"
	"github.com/ryunosuke121/muscle-SNS/utils"
	"gorm.io/gorm"
)

type userRepository struct {
	db              *gorm.DB
	s3Client        *s3.Client
	s3PresignClient *s3.PresignClient
}

func NewUserRepository(db *gorm.DB, s3Client *s3.Client, s3PresignClient *s3.PresignClient) domain.IUserRepository {
	return &userRepository{db, s3Client, s3PresignClient}
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
		return result.Error
	}
	return nil
}

// IDのリストからユーザーを取得する
func (ur *userRepository) GetUsersByIds(ctx context.Context, userIds []domain.UserID) ([]*domain.User, error) {
	var users []*User
	result := ur.db.WithContext(ctx).Where(userIds).Joins("UserGroup").Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}

	url, err := ur.GetUserImageUrlsByIds(ctx, userIds)
	if err != nil {
		return nil, err
	}

	var domainUsers []*domain.User
	for _, user := range users {
		domainUser := domain.User{
			ID:        domain.UserID(user.ID),
			Name:      domain.UserName(user.Name),
			Email:     user.Email,
			AvatarUrl: url[domain.UserID(user.ID)],
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
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("user not found or name is same")
	}

	return nil
}

// ユーザーの所属するグループを変更する
func (ur *userRepository) ChangeUserGroup(ctx context.Context, userId domain.UserID, groupId domain.UserGroupID) error {
	var user User
	result := ur.db.WithContext(ctx).Model(&user).Where("id = ?", userId).Update("training_group_id", groupId)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("user not found or group is same")
	}

	return nil
}

// IDのリストからユーザーの画像のURLを取得する
func (ur *userRepository) GetUserImageUrlsByIds(ctx context.Context, userIds []domain.UserID) (map[domain.UserID]string, error) {
	var users []*User
	result := ur.db.WithContext(ctx).Select("id", "image_url").Find(&users, userIds)
	if result.Error != nil {
		return nil, result.Error
	}

	var fileNames []string
	for _, user := range users {
		fileNames = append(fileNames, user.ImageUrl)
	}

	fileNameUrlMap := ur.getImageUrlByFileName(fileNames)
	log.Print(fileNameUrlMap)
	userIDUrlMap := make(map[domain.UserID]string)
	for _, user := range users {
		userIDUrlMap[domain.UserID(user.ID)] = fileNameUrlMap[user.ImageUrl]
	}
	return userIDUrlMap, nil
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
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
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

type Result struct {
	mu        sync.Mutex
	resultMap map[string]string
}

func (r *Result) Set(key string, value string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.resultMap[key] = value
}

// 並行処理で画像のURLを取得する
func (ur *userRepository) getImageUrlByFileName(fileNames []string) map[string]string {
	var wg sync.WaitGroup
	r := Result{resultMap: make(map[string]string)}

	for _, fileName := range fileNames {
		wg.Add(1)
		go func(fileName string) {
			if fileName == "" {
				r.Set(fileName, "")
				wg.Done()
				return
			}
			// s3から画像を取得
			param := s3.GetObjectInput{
				Bucket: aws.String(os.Getenv("BUCKET_NAME")),
				Key:    aws.String(fileName),
			}
			res, err := ur.s3PresignClient.PresignGetObject(context.Background(), &param)
			if err != nil {
				return
			}
			r.Set(fileName, res.URL)
			wg.Done()
		}(fileName)
	}
	wg.Wait()
	return r.resultMap
}
