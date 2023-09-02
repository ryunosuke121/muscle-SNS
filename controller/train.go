package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/ryunosuke121/muscle-SNS/db"
	"github.com/ryunosuke121/muscle-SNS/model"
)

type RequestPost struct {
	UserID     uint           `json:"user_id"`
	TrainingID uint           `json:"training_id"`
	Training   model.Training `json:"training"`
	Comment    string         `json:"comment"`
}
type ResponsePosts struct {
	mu    sync.Mutex
	Posts []model.Post `json:"posts"`
}

type Response struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type TrainingController struct {
	s3Client      *s3.Client
	presignClient *s3.PresignClient
}

func NewTrainingController(s3client *s3.Client, presignClient *s3.PresignClient) *TrainingController {
	return &TrainingController{s3client, presignClient}
}

// トレーニング作成
func (tc *TrainingController) CreateTraining(c echo.Context) error {
	training := new(model.Training)
	if err := c.Bind(training); err != nil {
		return err
	}
	training.CreatedAt = time.Now()
	db := db.NewDB()
	db.Create(&training)
	db.Preload("User").Preload("Menu").Find(&training)
	res := Response{
		Message: "success",
		Data:    training,
	}
	return c.JSON(200, res)
}

// トレーニング取得
func (tc *TrainingController) GetTraining(c echo.Context) error {
	id := c.Param("training_id")
	training := new(model.Training)
	db := db.NewDB()
	db.Preload("User").Preload("Menu").Where("id = ?", id).First(&training)
	res := Response{
		Message: "success",
		Data:    training,
	}
	return c.JSON(200, res)
}

// ユーザーのトレーニング取得
func (tc *TrainingController) GetUserTrainings(c echo.Context) error {
	id := c.Param("user_id")
	trainings := new([]model.Training)
	db := db.NewDB()
	db.Preload("User").Preload("Menu").Where("user_id = ?", id).Find(&trainings)
	res := Response{
		Message: "success",
		Data:    trainings,
	}
	return c.JSON(200, res)
}

// 投稿作成
func (tc *TrainingController) CreatePost(c echo.Context) error {
	//　投稿データの受取
	post := c.FormValue("post_info")
	var requestPost RequestPost
	err := json.Unmarshal([]byte(post), &requestPost)
	if err != nil {
		return err
	}
	// 画像の受取
	imageFile, err := c.FormFile("post_image")
	imgUrl := ""
	//画像がある場合はS3に保存
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
		imgUrl = fmt.Sprintf("post_image/%s%s", u.String(), imageFile.Filename)
		// s3に画像を保存
		param := &s3.PutObjectInput{
			Bucket: aws.String(os.Getenv("BUCKET_NAME")),
			Key:    aws.String(imgUrl),
			Body:   src,
		}
		_, err = tc.s3Client.PutObject(context.TODO(), param)
		if err != nil {
			return err
		}
	}
	// 投稿データの作成
	savePost := model.Post{
		UserID:     requestPost.UserID,
		TrainingID: requestPost.TrainingID,
		Training:   requestPost.Training,
		Comment:    requestPost.Comment,
		ImageUrl:   imgUrl,
		CreatedAt:  time.Now(),
	}
	db := db.NewDB()
	db.Create(&savePost)
	res := Response{
		Message: "success",
		Data:    savePost,
	}
	return c.JSON(200, res)
}

// ユーザーの投稿取得
func (tc *TrainingController) GetUserPosts(c echo.Context) error {
	id := c.Param("user_id")
	var user model.User
	db := db.NewDB()
	db.Preload("Posts.Training.Menu").Find(&user, id)
	posts := user.Posts
	var resPosts ResponsePosts
	var wg sync.WaitGroup
	for _, post := range posts {
		wg.Add(1)
		go func(post model.Post) {
			defer wg.Done()
			if post.ImageUrl == "" {
				resPosts.mu.Lock()
				defer resPosts.mu.Unlock()
				resPosts.Posts = append(resPosts.Posts, post)
				return
			}
			// 画像を取得
			param := &s3.GetObjectInput{
				Bucket: aws.String(os.Getenv("BUCKET_NAME")),
				Key:    aws.String(post.ImageUrl),
			}
			rq, err := tc.presignClient.PresignGetObject(context.Background(), param)
			if err != nil {
				return
			}
			fmt.Println(rq)
			post.ImageUrl = rq.URL
			resPosts.mu.Lock()
			defer resPosts.mu.Unlock()
			resPosts.Posts = append(resPosts.Posts, post)
		}(post)
	}
	wg.Wait()
	res := Response{
		Message: "success",
		Data:    resPosts.Posts,
	}
	return c.JSON(200, res)
}
