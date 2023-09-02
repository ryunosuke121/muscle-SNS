package controller

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/ryunosuke121/muscle-SNS/model"
	"github.com/ryunosuke121/muscle-SNS/usecase"
)

type ITrainController interface {
	GetTrainingById(c echo.Context) error
	GetUserTrainings(c echo.Context) error
	CreatePost(c echo.Context) error
	GetUserPostsById(c echo.Context) error
}

type TrainController struct {
	tu usecase.ITrainUsecase
}

func NewTrainController(tu usecase.ITrainUsecase) ITrainController {
	return &TrainController{tu: tu}
}

// トレーニング取得
func (tc *TrainController) GetTrainingById(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, "training_id is empty")
	}
	trainingId, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	res, err := tc.tu.GetTrainingById(uint(trainingId))
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, res)
}

// ユーザーのトレーニング取得
func (tc *TrainController) GetUserTrainings(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, "user_id is empty")
	}
	userId, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	res, err := tc.tu.GetUserTrainings(uint(userId))
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, res)
}

// 投稿の作成
func (tc *TrainController) CreatePost(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, "user_id is empty")
	}
	userId, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	post := model.Post{}
	if err := c.Bind(&post); err != nil {
		return err
	}
	post.UserID = uint(userId)
	if err := tc.tu.CreatePost(&post); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, post)
}

// あるユーザーの投稿を複数取得
func (tc *TrainController) GetUserPostsById(c echo.Context) error {
	return nil
}
