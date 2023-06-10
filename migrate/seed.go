package main

import (
	"fmt"
	"github.com/ryunosuke121/muscle-SNS/db"
	"github.com/ryunosuke121/muscle-SNS/model"
)

func main() {
	db := db.NewDB()
	menus := []model.Menu{}
	menus.append(model.Menu{id: 1, name: "ベンチプレス"})
	menus.append(model.Menu{id: 2, name: "スクワット"})
	menus.append(model.Menu{id: 3, name: "デッドリフト"})
	db.Create(&menus)
	users := []model.User{}
	users.append(model.User{id: 1, name: "田中太郎", email: "tekitou@example.com", password: "tekitou", training_group_id: 1})
	users.append(model.User{id: 2, name: "佐藤花子", email: "tekitou@example.com", password: "tekitou", training_group_id: 1})
	users.append(model.User{id: 3, name: "山田次郎", email: "tekitou@example.com", password: "tekitou", training_group_id: 2})
	db.Create(&users)
	training_groups := []model.TrainingGroup{}
	training_groups.append(model.TrainingGroup{id: 1, name: "筋トレグループ1"})
	training_groups.append(model.TrainingGroup{id: 2, name: "筋トレグループ2"})
	db.Create(&training_groups)
	fmt.Println("Successfully Seeded")
}
