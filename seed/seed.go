package main

import (
	"fmt"

	"github.com/ryunosuke121/muscle-SNS/src/db"
	"github.com/ryunosuke121/muscle-SNS/src/model"
)

func main() {
	db := db.NewDB()
	menus := []model.Menu{}
	menus = append(menus, model.Menu{ID: 1, Name: "ベンチプレス"})
	menus = append(menus, model.Menu{ID: 2, Name: "スクワット"})
	menus = append(menus, model.Menu{ID: 3, Name: "デッドリフト"})
	db.Create(&menus)
	training_groups := []model.TrainingGroup{}
	training_groups = append(training_groups, model.TrainingGroup{ID: 1, Name: "筋トレ好きあつまれ", ImageUrl: ""})
	training_groups = append(training_groups, model.TrainingGroup{ID: 2, Name: "みんなで筋トレ", ImageUrl: ""})
	training_groups = append(training_groups, model.TrainingGroup{ID: 3, Name: "マッチョになろうぜ。", ImageUrl: ""})
	training_groups = append(training_groups, model.TrainingGroup{ID: 4, Name: "ダンベル好きあつまれ", ImageUrl: ""})
	db.Create(&training_groups)
	users := []model.User{}
	users = append(users, model.User{ID: 1, Name: "田中太郎", Email: "tanaka@example.com", Password: "tekitou", TrainingGroupID: 1})
	users = append(users, model.User{ID: 2, Name: "佐藤花子", Email: "satou@example.com", Password: "tekitou", TrainingGroupID: 1})
	users = append(users, model.User{ID: 3, Name: "山田次郎", Email: "yamada@example.com", Password: "tekitou", TrainingGroupID: 2})
	db.Create(&users)

	fmt.Println("Successfully Seeded")
}
