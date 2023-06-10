package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/ryunosuke121/muscle-SNS/handler"
	"net/http"
	"os"
)

func main() {

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/", func(c echo.Context) error {
		return c.HTML(http.StatusOK, "Hello, Docker! <3")
	})

	e.GET("/ping", func(c echo.Context) error {
		return c.JSON(http.StatusOK, struct{ Status string }{Status: "OK"})
	})

	//ユーザーの作成
	e.POST("/user", handler.CreateUser)
	//ユーザーの取得
	e.GET("/user/:id", handler.GetUser)
	//ユーザーを複数取得
	e.GET("/users", handler.GetUsers)
	//ユーザーの更新
	e.PUT("/user/:id", handler.UpdateUser)

	//トレーニングの作成
	e.POST("/training", handler.CreateTraining)
	//トレーニングの取得
	e.GET("/training/:training_id", handler.GetTraining)
	//あるユーザーのトレーニングを複数取得
	e.GET("/user/trainings/:user_id", handler.GetUserTrainings)

	//投稿の作成
	e.POST("/post", handler.CreatePost)
	//あるユーザーの投稿を複数取得
	e.GET("/user/post/:user_id", handler.GetUserPosts)
	//あるグループの投稿の取得
	e.GET("/post/:group_id", handler.GetGroupPosts)

	httpPort := os.Getenv("HTTP_PORT")
	if httpPort == "" {
		httpPort = "8080"
	}

	e.Logger.Fatal(e.Start(":" + httpPort))
}
