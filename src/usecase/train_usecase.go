package usecase

import (
	"github.com/ryunosuke121/muscle-SNS/src/model"
	"github.com/ryunosuke121/muscle-SNS/src/repository"
)

type ITrainUsecase interface {
	GetTrainingById(id uint) (*model.Training, error)
	GetUserTrainings(id uint) (*[]model.Training, error)
	CreatePost(post *model.Post) error
	GetUserPosts(id uint) (*[]model.Post, error)
}

type TrainUsecase struct {
	tr repository.ITrainRepository
}

func NewTrainUsecase(tr repository.ITrainRepository) ITrainUsecase {
	return &TrainUsecase{tr: tr}
}

func (tu *TrainUsecase) GetTrainingById(id uint) (*model.Training, error) {
	return tu.tr.GetTrainingById(id)
}

func (tu *TrainUsecase) GetUserTrainings(id uint) (*[]model.Training, error) {
	return tu.tr.GetUserTrainings(id)
}

func (tu *TrainUsecase) CreatePost(post *model.Post) error {
	return tu.tr.CreatePost(post)
}

func (tu *TrainUsecase) GetUserPosts(id uint) (*[]model.Post, error) {
	return tu.tr.GetUserPosts(id)
}
