package models

import "time"

type User struct {
	Id        int
	Email     string
	Password  string
	IsAdmin   bool
	CreatedAt time.Time
	UpdatedAt *time.Time
}
