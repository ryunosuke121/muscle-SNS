package main

import (
	"fmt"
	"github.com/ryunosuke121/go-rest/db"
	"github.com/ryunosuke121/go-rest/model"
)

func main() {
	dbConn := db.NewDB()
	defer fmt.Println("Successfully Migrated")
	defer db.CloseDB(dbConn)
	dbConn.AutoMigrate(&model.User{}, &model.Task{})
}
