package repository

import (
	"time"
)

// トレーニング記録
type Training struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UserID    string    `json:"user_id" gorm:"not null " validate:"required"`
	User      *User     `json:"user" gorm:"foreignKey:UserID;OnUpdate:CASCADE,OnDelete:CASCADE;"`
	MenuID    uint      `json:"menu_id" gorm:"not null" validate:"required"`
	Menu      *Menu     `json:"menu" gorm:"foreignKey:MenuID;OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Times     uint      `json:"times"`
	Weight    uint      `json:"weight"`
	Sets      uint      `json:"sets"`
	CreatedAt time.Time `json:"created_at"`
}

// メニュー
type Menu struct {
	ID   uint   `json:"id" gorm:"primaryKey"`
	Name string `json:"name" gorm:"not null" validate:"required"`
}

// 投稿
type Post struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	UserID     string    `json:"user_id" gorm:"not null" validate:"required"`
	User       *User     `json:"user" gorm:"foreignKey:UserID;OnUpdate:CASCADE,OnDelete:CASCADE;"`
	TrainingID uint      `json:"training_id" gorm:"not null" validate:"required"`
	Training   *Training `json:"training" gorm:"foreignKey:TrainingID;OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Comment    string    `json:"comment" gorm:"not null"`
	CreatedAt  time.Time `json:"created_at" gorm:"not null"`
	ImageUrl   string    `json:"image_url"`
}
