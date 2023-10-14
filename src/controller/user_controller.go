package controller

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/ryunosuke121/muscle-SNS/src/middleware"
	"github.com/ryunosuke121/muscle-SNS/src/model"
	"github.com/ryunosuke121/muscle-SNS/src/usecase"
)

type IUserController interface {
	SignUp(c echo.Context) error
	Login(c echo.Context) error
	Logout(c echo.Context) error
	GetUsersById(c echo.Context) error
	UpdateUserName(c echo.Context) error
	UpdateUserTrainingGroup(c echo.Context) error
	UpdateUserImage(c echo.Context) error
}

type userController struct {
	uu usecase.IUserUseCase
}

func NewUserController(uu usecase.IUserUseCase) IUserController {
	return &userController{uu}
}

func (uc *userController) SignUp(c echo.Context) error {
	req := SignUpRequestSchema{}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	decodedToken := middleware.GetDecodedToken(c.Request().Context())
	if decodedToken == nil {
		return c.JSON(http.StatusAlreadyReported, errors.New("token is empty").Error())
	}
	user_id := (*decodedToken).Claims["user_id"].(string)
	email := (*decodedToken).Claims["email"].(string)

	user := model.User{
		ID:              user_id,
		Name:            req.Name,
		Email:           email,
		TrainingGroupID: req.TrainingGroupID,
	}
	userRes, err := uc.uu.SignUp(user)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, userRes)
}

func (uc *userController) Login(c echo.Context) error {
	decodedToken := middleware.GetDecodedToken(c.Request().Context())
	if decodedToken == nil {
		return c.JSON(http.StatusAlreadyReported, errors.New("token is empty").Error())
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(*decodedToken); err != nil {
		panic(err)
	}

	return c.JSON(http.StatusOK, decodedToken)
}

func (uc *userController) Logout(c echo.Context) error {
	return nil
}

func (uc *userController) GetUsersById(c echo.Context) error {
	ids := c.QueryParams()["id"]

	var users []model.UserResponse
	for _, id := range ids {
		user := model.User{}
		res, err := uc.uu.GetUserById(&user, id)
		if err == nil {
			users = append(users, res)
		}
	}
	return c.JSON(http.StatusOK, users)
}

func (uc *userController) UpdateUserName(c echo.Context) error {
	userId := c.Param("id")
	if userId == "" {
		return c.JSON(http.StatusBadRequest, errors.New("userId is empty").Error())
	}

	name := c.FormValue("name")
	if name == "" {
		return c.JSON(http.StatusBadRequest, errors.New("name is empty").Error())
	}

	res, err := uc.uu.UpdateUserName(name, userId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, res)
}

func (uc *userController) UpdateUserTrainingGroup(c echo.Context) error {
	userId := c.Param("id")
	if userId == "" {
		return c.JSON(http.StatusBadRequest, errors.New("userId is empty").Error())
	}

	groupIdstr := c.FormValue("group_id")
	if groupIdstr == "" {
		return c.JSON(http.StatusBadRequest, errors.New("groupId is empty").Error())
	}
	groupId, err := strconv.ParseUint(groupIdstr, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	res, err := uc.uu.UpdateUserTrainingGroup(uint(groupId), userId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, res)
}

func (uc *userController) UpdateUserImage(c echo.Context) error {
	userId := c.Param("id")
	if userId == "" {
		return c.JSON(http.StatusBadRequest, errors.New("userId is empty").Error())
	}

	imageFile, err := c.FormFile("user_image")
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	res, err := uc.uu.UpdateUserImage(imageFile, userId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, res)
}

func strConvertUint(strIds []string) ([]uint64, error) {
	var resultIds []uint64
	for _, str := range strIds {
		id, err := strconv.ParseUint(str, 10, 64)
		if err != nil {
			return make([]uint64, 0), err
		}
		resultIds = append(resultIds, id)
	}
	return resultIds, nil
}
