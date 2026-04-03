package model

import "time"

type PaymentCard struct {
	ID         int       `json:"id"`
	AdminID    int       `json:"admin_id"`
	Label      string    `json:"label"`
	CardNumber string    `json:"card_number"`
	Cardholder string    `json:"cardholder,omitempty"`
	IsActive   bool      `json:"is_active"`
	SortOrder  int       `json:"sort_order"`
	CreatedAt  time.Time `json:"created_at"`
}
