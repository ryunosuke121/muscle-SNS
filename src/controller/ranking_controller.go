package controller

import (
	"log"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/ryunosuke121/muscle-SNS/src/application"
	"github.com/ryunosuke121/muscle-SNS/src/domain"
)

type IRankingController interface {
	GetUserTotalWeightInMonth(c echo.Context) error
	GetMonthRankingInGroup(c echo.Context) error
	GetMonthRankingInGroupByMenu(c echo.Context) error
}

type RankingController struct {
	rs application.IRankingService
}

func NewRankingController(rs application.IRankingService) IRankingController {
	return &RankingController{rs}
}

// ユーザーの総重量を取得する
func (rc *RankingController) GetUserTotalWeightInMonth(c echo.Context) error {
	userId := c.Param("user_id")
	req := new(GetUserTotalWeightInMonthRequestSchema)
	if err := c.Bind(req); err != nil {
		return c.String(http.StatusBadRequest, "invalid request")
	}
	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	utw, err := rc.rs.GetUserTotalWeightInMonth(c.Request().Context(), domain.UserID(userId), int(req.Year), int(req.Month))
	if err != nil {
		log.Println(err)
		return c.String(http.StatusInternalServerError, "failed to get user total weight")
	}

	return c.JSON(http.StatusOK, utw)
}

// グループ内のユーザーの月間ランキングを取得する
func (rc *RankingController) GetMonthRankingInGroup(c echo.Context) error {
	groupIdstr := c.Param("group_id")
	groupId, err := strconv.Atoi(groupIdstr)
	if err != nil {
		return c.String(http.StatusBadRequest, "invalid group id")
	}

	req := new(GetMonthRankingInGroupRequestSchema)
	if err := c.Bind(req); err != nil {
		return c.String(http.StatusBadRequest, "invalid request")
	}
	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	users, err := rc.rs.GetMonthRankingInGroup(c.Request().Context(), domain.UserGroupID(groupId), int(req.Year), int(req.Month))
	if err != nil {
		log.Println(err)
		return c.String(http.StatusInternalServerError, "failed to get month ranking in group")
	}

	return c.JSON(http.StatusOK, users)
}

// グループ内のユーザーの月間ランキングを取得する
func (rc *RankingController) GetMonthRankingInGroupByMenu(c echo.Context) error {
	groupIdstr := c.Param("group_id")
	groupId, err := strconv.Atoi(groupIdstr)
	if err != nil {
		return c.String(http.StatusBadRequest, "invalid group id")
	}
	req := new(GetMonthRankingInGroupByMenuResponseSchema)
	if err := c.Bind(req); err != nil {
		return c.String(http.StatusBadRequest, "invalid request")
	}
	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	users, err := rc.rs.GetMonthRankingInGroupByMenu(c.Request().Context(), domain.UserGroupID(groupId), domain.MenuID(req.MenuID), int(req.Year), int(req.Month))
	if err != nil {
		log.Println(err)
		return c.String(http.StatusInternalServerError, "failed to get month ranking in group by menu")
	}

	return c.JSON(http.StatusOK, users)
}
