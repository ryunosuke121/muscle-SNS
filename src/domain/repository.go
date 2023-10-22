package domain

import (
	"context"
	"mime/multipart"
)

type IUserRepository interface {
	CreateUser(ctx context.Context, user *User) error
	GetUsersByIds(ctx context.Context, userIds []UserID) ([]*User, error)
	ChangeUserName(ctx context.Context, userId UserID, userName UserName) error
	ChangeUserGroup(ctx context.Context, userId UserID, groupId UserGroupID) error
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	GetUserImageUrlsByIds(ctx context.Context, userIds []UserID) (map[UserID]string, error)
	ChangeUserImage(ctx context.Context, userId UserID, file *multipart.FileHeader) error
}

type IPostRepository interface {
	GetPostsByIds(ctx context.Context, ids []PostID) ([]*Post, error)
	GetUserPosts(ctx context.Context, userId UserID) ([]*Post, error)
	CreatePost(ctx context.Context, post *Post) (*Post, error)
	DeletePost(ctx context.Context, postId PostID) error
	GetGroupPosts(ctx context.Context, groupId UserGroupID) ([]*Post, error)
	GetTrainingsByIds(ctx context.Context, ids []TrainingID) ([]*Training, error)
	GetUserTrainings(ctx context.Context, userId UserID) ([]*Training, error)
	SavePostImage(ctx context.Context, file *multipart.FileHeader) (fileName string, err error)
}
