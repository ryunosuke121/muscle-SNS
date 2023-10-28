package domain

import (
	"errors"
	"fmt"
	"time"
)

type Post struct {
	ID        PostID
	UserID    UserID
	Training  *Training
	Comment   string
	CreatedAt time.Time
	ImageName string
}

type (
	PostID uint
)

func (pid PostID) String() string {
	return fmt.Sprintf("%d", pid)
}

var (
	ErrForbidden = errors.New("this action is forbidden")
)
