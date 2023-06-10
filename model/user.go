package model

type User struct {
	ID       uint   `json:"id" gorm:"primaryKey"`
	NAME     string `json:"name" gorm:"not null" validate:"required"`
	Email    string `json:"email" gorm:"unique" validate:"required,email"`
	Password string `json:"password" validate:"required,min=4,max=32"`
	GroupId  uint   `json:"group_id"`
	ImageUrl string `json:"image_url"`
}
