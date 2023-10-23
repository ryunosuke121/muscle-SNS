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
