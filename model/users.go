package model

import "time"

type UsersModel struct {
	Id        int        `json:"id" db:"id"`
	Username  string     `json:"username" db:"username"`
	Password  string     `json:"password" db:"password"`
	CreatedAt *time.Time `json:"created_at" db:"created_at"`
	UpdatedAt *time.Time `json:"updated_at" db:"updated_at"`
}

type RequestUsersModel struct {
	Username string `json:"username" db:"username" validate:"required,alphanum"`
	Password string `json:"password" db:"password" validate:"required,alphanum"`
}
