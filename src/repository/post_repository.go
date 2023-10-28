package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"mime/multipart"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/ryunosuke121/muscle-SNS/src/domain"
	"github.com/ryunosuke121/muscle-SNS/utils"
	"gorm.io/gorm"
)

type PostRepository struct {
	db          *gorm.DB
	s3Client    *s3.Client
	redisClient *redis.Client
}

func NewPostRepository(db *gorm.DB, s3Client *s3.Client, redisClient *redis.Client) domain.IPostRepository {
	return &PostRepository{db, s3Client, redisClient}
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
		post.ImageName = pr.getImageUrlByFileName(post.ImageName)
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
		post.ImageName = pr.getImageUrlByFileName(post.ImageName)
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
		Comment:   post.Comment,
		ImageName: post.ImageName,
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
		post.ImageName = pr.getImageUrlByFileName(post.ImageName)
		domainPost := convertToPost(post)
		domainPosts = append(domainPosts, domainPost)
	}
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

// メニュー別のユーザーの総重量を取得する
func (pr *PostRepository) GetUsersTotalWeightByMenuInMonth(ctx context.Context, userIds []domain.UserID, menuId domain.MenuID, year int, month int) (map[domain.UserID]uint, error) {
	exists, err := pr.redisClient.Exists(ctx, fmt.Sprintf("total_weight_menu_%d_%d_%d", menuId, year, month)).Result()
	if err != nil {
		return nil, err
	}

	if exists == 1 {
		jsonData, err := pr.redisClient.Get(ctx, fmt.Sprintf("total_weight_menu_%d_%d_%d", menuId, year, month)).Result()
		if err != nil {
			return nil, err
		}
		var totalWeightMap map[domain.UserID]uint
		if err = json.Unmarshal([]byte(jsonData), &totalWeightMap); err != nil {
			return nil, err
		}
		return totalWeightMap, nil
	}

	var results []struct {
		UserID     string
		TotalCount uint
	}
	result := pr.db.WithContext(ctx).Table("trainings").Select("user_id, sum(weight * times * sets) as total_count").Where("user_id IN ?", userIds).Where("menu_id = ?", menuId).Where("extract(year from created_at) = ?", year).Where("extract(month from created_at) = ?", month).Group("user_id").Scan(&results)
	if result.Error != nil {
		return nil, result.Error
	}

	totalWeightMap := make(map[domain.UserID]uint)
	for _, result := range results {
		totalWeightMap[domain.UserID(result.UserID)] = result.TotalCount
	}

	go func(ctx context.Context, twm map[domain.UserID]uint, menuId domain.MenuID, year int, month int) {
		jsonData, err := json.Marshal(twm)
		if err != nil {
			log.Println(err)
			return
		}
		err = pr.redisClient.Set(ctx, fmt.Sprintf("total_weight_menu_%d_%d_%d", menuId, year, month), jsonData, time.Duration(600)).Err()
		if err != nil {
			log.Println(err)
			return
		}
	}(ctx, totalWeightMap, menuId, year, month)

	return totalWeightMap, nil
}

// 月内のユーザーの総重量を取得する
func (pr *PostRepository) GetUsersTotalWeightInMonth(ctx context.Context, userIds []domain.UserID, year int, month int) (map[domain.UserID]uint, error) {
	// キャッシュの存在を確認する
	exists, err := pr.redisClient.Exists(ctx, fmt.Sprintf("total_weight_%d_%d", year, month)).Result()
	if err != nil {
		return nil, err
	}

	//　存在する場合はそのまま返す
	if exists == 1 {
		jsonData, err := pr.redisClient.Get(ctx, fmt.Sprintf("total_weight_%d_%d", year, month)).Result()
		if err != nil {
			return nil, err
		}
		var totalWeightMap map[domain.UserID]uint
		if err = json.Unmarshal([]byte(jsonData), &totalWeightMap); err != nil {
			return nil, err
		}
		return totalWeightMap, nil
	}

	var results []struct {
		UserID     string
		TotalCount uint
	}
	result := pr.db.WithContext(ctx).Table("trainings").Select("user_id, sum(weight * times * sets) as total_count").Where("user_id IN ?", userIds).Where("extract(year from created_at) = ?", year).Where("extract(month from created_at) = ?", month).Group("user_id").Scan(&results)
	if result.Error != nil {
		return nil, result.Error
	}

	totalWeightMap := make(map[domain.UserID]uint)
	for _, result := range results {
		totalWeightMap[domain.UserID(result.UserID)] = result.TotalCount
	}

	// Redisに保存
	go func(ctx context.Context, twm map[domain.UserID]uint, year int, month int) {
		jsonData, err := json.Marshal(twm)
		if err != nil {
			log.Println(err)
			return
		}
		err = pr.redisClient.Set(ctx, fmt.Sprintf("total_weight_%d_%d", year, month), jsonData, time.Duration(600)).Err()
		if err != nil {
			log.Println(err)
			return
		}
	}(ctx, totalWeightMap, year, month)

	return totalWeightMap, nil
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

func (pr *PostRepository) getImageUrlByFileName(fileName string) string {
	return fmt.Sprintf("https://%s.s3-%s.amazonaws.com/%s", os.Getenv("BUCKET_NAME"), os.Getenv("AWS_REGION"), fileName)
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
		ImageName: post.ImageName,
	}
	return &domainPost
}
