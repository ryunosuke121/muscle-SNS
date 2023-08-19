package controller

import (
	"net/http"

	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/ryunosuke121/muscle-SNS/db"
	"github.com/ryunosuke121/muscle-SNS/model"
	"github.com/ryunosuke121/muscle-SNS/s3client"
)

type Response struct {
	Message string `json:"message"`
	Data    any    `json:"data"`
}

type CreateUserRequest struct {
	Name            string `json:"name" validate:"required"`
	Email           string `json:"email" validate:"required,email"`
	Password        string `json:"password" validate:"required,min=8,max=32"`
	TrainingGroupID uint   `json:"training_group_id"`
}

type UpdateUserRequest struct {
	Name            string `json:"name" validate:"required"`
	Email           string `json:"email" validate:"required,email"`
	TrainingGroupID uint   `json:"training_group_id"`
}

func CreateUser(c echo.Context) error {
	user_info := c.FormValue("user_info")
	var create_user_request CreateUserRequest
	err := json.Unmarshal([]byte(user_info), &create_user_request)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	user := model.User{
		Name:            create_user_request.Name,
		Email:           create_user_request.Email,
		Password:        create_user_request.Password,
		TrainingGroupID: create_user_request.TrainingGroupID,
	}

	validate := validator.New()
	if err := validate.Struct(user); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	if user.TrainingGroupID == 0 {
		user.TrainingGroupID = 1
	}
	encryptPw, err := crypto.PasswordEncrypt(user.Password)
	if err != nil {
		fmt.Println("パスワード暗号化中にエラーが発生しました。：", err)
		return err
	}
	user.Password = encryptPw

	imageFile, err := c.FormFile("user_image")
	imgUrl := ""
	if err == nil {
		src, err := imageFile.Open()
		if err != nil {
			return err
		}
		defer src.Close()
		//キーにuuidを含める
		u, err := uuid.NewRandom()
		if err != nil {
			return err
		}

		imgUrl = fmt.Sprintf("user_image/%s%s", u.String(), imageFile.Filename)

		// s3に画像を保存
		param := &s3.PutObjectInput{
			Bucket: aws.String(os.Getenv("BUCKET_NAME")),
			Key:    aws.String(imgUrl),
			Body:   src,
		}
		_, err = s3client.S3Client.PutObject(context.TODO(), param)
		if err != nil {
			return err
		}
	}
	user.ImageUrl = imgUrl
	db := db.NewDB()
	err = db.Create(&user).Error
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	db.Preload("TrainingGroup").Find(&user)
	res := Response{
		Message: "success",
		Data:    user,
	}
	return c.JSON(200, res)
}

// ユーザー取得
func GetUser(c echo.Context) error {
	id := c.Param("id")
	user := new(model.User)
	db := db.NewDB()
	db.First(&user, id)
	db.Preload("TrainingGroup").Find(&user)
	//s3から画像を取得
	if user.ImageUrl != "" {
		param := &s3.GetObjectInput{
			Bucket: aws.String(os.Getenv("BUCKET_NAME")),
			Key:    aws.String(user.ImageUrl),
		}
		res, err := s3client.PresignClient.PresignGetObject(context.Background(), param)
		if err != nil {
			return err
		}
		user.ImageUrl = res.URL
	}
	res := Response{
		Message: "success",
		Data:    user,
	}
	return c.JSON(200, res)
}

// ユーザーを複数取得
func GetUsers(c echo.Context) error {
	ids := c.QueryParams()["id"]
	users := new([]model.User)
	db := db.NewDB()
	db.Find(&users, ids)
	res := Response{
		Message: "success",
		Data:    users,
	}
	return c.JSON(200, res)
}

// ユーザー更新
func UpdateUser(c echo.Context) error {
	id := c.Param("id")
	// リクエストボディを取得
	user_info := c.FormValue("user_info")
	var request UpdateUserRequest
	if err := json.Unmarshal([]byte(user_info), &request); err != nil {
		return err
	}

	user := new(model.User)
	db := db.NewDB()
	db.First(&user, id)
	if user.ID == 0 {
		return c.JSON(http.StatusBadRequest, "ユーザーが見つかりませんでした。")
	}

	data := map[string]interface{}{}

	if request.Name != "" {
		data["name"] = request.Name
	}
	if request.Email != "" {
		data["email"] = request.Email
	}
	if request.TrainingGroupID != 0 {
		data["training_group_id"] = request.TrainingGroupID
	}

	// TODO: 画像のアップデート処理
	imageFile, err := c.FormFile("user_image")
	if err == nil {
		src, err := imageFile.Open()
		if err != nil {
			return err
		}
		defer src.Close()
		//キーにuuidを含める
		u, err := uuid.NewRandom()
		if err != nil {
			return err
		}

		imgUrl := fmt.Sprintf("user_image/%s%s", u.String(), imageFile.Filename)

		// s3に画像を保存
		param := &s3.PutObjectInput{
			Bucket: aws.String(os.Getenv("BUCKET_NAME")),
			Key:    aws.String(imgUrl),
			Body:   src,
		}
		_, err = s3client.S3Client.PutObject(context.TODO(), param)
		if err != nil {
			return err
		}
		data["image_url"] = imgUrl
	}

	db.Model(&user).Updates(data)
	res := Response{
		Message: "success",
		Data:    user,
	}
	return c.JSON(200, res)
}
