package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/ryunosuke121/muscle-SNS/db"
	"github.com/ryunosuke121/muscle-SNS/model"
)

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
	db.Find(&groups)
	res := Response{
		Message: "success",
		Data:    groups,
	}
	return c.JSON(200, res)
}
