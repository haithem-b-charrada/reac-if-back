package models

import "time"

type User struct {
	ID                       int       `json:"id" gorm:"primary_key"`
	FirstName                string    `json:"first_name"`
	LastName                 string    `json:"last_name"`
	Email                    string    `json:"email" gorm:"unique"`
	CreatedAt                time.Time `json:"created_at"`
	Password                 string    `json:"password"`
	Phone                    string    `json:"phone"`
	Verified                 bool      `json:"verified"`
	Roles                    string    `json:"roles"`
	Gender                   string    `json:"gender"`
	Image                    string    `json:"image"`
	ResetPasswordToken       string    `json:"reset_password_token"`
	ResetPasswordExpiredDate time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"reset_pass_expired"`
	Conformation             bool      `json:"conformation"`
	WarehouseID              int       `json:"warehouseID"`
}
