package domain

import (
	"context"
	"mime/multipart"
)

type IUserRepository interface {
	CreateUser(ctx context.Context, user *User) error                                     // ユーザーを作成する
	GetUsersByIds(ctx context.Context, userIds []UserID) ([]*User, error)                 // IDのリストからユーザーを取得する
	ChangeUserName(ctx context.Context, userId UserID, userName UserName) error           // ユーザーの名前を更新する
	ChangeUserGroup(ctx context.Context, userId UserID, groupId UserGroupID) error        // ユーザーのグループを更新する
	GetUserByEmail(ctx context.Context, email string) (*User, error)                      // メールアドレスからユーザーを取得する
	ChangeUserImage(ctx context.Context, userId UserID, file *multipart.FileHeader) error // ユーザーの画像を更新する
	GetUsersInGroup(ctx context.Context, groupId UserGroupID) ([]*User, error)            // グループに所属するユーザーを取得する
}

type IPostRepository interface {
	GetPostsByIds(ctx context.Context, ids []PostID) ([]*Post, error)                                                                    // 投稿を取得する
	GetUserPosts(ctx context.Context, userId UserID) ([]*Post, error)                                                                    // ユーザーの投稿を取得する
	CreatePost(ctx context.Context, post *Post) (*Post, error)                                                                           // 投稿を作成する
	DeletePost(ctx context.Context, postId PostID) error                                                                                 // 投稿を削除する
	GetGroupPosts(ctx context.Context, groupId UserGroupID) ([]*Post, error)                                                             // グループの投稿を取得する
	GetTrainingsByIds(ctx context.Context, ids []TrainingID) ([]*Training, error)                                                        // トレーニングを取得する
	GetUserTrainings(ctx context.Context, userId UserID) ([]*Training, error)                                                            // ユーザーのトレーニングを取得する
	GetUsersTotalWeightByMenuInMonth(ctx context.Context, userIds []UserID, menuId MenuID, year int, month int) (map[UserID]uint, error) // メニュー別のユーザーの総重量を取得する
	GetUsersTotalWeightInMonth(ctx context.Context, userIds []UserID, year int, month int) (map[UserID]uint, error)                      // ユーザー別の総重量を取得する
	SavePostImage(ctx context.Context, file *multipart.FileHeader) (fileName string, err error)                                          // 投稿の画像を保存する
}

type IMenuRepository interface {
	GetMenuById(ctx context.Context, id MenuID) (Menu, error)
}
