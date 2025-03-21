package models

import (
	"time"
)

type Book struct {
	ID         int
	Title      string
	Year       int
	Author     string
	Price      int
	Amount     int
	CategoryID int
	CreatedAt  time.Time
	UpdatedAt  *time.Time
}
