package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/ryunosuke121/muscle-SNS/db"
	"github.com/ryunosuke121/muscle-SNS/model"
	"net/http"
	"time"
)

// 　トレーニング作成
func CreateTraining(c echo.Context) error {
	training := new(model.Training)
	if err := c.Bind(training); err != nil {
		return err
	}
	training.CreatedAt = time.Now()
	db := db.NewDB()
	db.Create(&training)
	res := Response{
		Message: "success",
		Data:    training,
	}
	return c.JSON(200, res)
}

// トレーニング取得
func GetTraining(c echo.Context) error {
	id := c.Param("id")
	training := new(model.Training)
	db := db.NewDB()
	db.First(&training, id)
	res := Response{
		Message: "success",
		Data:    training,
	}
	return c.JSON(200, res)
}

// ユーザーのトレーニング取得
func GetUserTrainings(c echo.Context) error {
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
func CreatePost(c echo.Context) error {
	post := new(model.Post)
	if err := c.Bind(post); err != nil {
		return err
	}

	post.CreatedAt = time.Now()
	db := db.NewDB()
	db.Create(&post)
	res := Response{
		Message: "success",
		Data:    post,
	}
	return c.JSON(200, res)
}

// ユーザーの投稿取得
func GetUserPosts(c echo.Context) error {
	id := c.Param("user_id")
	user := new(model.User)
	db := db.NewDB()
	db.Preload("Posts").First(&user, id)
	posts := new([]model.Post)
	posts = &user.Posts

	res := Response{
		Message: "success",
		Data:    posts,
	}
	return c.JSON(200, res)
}

// あるグループの投稿取得
func GetGroupPosts(c echo.Context) error {
	id := c.Param("group_id")
	users := new([]model.User)
	db := db.NewDB()
	err := db.Preload("Posts").Where("training_group_id = ?", id).Find(&users).Error
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	var posts []model.Post
	for _, user := range *users {
		for _, post := range user.Posts {
			posts = append(posts, post)
		}
	}
	res := Response{
		Message: "success",
		Data:    posts,
	}
	return c.JSON(200, res)
}
