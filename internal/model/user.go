package model

import "time"

type User struct {
	ID          int       `json:"id"`
	TelegramID  int64     `json:"telegram_id"`
	Username    string    `json:"username,omitempty"`
	FirstName   string    `json:"first_name,omitempty"`
	LastName    string    `json:"last_name,omitempty"`
	Role        string    `json:"role"`
	IsBlocked   bool      `json:"is_blocked"`
	Phone       string    `json:"phone,omitempty"`
	CompanyName string    `json:"company_name,omitempty"`
	Description string    `json:"description,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

const (
	RoleClient     = "client"
	RoleLandlord   = "landlord"
	RoleAdmin      = "admin"
	RoleSuperAdmin = "super_admin"
)
