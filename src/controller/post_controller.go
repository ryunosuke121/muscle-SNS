package controller

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/ryunosuke121/muscle-SNS/src/application"
	"github.com/ryunosuke121/muscle-SNS/src/domain"
	"github.com/ryunosuke121/muscle-SNS/utils/middleware"
)

type IPostController interface {
	GetPostsByIds(c echo.Context) error
	GetMyPosts(c echo.Context) error
	GetUserPosts(c echo.Context) error
	CreatePost(c echo.Context) error
	DeletePost(c echo.Context) error
	GetGroupPosts(c echo.Context) error
}

type PostController struct {
	ps application.IPostService
}

func NewPostController(ps application.IPostService) IPostController {
	return &PostController{ps}
}

// 投稿を複数件取得する
func (pc *PostController) GetPostsByIds(c echo.Context) error {
	ids := c.QueryParams()["id"]
	postIds := make([]domain.PostID, len(ids))
	for i, id := range ids {
		postId, err := strconv.Atoi(id)
		if err != nil {
			return c.String(http.StatusBadRequest, "invalid id")
		}
		postIds[i] = domain.PostID(postId)
	}

	posts, err := pc.ps.GetPostsByIds(c.Request().Context(), postIds)
	if err != nil {
		return c.String(http.StatusInternalServerError, "failed to get posts")
	}

	return c.JSON(http.StatusOK, posts)
}

// 自分の投稿を複数取得する
func (pc *PostController) GetMyPosts(c echo.Context) error {
	ctx := c.Request().Context()
	userId, err := middleware.GetUserId(ctx)
	if err != nil {
		return c.String(http.StatusBadRequest, "invalid user id")
	}

	posts, err := pc.ps.GetUserPosts(ctx, userId)
	if err != nil {
		return c.String(http.StatusInternalServerError, "failed to get posts")
	}

	return c.JSON(http.StatusOK, posts)
}

// ユーザーの投稿を複数取得する
func (pc *PostController) GetUserPosts(c echo.Context) error {
	userId := c.Param("user_id")

	posts, err := pc.ps.GetUserPosts(c.Request().Context(), domain.UserID(userId))
	if err != nil {
		return c.String(http.StatusInternalServerError, "failed to get posts")
	}

	return c.JSON(http.StatusOK, posts)
}

// 投稿を作成する
func (pc *PostController) CreatePost(c echo.Context) error {
	ctx := c.Request().Context()
	userId, err := middleware.GetUserId(ctx)
	if err != nil {
		return c.String(http.StatusBadRequest, "invalid user id")
	}

	image, err := c.FormFile("image")
	if err != nil {
		return c.String(http.StatusBadRequest, "invalid image")
	}

	post := new(CreatePostRequestSchema)
	if err := json.Unmarshal([]byte(c.FormValue("post")), &post); err != nil {
		log.Print(err.Error())
		return c.String(http.StatusBadRequest, "invalid post")
	}

	createPostReq := &application.CreatePostRequest{
		UserID:  userId,
		Comment: post.Comment,
		Image:   image,
		Training: &application.TrainingRequest{
			MenuID: domain.MenuID(post.Training.MenuID),
			Times:  post.Training.Times,
			Weight: post.Training.Weight,
			Sets:   post.Training.Sets,
		},
	}

	res, err := pc.ps.CreatePost(ctx, createPostReq)
	if err != nil {
		log.Print(err.Error())
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, res)
}

// 投稿を削除する
func (pc *PostController) DeletePost(c echo.Context) error {
	postId, err := strconv.Atoi(c.Param("post_id"))
	if err != nil {
		return c.String(http.StatusBadRequest, "invalid post id")
	}

	userId, err := middleware.GetUserId(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	err = pc.ps.DeletePost(c.Request().Context(), userId, domain.PostID(postId))
	if err != nil {
		if err == domain.ErrForbidden {
			return c.JSON(http.StatusForbidden, err.Error())
		}
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	return c.String(http.StatusOK, "success")
}

// グループの投稿を取得する
func (pc *PostController) GetGroupPosts(c echo.Context) error {
	groupId, err := strconv.Atoi(c.Param("group_id"))
	if err != nil {
		return c.String(http.StatusBadRequest, "invalid group id")
	}

	posts, err := pc.ps.GetGroupPosts(c.Request().Context(), domain.UserGroupID(groupId))
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, posts)
}
