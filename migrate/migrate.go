package main

import (
	"fmt"
	"github.com/ryunosuke121/muscle-SNS/db"
	"github.com/ryunosuke121/muscle-SNS/model"
)

func main() {
	dbConn := db.NewDB()
	defer fmt.Println("Successfully Migrated")
	defer db.CloseDB(dbConn)
	dbConn.AutoMigrate(&model.TrainingGroup{}, &model.User{}, &model.Training{}, &model.Menu{}, &model.Post{})
}
