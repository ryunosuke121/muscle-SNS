package handler

import (
	"net/http"
	"sort"
	"strconv"
	"sync"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/ryunosuke121/muscle-SNS/db"
	"github.com/ryunosuke121/muscle-SNS/model"
)

type RankMap struct {
	mu  sync.Mutex
	utl UserTotalWeightList
}

type UserTotalWeight struct {
	User         model.User
	Total_weight uint
}

type UserTotalWeightList []UserTotalWeight

func (utl UserTotalWeightList) Len() int {
	return len(utl)
}

func (utl UserTotalWeightList) Less(i, j int) bool {
	return utl[i].Total_weight > utl[j].Total_weight
}

func (utl UserTotalWeightList) Swap(i, j int) {
	utl[i], utl[j] = utl[j], utl[i]
}

// ユーザーの総重量を取得
func TotalWeight(c echo.Context) error {
	user_id, err := strconv.Atoi(c.Param("user_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	var user model.User
	db := db.NewDB()
	db.Preload("Trainings").Find(&user, user_id)
	total_weight := calcUserWeight(user)
	res := Response{
		Message: "success",
		Data:    total_weight}
	return c.JSON(200, res)
}

// グループ内のランキングを取得
func GroupRanking(c echo.Context) error {
	group_id := c.Param("group_id")
	users := new([]model.User)
	db := db.NewDB()
	db.Where("training_group_id = ?", group_id).Find(&users)
	rankMap := RankMap{utl: []UserTotalWeight{}}

	var wg sync.WaitGroup
	//　並行処理でユーザーの総重量を計算
	for _, user := range *users {
		wg.Add(1)
		go func(user model.User) {
			defer wg.Done()
			total_weight := calcUserWeight(user)
			rankMap.mu.Lock()
			rankMap.utl = append(rankMap.utl, UserTotalWeight{user, total_weight})
			rankMap.mu.Unlock()
		}(user)
	}
	wg.Wait()
	sort.Sort(rankMap.utl)
	res := Response{
		Message: "success",
		Data:    rankMap.utl}
	return c.JSON(200, res)
}

// ユーザーの1ヶ月に持ち上げた総重量を計算
func calcUserWeight(user model.User) uint {
	db := db.NewDB()
	db.Preload("Trainings").Find(&user)
	oneMonthAgo := time.Now().AddDate(0, -1, 0)

	var total_weight uint = 0
	for _, training := range user.Trainings {
		if training.CreatedAt.After(oneMonthAgo) {
			total_weight += training.Weight * training.Times * training.Sets
		}
	}
	return total_weight
}
