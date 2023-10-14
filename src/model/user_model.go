package model

type User struct {
	ID              string        `json:"id" gorm:"primaryKey;size:64"`
	Name            string        `json:"name" gorm:"not null" validate:"required"`
	Email           string        `json:"email" gorm:"unique" validate:"required,email"`
	TrainingGroupID uint          `json:"training_group_id"`
	TrainingGroup   TrainingGroup `json:"training_group"`
	ImageUrl        string        `json:"image_url"`
	Posts           []Post        `json:"posts"`
	Trainings       []Training    `json:"trainings"`
}

type UserResponse struct {
	ID              string `json:"id"`
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
