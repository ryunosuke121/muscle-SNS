package usecase

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/ryunosuke121/muscle-SNS/model"
	"github.com/ryunosuke121/muscle-SNS/repository"
	"github.com/ryunosuke121/muscle-SNS/validator"
	"golang.org/x/crypto/bcrypt"
)

type IUserUseCase interface {
	SignUp(user model.User) (model.UserResponse, error)
	Login(user model.User) (string, error)
	SetUserImage(user model.User, file []byte) (string, error)
	GetUserImageUrlById(userId uint) (string, error)
}

type userUsecase struct {
	ur repository.IUserRepository
	uv validator.IUserValidator
}

func NewUserUsecase(ur repository.IUserRepository, uv validator.IUserValidator) IUserUseCase {
	return &userUsecase{ur, uv}
}

func (uu *userUsecase) SignUp(user model.User) (model.UserResponse, error) {
	if err := uu.uv.UserValidator(user); err != nil {
		return model.UserResponse{}, err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	if err != nil {
		return model.UserResponse{}, err
	}

	newUser := model.User{
		Name:            user.Name,
		Email:           user.Email,
		Password:        string(hash),
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

	if err := bcrypt.CompareHashAndPassword([]byte(storedUser.Password), []byte(user.Password)); err != nil {
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

func (uu *userUsecase) SetUserImage(user model.User, file []byte) (string, error) {
	imgUrl, err := uu.ur.SetUserImage(&user, file)
	if err != nil {
		return "", err
	}
	return imgUrl, nil
}

func (uu *userUsecase) GetUserImageUrlById(userId uint) (string, error) {
	imgUrl, err := uu.ur.GetUserImageUrlById(userId)
	if err != nil {
		return "", err
	}
	return imgUrl, nil
}
