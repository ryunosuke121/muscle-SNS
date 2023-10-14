package router

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/ryunosuke121/muscle-SNS/src/controller"
	"github.com/ryunosuke121/muscle-SNS/src/middleware"
)

func NewRouter(uc controller.IUserController, tc controller.TrainingController, gc controller.GroupController) *echo.Echo {
	e := echo.New()

	client, err := middleware.NewAuthClient()
	if err != nil {
		e.Logger.Fatal(err)
	}

	e.POST("/signup", uc.SignUp, client.CheckToken)
	e.POST("/login", uc.Login, client.CheckToken)
	e.POST("/logout", uc.Logout)

	e.GET("/ping", func(c echo.Context) error {
		return c.JSON(http.StatusOK, struct{ Status string }{Status: "OK"})
	})

	//ユーザーを複数取得
	e.GET("/user", uc.GetUsersById)

	u := e.Group("/user", client.CheckToken)
	//ユーザー名の更新
	u.PUT("/name", uc.UpdateUserName)
	//ユーザーのトレーニンググループの更新
	u.PUT("/training_group", uc.UpdateUserTrainingGroup)
	//ユーザーの画像の更新
	u.PUT("/image", uc.UpdateUserImage)

	//トレーニングの作成
	e.POST("/training", tc.CreateTraining)
	//トレーニングの取得
	e.GET("/training/:training_id", tc.GetTraining)
	//あるユーザーのトレーニングを複数取得
	e.GET("/user/trainings/:user_id", tc.GetUserTrainings)

	//投稿の作成
	e.POST("/post", tc.CreatePost)
	//あるユーザーの投稿を複数取得
	e.GET("/user/post/:user_id", tc.GetUserPosts)
	//あるグループの投稿の取得
	e.GET("/group/posts/:group_id", gc.GetGroupPosts)

	//グループ一覧の取得
	e.GET("/groups", gc.GetGroups)

	//グループ内のランキングを取得
	e.GET("/group/ranking/:group_id", controller.GroupRanking)
	//ユーザーの総重量を取得
	e.GET("/user/total_weight/:user_id", controller.TotalWeight)

	return e
}
