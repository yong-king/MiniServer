package models

import "time"

type Shelf struct {
	ID          int64     `json:"id" db:"id"`
	Theme       string    `json:"theme" db:"theme"`
	Size        int64     `json:"size" db:"size"`
	Create_time time.Time `json:"create_time" db:"create_time"`
	Update_time time.Time `json:"update_time" db:"update_time"`
}
