package repository

import (
	"context"
	"fmt"
	"log"
	"mime/multipart"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
	"github.com/ryunosuke121/muscle-SNS/src/domain"
	"github.com/ryunosuke121/muscle-SNS/utils"
	"gorm.io/gorm"
)

type PostRepository struct {
	db              *gorm.DB
	s3Client        *s3.Client
	s3PresignClient *s3.PresignClient
}

func NewPostRepository(db *gorm.DB, s3Client *s3.Client, s3PresignClient *s3.PresignClient) domain.IPostRepository {
	return &PostRepository{db, s3Client, s3PresignClient}
}

// 投稿を取得する
func (pr *PostRepository) GetPostsByIds(ctx context.Context, ids []domain.PostID) ([]*domain.Post, error) {
	var posts []*Post
	result := pr.db.WithContext(ctx).Where(ids).Joins("Training").Joins("Training.Menu").Find(&posts)
	if result.Error != nil {
		return nil, result.Error
	}

	var domainPosts []*domain.Post
	for _, post := range posts {
		domainPost := convertToPost(post)
		domainPosts = append(domainPosts, domainPost)
	}
	return domainPosts, nil
}

// ユーザーの投稿を取得する
func (pr *PostRepository) GetUserPosts(ctx context.Context, id domain.UserID) ([]*domain.Post, error) {
	var posts []*Post
	result := pr.db.WithContext(ctx).Where("posts.user_id = ?", id).Joins("Training").Joins("Training.Menu").Find(&posts)
	if result.Error != nil {
		return nil, result.Error
	}

	var domainPosts []*domain.Post
	for _, post := range posts {
		domainPost := convertToPost(post)
		domainPosts = append(domainPosts, domainPost)
	}
	return domainPosts, nil
}

// 投稿を作成する
func (pr *PostRepository) CreatePost(ctx context.Context, post *domain.Post) (*domain.Post, error) {
	newPost := Post{
		UserID: post.UserID.String(),
		Training: &Training{
			UserID: post.Training.UserID.String(),
			MenuID: uint(post.Training.Menu.ID),
			Times:  post.Training.Times,
			Weight: post.Training.Weight,
			Sets:   post.Training.Sets,
		},
		Comment:  post.Comment,
		ImageUrl: post.ImageUrl,
	}

	result := pr.db.WithContext(ctx).Create(&newPost)
	if result.Error != nil {
		return nil, result.Error
	}

	createdPost, error := pr.GetPostsByIds(ctx, []domain.PostID{domain.PostID(newPost.ID)})
	if error != nil {
		return nil, error
	}

	return createdPost[0], nil
}

// 投稿を削除する
func (pr *PostRepository) DeletePost(ctx context.Context, id domain.PostID) error {
	result := pr.db.WithContext(ctx).Delete(&Post{}, id)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// グループの投稿を取得する\
func (pr *PostRepository) GetGroupPosts(ctx context.Context, id domain.UserGroupID) ([]*domain.Post, error) {
	var posts []*Post
	result := pr.db.WithContext(ctx).Joins("LEFT JOIN users ON users.id = posts.user_id").Where("users.user_group_id = ?", id).Joins("Training").Joins("Training.Menu").Find(&posts)
	if result.Error != nil {
		return nil, result.Error
	}

	var domainPosts []*domain.Post
	for _, post := range posts {
		domainPost := convertToPost(post)
		domainPosts = append(domainPosts, domainPost)
	}
	log.Print(len(domainPosts))
	return domainPosts, nil
}

// トレーニングを取得する
func (pr *PostRepository) GetTrainingsByIds(ctx context.Context, ids []domain.TrainingID) ([]*domain.Training, error) {
	var trainings []*Training
	result := pr.db.WithContext(ctx).Where(ids).Find(&trainings)
	if result.Error != nil {
		return nil, result.Error
	}

	var domainTrainings []*domain.Training
	for _, training := range trainings {
		domainTraining := domain.Training{
			ID:     domain.TrainingID(training.ID),
			UserID: domain.UserID(training.UserID),
			Menu: &domain.Menu{
				ID:   domain.MenuID(training.MenuID),
				Name: training.Menu.Name,
			},
			Times:     training.Times,
			Weight:    training.Weight,
			Sets:      training.Sets,
			CreatedAt: training.CreatedAt,
		}
		domainTrainings = append(domainTrainings, &domainTraining)
	}
	return domainTrainings, nil
}

// ユーザーのトレーニングを取得する
func (pr *PostRepository) GetUserTrainings(ctx context.Context, id domain.UserID) ([]*domain.Training, error) {
	var trainings []*Training
	result := pr.db.WithContext(ctx).Where("user_id = ?", id).Find(&trainings)
	if result.Error != nil {
		return nil, result.Error
	}

	var domainTrainings []*domain.Training
	for _, training := range trainings {
		domainTraining := domain.Training{
			ID:     domain.TrainingID(training.ID),
			UserID: domain.UserID(training.UserID),
			Menu: &domain.Menu{
				ID:   domain.MenuID(training.MenuID),
				Name: training.Menu.Name,
			},
			Times:     training.Times,
			Weight:    training.Weight,
			Sets:      training.Sets,
			CreatedAt: training.CreatedAt,
		}
		domainTrainings = append(domainTrainings, &domainTraining)
	}

	return domainTrainings, nil
}

// TODO: コンテキストを持たせてキャンセルできるようにする
func (pr *PostRepository) SavePostImage(ctx context.Context, file *multipart.FileHeader) (fileName string, err error) {
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

	fileName = fmt.Sprintf("post_image/%s%s", u.String(), extension)

	param := s3.PutObjectInput{
		Bucket: aws.String(os.Getenv("BUCKET_NAME")),
		Key:    aws.String(fileName),
		Body:   src,
	}
	_, err = pr.s3Client.PutObject(ctx, &param)
	if err != nil {
		return "", err
	}

	return fileName, nil
}

func convertToPost(post *Post) *domain.Post {
	domainTraining := domain.Training{
		ID:     domain.TrainingID(post.Training.ID),
		UserID: domain.UserID(post.Training.UserID),
		Menu: &domain.Menu{
			ID:   domain.MenuID(post.Training.MenuID),
			Name: post.Training.Menu.Name,
		},
		Times:     post.Training.Times,
		Weight:    post.Training.Weight,
		Sets:      post.Training.Sets,
		CreatedAt: post.Training.CreatedAt,
	}

	domainPost := domain.Post{
		ID:        domain.PostID(post.ID),
		UserID:    domain.UserID(post.UserID),
		Training:  &domainTraining,
		Comment:   post.Comment,
		CreatedAt: post.CreatedAt,
		ImageUrl:  post.ImageUrl,
	}
	return &domainPost
}
