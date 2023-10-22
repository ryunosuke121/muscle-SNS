package main

import (
	"fmt"

	"github.com/ryunosuke121/muscle-SNS/src/repository"
	"github.com/ryunosuke121/muscle-SNS/src/repository/db"
)

func main() {
	fmt.Println("Start Seeding")
	db := db.NewDB()
	menus := []repository.Menu{}
	menus = append(menus, repository.Menu{ID: 1, Name: "ベンチプレス"})
	menus = append(menus, repository.Menu{ID: 2, Name: "スクワット"})
	menus = append(menus, repository.Menu{ID: 3, Name: "デッドリフト"})
	result := db.Create(&menus)
	if result.Error != nil {
		fmt.Println(result.Error)
		return
	}

	userGroups := []repository.UserGroup{}
	userGroups = append(userGroups, repository.UserGroup{ID: 1, Name: "筋トレ好きあつまれ", ImageUrl: ""})
	userGroups = append(userGroups, repository.UserGroup{ID: 2, Name: "みんなで筋トレ", ImageUrl: ""})
	userGroups = append(userGroups, repository.UserGroup{ID: 3, Name: "マッチョになろうぜ。", ImageUrl: ""})
	userGroups = append(userGroups, repository.UserGroup{ID: 4, Name: "ダンベル好きあつまれ", ImageUrl: ""})
	result = db.Create(&userGroups)
	if result.Error != nil {
		fmt.Println(result.Error)
		return
	}

	fmt.Println("Successfully Seeded")
}
