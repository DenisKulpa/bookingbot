package model

import "time"

const (
	BookingStatusPendingApproval = "pending_approval"
	BookingStatusApproved        = "approved"
	BookingStatusPaymentClaimed  = "payment_claimed"
	BookingStatusConfirmed       = "confirmed"
	BookingStatusRejected        = "rejected"
	BookingStatusCancelled       = "cancelled"
)

type Booking struct {
	ID          int       `json:"id"`
	ApartmentID int       `json:"apartment_id"`
	ClientID    int       `json:"client_id"`
	CheckIn     time.Time `json:"check_in"`
	CheckOut    time.Time `json:"check_out"`
	GuestsCount int       `json:"guests_count"`
	TotalPrice  float64   `json:"total_price"`
	Status      string    `json:"status"`
	AdminNote   string    `json:"admin_note,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
