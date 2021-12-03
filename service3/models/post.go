package models

import "time"

type Post struct {
	ID        int       `json:"id" gorm:"primary_key"`
	Firstname string    `json:"firstname"`
	Lastname  string    `json:"lastname"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedAt time.Time `json:"created_at"`
	Context   string    `json:"context"`
	Status    string    `json:"status"`
}
