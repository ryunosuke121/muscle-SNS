package handler

import (
	"net/http"

	"context"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/labstack/echo/v4"
	"github.com/ryunosuke121/muscle-SNS/db"
	"github.com/ryunosuke121/muscle-SNS/model"
	"github.com/ryunosuke121/muscle-SNS/s3client"
)

// あるグループの投稿取得
func GetGroupPosts(c echo.Context) error {
	id := c.Param("group_id")
	users := new([]model.User)
	db := db.NewDB()
	err := db.Preload("Posts.User").Preload("Posts.Training.Menu").Where("training_group_id = ?", id).Find(&users).Error
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	var posts []model.Post
	for _, user := range *users {
		posts = append(posts, user.Posts...)
	}

	res := Response{
		Message: "success",
		Data:    posts,
	}
	return c.JSON(200, res)
}

// グループ一覧取得
func GetGroups(c echo.Context) error {
	groups := new([]model.TrainingGroup)
	db := db.NewDB()
	db.Preload("Users").Find(&groups)

	// 画像のURLを取得
	for _, group := range *groups {
		if group.ImageUrl != "" {
			param := &s3.GetObjectInput{
				Bucket: aws.String(os.Getenv("BUCKET_NAME")),
				Key:    aws.String(group.ImageUrl),
			}
			rq, err := s3client.PresignClient.PresignGetObject(context.Background(), param)
			if err != nil {
				return c.JSON(http.StatusBadRequest, err.Error())
			}
			group.ImageUrl = rq.URL
		}
	}

	res := Response{
		Message: "success",
		Data:    *groups,
	}
	return c.JSON(200, res)
}
