package models

import "time"

type Item struct {
	ID          uint      `json:"id" gorm:"primarykey"`
	OrderID     uint      `json:"order_id"`
	ItemCode    string    `json:"item_code"`
	Description string    `json:"description"`
	Quantity    uint      `json:"quantity"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
