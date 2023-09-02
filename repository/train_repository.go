package repository

import (
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/ryunosuke121/muscle-SNS/model"
	"gorm.io/gorm"
)

type ITrainRepository interface {
	GetTrainingById(id uint) (*model.Training, error)
	GetUserTrainings(id uint) (*[]model.Training, error)
	CreatePost(post *model.Post) error
	GetUserPosts(id uint) (*[]model.Post, error)
}

type TrainRepository struct {
	db              *gorm.DB
	s3Client        *s3.Client
	s3PresignClient *s3.PresignClient
}

func NewTrainRepository(db *gorm.DB, s3Client *s3.Client, s3PresignClient *s3.PresignClient) ITrainRepository {
	return &TrainRepository{db, s3Client, s3PresignClient}
}

func (tr *TrainRepository) GetTrainingById(id uint) (*model.Training, error) {
	training := new(model.Training)
	result := tr.db.First(&training, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return training, nil
}

func (tr *TrainRepository) GetUserTrainings(id uint) (*[]model.Training, error) {
	trainings := new([]model.Training)
	result := tr.db.Where("user_id = ?", id).Find(&trainings)
	if result.Error != nil {
		return nil, result.Error
	}
	return trainings, nil
}

func (tr *TrainRepository) CreatePost(post *model.Post) error {
	result := tr.db.Create(&post)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (tr *TrainRepository) GetUserPosts(id uint) (*[]model.Post, error) {
	posts := new([]model.Post)
	result := tr.db.Where("user_id = ?", id).Find(&posts)
	if result.Error != nil {
		return nil, result.Error
	}
	return posts, nil
}
