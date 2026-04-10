package dto

import "time"

type UserRegisterRes struct {
	Username string `json:"username"`
}

type UserLoginRes struct {
	Username string `json:"username"`
	Token    string `json:"token"`
}

type UserProfileRes struct {
	ID        int64     `json:"id"`
	Email     string    `json:"email"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
