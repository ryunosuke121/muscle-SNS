package domain

import "time"

type Training struct {
	ID        TrainingID
	UserID    UserID
	Menu      *Menu
	Times     uint
	Weight    uint
	Sets      uint
	CreatedAt time.Time
}

type Menu struct {
	ID   MenuID
	Name string
}

type (
	TrainingID uint
	MenuID     uint
)
