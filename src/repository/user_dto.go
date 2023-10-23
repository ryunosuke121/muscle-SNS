package repository

import (
	"time"
)

type User struct {
	ID          string      `json:"id" gorm:"primaryKey;size:64"`
	Name        string      `json:"name" gorm:"not null" validate:"required"`
	Email       string      `json:"email" gorm:"unique" validate:"required,email"`
	UserGroupID uint        `json:"user_group_id" gorm:"default:1"`
	UserGroup   UserGroup   `json:"user_group" gorm:"foreignKey:UserGroupID"`
	ImageUrl    string      `json:"image_url"`
	Posts       []*Post     `json:"posts" gorm:"foreignKey:UserID"`
	Trainings   []*Training `json:"trainings" gorm:"foreignKey:UserID"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
}

type UserGroup struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name" gorm:"not null" validate:"required"`
	ImageUrl  string    `json:"image_url"`
	Users     []*User   `json:"users" gorm:"foreignKey:UserGroupID"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
