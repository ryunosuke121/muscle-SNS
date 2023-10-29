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
	GetPosts(c echo.Context) error
	CreatePost(c echo.Context) error
	DeletePost(c echo.Context) error
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

	options := &domain.GetPostsOptions{
		UserId: &userId,
	}

	posts, err := pc.ps.GetPostsByOptions(ctx, options)
	if err != nil {
		return c.String(http.StatusInternalServerError, "failed to get posts")
	}

	return c.JSON(http.StatusOK, posts)
}

// クエリを指定して投稿を取得する
func (pc *PostController) GetPosts(c echo.Context) error {
	options := &domain.GetPostsOptions{}

	if userId := c.QueryParam("user_id"); userId == "" {
		uid := domain.UserID(userId)
		options.UserId = &uid
	}

	if menuId := c.QueryParam("menu_id"); menuId == "" {
		m, err := strconv.Atoi(menuId)
		if err != nil {
			return c.String(http.StatusBadRequest, "invalid menu id")
		}
		mid := domain.MenuID(m)
		options.MenuId = &mid
	}

	if groupId := c.QueryParam("group_id"); groupId == "" {
		g, err := strconv.Atoi(groupId)
		if err != nil {
			return c.String(http.StatusBadRequest, "invalid group id")
		}
		gid := domain.UserGroupID(g)
		options.GroupId = &gid
	}

	if year := c.QueryParam("year"); year == "" {
		y, err := strconv.Atoi(year)
		if err != nil {
			return c.String(http.StatusBadRequest, "invalid year")
		}
		options.Year = &y
	}

	if month := c.QueryParam("month"); month == "" {
		m, err := strconv.Atoi(month)
		if err != nil {
			return c.String(http.StatusBadRequest, "invalid month")
		}
		options.Month = &m
	}

	if limit := c.QueryParam("limit"); limit == "" {
		l, err := strconv.Atoi(limit)
		if err != nil {
			return c.String(http.StatusBadRequest, "invalid limit")
		}
		options.Limit = &l
	}

	if cursor := c.QueryParam("cursor"); cursor == "" {
		cs, err := strconv.Atoi(cursor)
		if err != nil {
			return c.String(http.StatusBadRequest, "invalid cursor")
		}
		options.Cursor = &cs
	}

	posts, err := pc.ps.GetPostsByOptions(c.Request().Context(), options)
	if err != nil {
		log.Print(err.Error())
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
