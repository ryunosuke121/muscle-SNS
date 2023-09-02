package model

type User struct {
	ID              uint          `json:"id" gorm:"primaryKey"`
	Name            string        `json:"name" gorm:"not null" validate:"required"`
	Email           string        `json:"email" gorm:"unique" validate:"required,email"`
	Password        string        `json:"password" validate:"required,min=4,max=32"`
	TrainingGroupID uint          `json:"training_group_id"`
	TrainingGroup   TrainingGroup `json:"training_group"`
	ImageUrl        string        `json:"image_url"`
	Posts           []Post        `json:"posts"`
	Trainings       []Training    `json:"trainings"`
}

type UserResponse struct {
	ID              uint   `json:"id"`
	Name            string `json:"name"`
	Email           string `json:"email"`
	TrainingGroupID uint   `json:"training_group_id"`
	ImageUrl        string `json:"image_url"`
}

type TrainingGroup struct {
	ID       uint   `json:"id" gorm:"primaryKey"`
	Name     string `json:"name" gorm:"not null"`
	ImageUrl string `json:"image_url"`
	Users    []User `json:"users"`
}
