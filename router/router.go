package router

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/ryunosuke121/muscle-SNS/src/controller"
	"github.com/ryunosuke121/muscle-SNS/utils/middleware"
)

func NewRouter(uc controller.IUserController, pc controller.IPostController) *echo.Echo {
	e := echo.New()

	client, err := middleware.NewAuthClient()
	if err != nil {
		e.Logger.Fatal(err)
	}

	e.POST("/signup", uc.SignUp, client.CheckToken)

	e.GET("/ping", func(c echo.Context) error {
		return c.JSON(http.StatusOK, struct{ Status string }{Status: "OK"})
	})

	u := e.Group("/user", client.CheckToken)
	//ユーザーを複数取得
	u.GET("", uc.GetUsersByIds)
	//ユーザー名の更新
	u.PUT("/name", uc.UpdateUserName)
	//ユーザーのトレーニンググループの更新
	u.PUT("/user_group", uc.UpdateUserGroup)
	//ユーザーの画像の更新
	u.PUT("/image", uc.UpdateUserImage)
	//あるユーザーの投稿を複数取得
	u.GET("/post/:user_id", pc.GetUserPosts)
	// 自分の投稿を複数取得する
	u.GET("/post/my", pc.GetMyPosts)

	p := e.Group("/post", client.CheckToken)
	// 投稿を複数件取得する
	p.GET("", pc.GetPostsByIds)
	//投稿の作成
	p.POST("", pc.CreatePost)
	//あるグループの投稿の取得
	p.GET("/group/:group_id", pc.GetGroupPosts)
	// 投稿を削除する
	p.DELETE("/:post_id", pc.DeletePost)
	// //グループ一覧の取得
	// e.GET("/groups", gc.GetGroups)

	// //グループ内のランキングを取得
	// e.GET("/group/ranking/:group_id", controller.GroupRanking)
	// //ユーザーの総重量を取得
	// e.GET("/user/total_weight/:user_id", controller.TotalWeight)

	return e
}
