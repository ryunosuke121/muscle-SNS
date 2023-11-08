package domain

import "errors"

type User struct {
	ID        UserID
	Name      UserName
	Email     string
	AvatarUrl string
	UserGroup *UserGroup
	Posts     []*Post
	Trainings []*Training
}

type UserGroup struct {
	ID       UserGroupID
	Name     string
	ImageUrl string
	Users    []*User
}

type (
	UserID      string
	UserName    string
	UserGroupID uint
)

func (uid UserID) String() string {
	return string(uid)
}

func (un UserName) String() string {
	return string(un)
}

var (
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrUserNotFound       = errors.New("user not found")
	ErrUserInfoNotChanged = errors.New("user info not changed")
)
