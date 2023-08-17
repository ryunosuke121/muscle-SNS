package controller

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/ryunosuke121/muscle-SNS/db"
	"github.com/ryunosuke121/muscle-SNS/model"
)

func CreateMenu(c echo.Context) error {
	var menu model.Menu
	if err := c.Bind(&menu); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	db := db.NewDB()
	db.Create(&menu)
	res := Response{
		Message: "success",
		Data:    menu,
	}
	return c.JSON(200, res)
}
