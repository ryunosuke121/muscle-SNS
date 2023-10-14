package usecase

import (
	"mime/multipart"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/ryunosuke121/muscle-SNS/src/model"
	"github.com/ryunosuke121/muscle-SNS/src/repository"
	"github.com/ryunosuke121/muscle-SNS/src/validator"
)

type IUserUseCase interface {
	SignUp(user model.User) (model.UserResponse, error)
	Login(user model.User) (string, error)
	GetUserImageUrlById(userId string) (string, error)
	GetUserById(user *model.User, userId string) (model.UserResponse, error)
	UpdateUserName(name string, userId string) (model.UserResponse, error)
	UpdateUserTrainingGroup(groupId uint, userId string) (model.UserResponse, error)
	UpdateUserImage(file *multipart.FileHeader, userId string) (model.UserResponse, error)
}

type userUsecase struct {
	ur repository.IUserRepository
	uv validator.IUserValidator
}

func NewUserUsecase(ur repository.IUserRepository, uv validator.IUserValidator) IUserUseCase {
	return &userUsecase{ur, uv}
}

func (uu *userUsecase) SignUp(user model.User) (model.UserResponse, error) {
	newUser := model.User{
		ID:              user.ID,
		Name:            user.Name,
		Email:           user.Email,
		TrainingGroupID: user.TrainingGroupID,
	}

	if err := uu.ur.CreateUser(&newUser); err != nil {
		return model.UserResponse{}, err
	}

	resUser := model.UserResponse{
		ID:              newUser.ID,
		Name:            newUser.Name,
		Email:           newUser.Email,
		TrainingGroupID: newUser.TrainingGroupID,
	}

	return resUser, nil
}

func (uu *userUsecase) Login(user model.User) (string, error) {
	if err := uu.uv.UserValidator(user); err != nil {
		return "", err
	}

	storedUser := model.User{}
	if err := uu.ur.GetUserByEmail(&storedUser, user.Email); err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": storedUser.ID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET_KEY")))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (uu *userUsecase) GetUserImageUrlById(userId string) (string, error) {
	imgUrl, err := uu.ur.GetUserImageUrlById(userId)
	if err != nil {
		return "", err
	}
	return imgUrl, nil
}

func (uu *userUsecase) GetUserById(user *model.User, userId string) (model.UserResponse, error) {
	if err := uu.ur.GetUserById(user, userId); err != nil {
		return model.UserResponse{}, err
	}

	resUser := model.UserResponse{
		ID:              user.ID,
		Name:            user.Name,
		Email:           user.Email,
		TrainingGroupID: user.TrainingGroupID,
		ImageUrl:        user.ImageUrl,
	}

	return resUser, nil
}

func (uu *userUsecase) UpdateUserName(name string, userId string) (model.UserResponse, error) {
	user := &model.User{}
	if err := uu.ur.UpdateUserName(user, userId, name); err != nil {
		return model.UserResponse{}, err
	}

	resUser := model.UserResponse{
		ID:              user.ID,
		Name:            user.Name,
		Email:           user.Email,
		TrainingGroupID: user.TrainingGroupID,
		ImageUrl:        user.ImageUrl,
	}

	return resUser, nil
}

func (uu *userUsecase) UpdateUserTrainingGroup(groupId uint, userId string) (model.UserResponse, error) {
	user := &model.User{}

	if err := uu.ur.UpdateUserTrainingGroup(user, userId, groupId); err != nil {
		return model.UserResponse{}, err
	}

	resUser := model.UserResponse{
		ID:              user.ID,
		Name:            user.Name,
		Email:           user.Email,
		TrainingGroupID: user.TrainingGroupID,
		ImageUrl:        user.ImageUrl,
	}

	return resUser, nil
}

func (uu *userUsecase) UpdateUserImage(file *multipart.FileHeader, userId string) (model.UserResponse, error) {
	if err := uu.uv.UserImageValidator(file); err != nil {
		return model.UserResponse{}, err
	}

	user := &model.User{}
	if err := uu.ur.SetUserImage(user, userId, file); err != nil {
		return model.UserResponse{}, err
	}

	resUser := model.UserResponse{
		ID:              user.ID,
		Name:            user.Name,
		Email:           user.Email,
		TrainingGroupID: user.TrainingGroupID,
		ImageUrl:        user.ImageUrl,
	}

	return resUser, nil
}
