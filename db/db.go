package db

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"os"
)

// DBに接続する
func NewDB() *gorm.DB {
	url := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", os.Getenv("MYSQL_USER"),
		os.Getenv("MYSQL_PASSWORD"), os.Getenv("MYSQL_HOST"),
		os.Getenv("MYSQL_PORT"), os.Getenv("MYSQL_DB"))
	fmt.Println(url)
	db, err := gorm.Open(mysql.Open(url), &gorm.Config{})
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("Connceted")
	return db
}

// DBを閉じる
func CloseDB(db *gorm.DB) {
	sqlDB, _ := db.DB()
	if err := sqlDB.Close(); err != nil {
		log.Fatalln(err)
	}
}
