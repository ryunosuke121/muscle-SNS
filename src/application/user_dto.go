package application

import (
	"mime/multipart"

	"github.com/ryunosuke121/muscle-SNS/src/domain"
)

type SignUpResponse struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type GetUserPublic struct {
	ID        string           `json:"id"`
	Name      string           `json:"name"`
	Email     string           `json:"email"`
	UserGroup *UserGroupPublic `json:"user_group"`
	AvatarUrl string           `json:"avatar_url"`
}

type GetUserPublicResponse struct {
	Users []GetUserPublic `json:"users"`
}

type UserGroupPublic struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	ImageUrl string `json:"image_url"`
}

type GetUserGroupResponse struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	ImageUrl string `json:"image_url"`
}

type PostPublic struct {
	ID        string          `json:"id"`
	UserID    string          `json:"user_id"`
	Training  *TrainingPublic `json:"training"`
	Comment   string          `json:"comment"`
	CreatedAt string          `json:"created_at"`
	ImageUrl  string          `json:"image_url"`
}

type TrainingPublic struct {
	ID        uint   `json:"id"`
	UserID    string `json:"user_id"`
	Menu      *Menu  `json:"menu"`
	Times     uint   `json:"times"`
	Weight    uint   `json:"weight"`
	Sets      uint   `json:"sets"`
	CreatedAt string `json:"created_at"`
}

type Menu struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type CreatePostRequest struct {
	UserID   domain.UserID         `json:"user_id"`
	Comment  string                `json:"comment"`
	Image    *multipart.FileHeader `json:"image"`
	Training *TrainingRequest      `json:"training"`
}

type TrainingRequest struct {
	MenuID domain.MenuID `json:"menu_id"`
	Times  uint          `json:"times"`
	Weight uint          `json:"weight"`
	Sets   uint          `json:"sets"`
}
