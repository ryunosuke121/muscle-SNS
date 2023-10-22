package controller

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/ryunosuke121/muscle-SNS/src/application"
	"github.com/ryunosuke121/muscle-SNS/src/domain"
	"github.com/ryunosuke121/muscle-SNS/utils/middleware"
)

type IUserController interface {
	SignUp(c echo.Context) error
	GetUsersByIds(c echo.Context) error
	UpdateUserName(c echo.Context) error
	UpdateUserGroup(c echo.Context) error
	UpdateUserImage(c echo.Context) error
}

type userController struct {
	us application.IUserService
}

func NewUserController(us application.IUserService) IUserController {
	return &userController{us}
}

func (uc *userController) SignUp(c echo.Context) error {
	ctx := c.Request().Context()

	req := new(SignUpRequestSchema)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	// ユーザーID, Emailをコンテキストから取得
	user_id, err := middleware.GetUserId(ctx)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	email, err := middleware.GetEmail(ctx)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	userRes, err := uc.us.SignUp(ctx, user_id, domain.UserName(req.Name), email)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusCreated, userRes)
}

// ユーザーを複数取得
func (uc *userController) GetUsersByIds(c echo.Context) error {
	ids := c.QueryParams()["id"]
	uids := make([]domain.UserID, len(ids))
	for i, id := range ids {
		uids[i] = domain.UserID(id)
	}

	res, err := uc.us.GetUsersByIds(c.Request().Context(), uids)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, res)
}

// ユーザー名の更新
func (uc *userController) UpdateUserName(c echo.Context) error {
	ctx := c.Request().Context()

	req := new(UpdateUserNameRequestSchema)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	userId, err := middleware.GetUserId(ctx)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	res, err := uc.us.UpdateUserName(ctx, userId, domain.UserName(req.Name))
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, res)
}

func (uc *userController) UpdateUserGroup(c echo.Context) error {
	ctx := c.Request().Context()

	req := new(UpdateUserGroupRequestSchema)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	userId, err := middleware.GetUserId(ctx)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	res, err := uc.us.UpdateUserGroup(ctx, userId, domain.UserGroupID(req.GroupID))
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, res)
}

func (uc *userController) UpdateUserImage(c echo.Context) error {
	imageFile, err := c.FormFile("user_image")
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	userId, err := middleware.GetUserId(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	res, err := uc.us.UpdateUserImage(c.Request().Context(), userId, imageFile)
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
