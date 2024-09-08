package models

import (
	"time"
)

type Order struct {
	ID           uint      `json:"id" gorm:"primarykey"`
	CustomerName string    `json:"customer_name"`
	OrderedAt    time.Time `json:"ordered_at"`
	Items        []Item    `json:"items" gorm:"foreignKey:OrderID"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
