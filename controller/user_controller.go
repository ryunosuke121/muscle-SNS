package controller

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/ryunosuke121/muscle-SNS/model"
	"github.com/ryunosuke121/muscle-SNS/usecase"
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
	user := model.User{}
	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	userRes, err := uc.uu.SignUp(user)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, userRes)
}

func (uc *userController) Login(c echo.Context) error {
	return nil
}

func (uc *userController) Logout(c echo.Context) error {
	return nil
}

func (uc *userController) GetUsersById(c echo.Context) error {
	strIds := c.QueryParams()["id"]
	ids, err := strConvertUint(strIds)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	var users []model.UserResponse
	for _, id := range ids {
		user := model.User{}
		res, _ := uc.uu.GetUserById(&user, uint(id))
		users = append(users, res)
	}
	return c.JSON(http.StatusOK, users)
}

func (uc *userController) UpdateUserName(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, errors.New("userId is empty").Error())
	}
	userId, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	name := c.FormValue("name")
	if name == "" {
		return c.JSON(http.StatusBadRequest, errors.New("name is empty").Error())
	}

	res, err := uc.uu.UpdateUserName(name, uint(userId))
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, res)
}

func (uc *userController) UpdateUserTrainingGroup(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, errors.New("userId is empty").Error())
	}
	userId, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	groupIdstr := c.FormValue("group_id")
	if groupIdstr == "" {
		return c.JSON(http.StatusBadRequest, errors.New("groupId is empty").Error())
	}
	groupId, err := strconv.ParseUint(groupIdstr, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	res, err := uc.uu.UpdateUserTrainingGroup(uint(groupId), uint(userId))
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, res)
}

func (uc *userController) UpdateUserImage(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, errors.New("userId is empty").Error())
	}

	userId, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	imageFile, err := c.FormFile("user_image")
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	res, err := uc.uu.UpdateUserImage(imageFile, uint(userId))
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
