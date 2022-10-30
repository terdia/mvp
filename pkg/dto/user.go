package dto

import (
	"time"
)

type CreateUserRequest struct {
	Username string `json:"username"` // username
	Role     string `json:"role"`     // seller|buyer
	Password string `json:"password"` // minimum 6 bytes maximum 72 bytes
}

type UserResponse struct {
	User APIUser `json:"user"`
}

type APIUser struct {
	ID        int64     `json:"id"`
	Username  string    `json:"username"`
	Role      string    `json:"role"`
	Deposit   int       `json:"deposit"`
	CreatedAt time.Time `json:"created_at"`
}
