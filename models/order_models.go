package models

import (
	"time"

	"gorm.io/gorm"
)

type Order struct {
	gorm.Model
	UserID      int       `json:"userId"` // Foreign key to User
	OrderDate   time.Time `json:"orderDate"`
	TotalAmount float64   `json:"totalAmount"`
	Status      string    `json:"status"`
}

// OrderResponse represents the structure of the entire response
type OrderResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    struct {
		Orders []Order `json:"orders"`
	} `json:"data"`
}
