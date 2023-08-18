package router

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/ryunosuke121/muscle-SNS/controller"
)

func NewRouter(uc controller.IUserController) *echo.Echo {
	e := echo.New()

	e.POST("/signup", uc.SignUp)
	e.POST("/login", uc.Login)
	e.POST("/logout", uc.Logout)

	e.GET("/", func(c echo.Context) error {
		return c.HTML(http.StatusOK, "Hello, Docker! <3")
	})

	e.GET("/ping", func(c echo.Context) error {
		return c.JSON(http.StatusOK, struct{ Status string }{Status: "OK"})
	})

	//ユーザーの作成
	e.POST("/user", controller.CreateUser)
	//ユーザーの取得
	e.GET("/user/:id", controller.GetUser)
	//ユーザーを複数取得
	e.GET("/users", controller.GetUsers)
	//ユーザーの更新
	e.PUT("/user/:id", controller.UpdateUser)

	//トレーニングの作成
	e.POST("/training", controller.CreateTraining)
	//トレーニングの取得
	e.GET("/training/:training_id", controller.GetTraining)
	//あるユーザーのトレーニングを複数取得
	e.GET("/user/trainings/:user_id", controller.GetUserTrainings)

	//投稿の作成
	e.POST("/post", controller.CreatePost)
	//あるユーザーの投稿を複数取得
	e.GET("/user/post/:user_id", controller.GetUserPosts)
	//あるグループの投稿の取得
	e.GET("/group/posts/:group_id", controller.GetGroupPosts)

	//グループ一覧の取得
	e.GET("/groups", controller.GetGroups)

	//グループ内のランキングを取得
	e.GET("/group/ranking/:group_id", controller.GroupRanking)
	//ユーザーの総重量を取得
	e.GET("/user/total_weight/:user_id", controller.TotalWeight)

	return e
}