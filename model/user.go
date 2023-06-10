package model

import "time"

type User struct {
	ID       uint   `json:"id" gorm:"primaryKey"`
	NAME     string `json:"name" gorm:"not null"`
	Email    string `json:"email" gorm:"unique"`
	Password string `json:"password"`
	GroupId  uint   `json:"group_id"`
	ImageUrl string `json:"image_url"`
}

type UserResponse struct {
	ID    uint   `json:"id" gorm:"primaryKey"`
	Email string `json:"email" gorm:"unique"`
}
