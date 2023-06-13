package handler

import (
	"net/http"

	"context"
	"os"

	"fmt"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/labstack/echo/v4"
	"github.com/ryunosuke121/muscle-SNS/db"
	"github.com/ryunosuke121/muscle-SNS/model"
	"github.com/ryunosuke121/muscle-SNS/s3client"
)

type ResponseGroupPosts struct {
	mu    sync.Mutex
	Posts []model.Post `json:"posts"`
}

// あるグループの投稿取得
func GetGroupPosts(c echo.Context) error {
	id := c.Param("group_id")
	users := new([]model.User)
	db := db.NewDB()
	err := db.Preload("Posts.User").Preload("Posts.Training.Menu").Where("training_group_id = ?", id).Find(&users).Error
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	response_group_posts := ResponseGroupPosts{}
	var wg sync.WaitGroup
	for _, user := range *users {
		for _, post := range user.Posts {
			wg.Add(1)
			go func(post model.Post) {
				defer wg.Done()
				if post.ImageUrl == "" {
					response_group_posts.mu.Lock()
					defer response_group_posts.mu.Unlock()
					response_group_posts.Posts = append(response_group_posts.Posts, post)
					return
				}
				// 画像を取得
				param := &s3.GetObjectInput{
					Bucket: aws.String(os.Getenv("BUCKET_NAME")),
					Key:    aws.String(post.ImageUrl),
				}
				rq, err := s3client.PresignClient.PresignGetObject(context.Background(), param)
				if err != nil {
					return
				}
				fmt.Println(rq)
				post.ImageUrl = rq.URL
				response_group_posts.mu.Lock()
				defer response_group_posts.mu.Unlock()
				response_group_posts.Posts = append(response_group_posts.Posts, post)
			}(post)
		}
		wg.Wait()
	}

	res := Response{
		Message: "success",
		Data:    response_group_posts.Posts,
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
