package main

import (
	"fmt"

	"github.com/ryunosuke121/muscle-SNS/db"
	"github.com/ryunosuke121/muscle-SNS/model"
)

func main() {
	db := db.NewDB()
	menus := []model.Menu{}
	menus = append(menus, model.Menu{ID: 1, Name: "ベンチプレス"})
	menus = append(menus, model.Menu{ID: 2, Name: "スクワット"})
	menus = append(menus, model.Menu{ID: 3, Name: "デッドリフト"})
	db.Create(&menus)
	training_groups := []model.TrainingGroup{}
	training_groups = append(training_groups, model.TrainingGroup{ID: 1, Name: "筋トレ好きあつまれ", ImageUrl: "group_image/istockphoto-1209027669-1024x1024.jpeg"})
	training_groups = append(training_groups, model.TrainingGroup{ID: 2, Name: "みんなで筋トレ", ImageUrl: "group_image/macho_man.png"})
	training_groups = append(training_groups, model.TrainingGroup{ID: 3, Name: "マッチョになろうぜ。", ImageUrl: "group_image/o0320048214003980269-1-1.jpeg"})
	training_groups = append(training_groups, model.TrainingGroup{ID: 4, Name: "ダンベル好きあつまれ", ImageUrl: "group_image/c67d81caed121f7d0e224774b3b7e6fd-1024x683.jpeg"})
	db.Create(&training_groups)
	users := []model.User{}
	users = append(users, model.User{ID: 1, Name: "田中太郎", Email: "tanaka@example.com", Password: "tekitou", TrainingGroupID: 1})
	users = append(users, model.User{ID: 2, Name: "佐藤花子", Email: "satou@example.com", Password: "tekitou", TrainingGroupID: 1})
	users = append(users, model.User{ID: 3, Name: "山田次郎", Email: "yamada@example.com", Password: "tekitou", TrainingGroupID: 2})
	db.Create(&users)

	fmt.Println("Successfully Seeded")
}
