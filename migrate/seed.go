package migrate

import (
	"fmt"

	"github.com/ryunosuke121/muscle-SNS/db"
	"github.com/ryunosuke121/muscle-SNS/model"
)

func seed() {
	db := db.NewDB()
	menus := []model.Menu{}
	menus = append(menus, model.Menu{ID: 1, Name: "ベンチプレス"})
	menus = append(menus, model.Menu{ID: 2, Name: "スクワット"})
	menus = append(menus, model.Menu{ID: 3, Name: "デッドリフト"})
	db.Create(&menus)
	users := []model.User{}
	users = append(users, model.User{ID: 1, Name: "田中太郎", Email: "tekitou@example.com", Password: "tekitou", TrainingGroupID: 1})
	users = append(users, model.User{ID: 2, Name: "佐藤花子", Email: "tekitou@example.com", Password: "tekitou", TrainingGroupID: 1})
	users = append(users, model.User{ID: 3, Name: "山田次郎", Email: "tekitou@example.com", Password: "tekitou", TrainingGroupID: 2})
	db.Create(&users)
	training_groups := []model.TrainingGroup{}
	training_groups = append(training_groups, model.TrainingGroup{ID: 1, Name: "筋トレグループ1"})
	training_groups = append(training_groups, model.TrainingGroup{ID: 2, Name: "筋トレグループ2"})
	db.Create(&training_groups)
	fmt.Println("Successfully Seeded")
}
