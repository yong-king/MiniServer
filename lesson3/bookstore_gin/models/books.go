package models

import "time"

type Book struct {
	ID          int64     `json:"id" db:"id"`
	ShelfID     int64     `json:"shelf_id" db:"shelf_id"`
	Title       string    `json:"title" db:"title"`
	Author      string    `json:"author" db:"author"`
	Create_time time.Time `json:"create_time" db:"create_time"`
	Update_time time.Time `json:"update_time" db:"update_time"`
}
