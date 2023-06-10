package handler

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/ryunosuke121/muscle-SNS/db"
	"github.com/ryunosuke121/muscle-SNS/model"
	"net/http"
)

type Response struct {
	Message string `json:"message"`
	Data    any    `json:"data"`
}

func CreateUser(c echo.Context) error {
	user := new(model.User)
	if err := c.Bind(user); err != nil {
		return err
	}

	validate := validator.New()
	if err := validate.Struct(user); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	db := db.NewDB()
	db.Create(&user)
	res := Response{
		Message: "success",
		Data:    user,
	}
	return c.JSON(200, res)
}

func GetUser(c echo.Context) error {
	id := c.Param("id")
	user := new(model.User)
	db := db.NewDB()
	db.First(&user, id)
	res := Response{
		Message: "success",
		Data:    user,
	}
	return c.JSON(200, res)
}

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

func UpdateUser(c echo.Context) error {
	id := c.Param("id")
	// リクエストボディを取得
	request := new(model.User)
	if err := c.Bind(request); err != nil {
		return err
	}

	user := new(model.User)
	db := db.NewDB()
	db.First(&user, id)
	// パスワードが一致しない場合はエラー
	if user.Password != request.Password {
		return c.JSON(http.StatusBadRequest, "password is wrong")
	}
	data := map[string]interface{}{}

	if request.NAME != "" {
		data["name"] = request.NAME
	}
	if request.ImageUrl != "" {
		data["image_url"] = request.ImageUrl
	}
	if request.GroupId != 0 {
		data["group_id"] = request.GroupId
	}

	db.Model(&user).Updates(data)
	res := Response{
		Message: "success",
		Data:    user,
	}
	return c.JSON(200, res)
}
