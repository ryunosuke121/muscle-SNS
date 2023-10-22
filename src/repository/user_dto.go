package repository

type User struct {
	ID          string     `json:"id" gorm:"primaryKey;size:64"`
	Name        string     `json:"name" gorm:"not null" validate:"required"`
	Email       string     `json:"email" gorm:"unique" validate:"required,email"`
	UserGroupID uint       `json:"user_group_id"`
	UserGroup   UserGroup  `json:"user_group"`
	ImageUrl    string     `json:"image_url"`
	Posts       []Post     `json:"posts"`
	Trainings   []Training `json:"trainings"`
}

type UserGroup struct {
	ID       uint   `json:"id" gorm:"primaryKey"`
	Name     string `json:"name" gorm:"not null" validate:"required"`
	ImageUrl string `json:"image_url"`
	Users    []User `json:"users"`
}
