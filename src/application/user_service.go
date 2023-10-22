package application

import (
	"context"
	"mime/multipart"

	"github.com/ryunosuke121/muscle-SNS/src/domain"
)

type IUserService interface {
	SignUp(ctx context.Context, id domain.UserID, name domain.UserName, email string) (*SignUpResponse, error)
	GetUserById(ctx context.Context, userId domain.UserID) (*GetUserPublic, error)
	GetUsersByIds(ctx context.Context, userId []domain.UserID) ([]*GetUserPublic, error)
	UpdateUserName(ctx context.Context, userId domain.UserID, name domain.UserName) (*GetUserPublic, error)
	UpdateUserGroup(ctx context.Context, userId domain.UserID, groupId domain.UserGroupID) (*GetUserPublic, error)
	UpdateUserImage(ctx context.Context, userId domain.UserID, file *multipart.FileHeader) (*GetUserPublic, error)
}

type userService struct {
	ur domain.IUserRepository
}

func NewUserService(ur domain.IUserRepository) IUserService {
	return &userService{ur}
}

func (us *userService) SignUp(ctx context.Context, id domain.UserID, name domain.UserName, email string) (*SignUpResponse, error) {
	newUser := &domain.User{
		ID:    id,
		Name:  name,
		Email: email,
	}

	if err := us.ur.CreateUser(ctx, newUser); err != nil {
		return nil, err
	}

	resUser := &SignUpResponse{
		ID:    newUser.ID.String(),
		Name:  newUser.Name.String(),
		Email: newUser.Email,
	}

	return resUser, nil
}

func (us *userService) GetUserById(ctx context.Context, userId domain.UserID) (*GetUserPublic, error) {
	users, err := us.ur.GetUsersByIds(ctx, []domain.UserID{userId})
	if err != nil {
		return nil, err
	}

	resUser := &GetUserPublic{
		ID:    users[0].ID.String(),
		Name:  users[0].Name.String(),
		Email: users[0].Email,
		UserGroup: &UserGroupPublic{
			ID:       uint(users[0].UserGroup.ID),
			Name:     users[0].UserGroup.Name,
			ImageUrl: users[0].UserGroup.ImageUrl,
		},
		AvatarUrl: users[0].AvatarUrl,
	}

	return resUser, nil
}

func (us *userService) GetUsersByIds(ctx context.Context, userIds []domain.UserID) ([]*GetUserPublic, error) {
	users, err := us.ur.GetUsersByIds(ctx, userIds)
	if err != nil {
		return nil, err
	}

	resUsers := make([]*GetUserPublic, len(users))
	for i, user := range users {
		resUsers[i] = &GetUserPublic{
			ID:    user.ID.String(),
			Name:  user.Name.String(),
			Email: user.Email,
			UserGroup: &UserGroupPublic{
				ID:       uint(user.UserGroup.ID),
				Name:     user.UserGroup.Name,
				ImageUrl: user.UserGroup.ImageUrl,
			},
			AvatarUrl: user.AvatarUrl,
		}
	}

	return resUsers, nil
}

func (us *userService) UpdateUserName(ctx context.Context, userId domain.UserID, name domain.UserName) (*GetUserPublic, error) {
	if err := us.ur.ChangeUserName(ctx, userId, name); err != nil {
		return nil, err
	}

	user, err := us.GetUserById(ctx, userId)
	if err != nil {
		return nil, err
	}

	resUser := &GetUserPublic{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		UserGroup: user.UserGroup,
		AvatarUrl: user.AvatarUrl,
	}

	return resUser, nil
}

func (us *userService) UpdateUserGroup(ctx context.Context, userId domain.UserID, groupId domain.UserGroupID) (*GetUserPublic, error) {
	if err := us.ur.ChangeUserGroup(ctx, userId, groupId); err != nil {
		return nil, err
	}

	user, err := us.GetUserById(ctx, userId)
	if err != nil {
		return nil, err
	}

	resUser := &GetUserPublic{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		UserGroup: user.UserGroup,
		AvatarUrl: user.AvatarUrl,
	}

	return resUser, nil
}

func (us *userService) UpdateUserImage(ctx context.Context, userId domain.UserID, file *multipart.FileHeader) (*GetUserPublic, error) {
	err := us.ur.ChangeUserImage(ctx, userId, file)
	if err != nil {
		return nil, err
	}

	user, err := us.GetUserById(ctx, userId)
	if err != nil {
		return nil, err
	}

	resUser := &GetUserPublic{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		UserGroup: user.UserGroup,
		AvatarUrl: user.AvatarUrl,
	}

	return resUser, nil
}
