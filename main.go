package main

import (
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/ryunosuke121/muscle-SNS/handler"
	"github.com/ryunosuke121/muscle-SNS/s3client"
)

func main() {

	// AWS S3の設定
	s3client.InitS3Client()

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
	e.GET("/group/posts/:group_id", handler.GetGroupPosts)

	//グループ一覧の取得
	e.GET("/groups", handler.GetGroups)

	//グループ内のランキングを取得
	e.GET("/group/ranking/:group_id", handler.GroupRanking)
	//ユーザーの総重量を取得
	e.GET("/user/total_weight/:user_id", handler.TotalWeight)

	httpPort := os.Getenv("HTTP_PORT")
	if httpPort == "" {
		httpPort = "8080"
	}

	e.Logger.Fatal(e.Start(":" + httpPort))
}
