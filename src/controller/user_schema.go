package controller

type SignUpRequestSchema struct {
	Name string `json:"name" validate:"required"`
}

type UpdateUserNameRequestSchema struct {
	Name string `json:"name" validate:"required"`
}

type UpdateUserGroupRequestSchema struct {
	GroupID uint `json:"group_id" validate:"required"`
}

type UserResponse struct {
	ID              uint   `json:"id"`
	Name            string `json:"name"`
	Email           string `json:"email"`
	TrainingGroupID uint   `json:"training_group_id"`
	ImageUrl        string `json:"image_url"`
}
