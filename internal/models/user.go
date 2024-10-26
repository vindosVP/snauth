package models

import "time"

type User struct {
	Id        int64
	Email     string
	HPassword string
	CreatedAt time.Time
	IsBanned  bool
	IsDeleted bool
	IsAdmin   bool
}
