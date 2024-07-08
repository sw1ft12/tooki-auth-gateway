package models

import "time"

type User struct {
	Id        string    `json:"id" db:"id"`
	Email     string    `json:"email" db:"email"`
	Login     string    `json:"login" db:"login"`
	Password  string    `json:"password" db:"password"`
	Name      string    `json:"name" db:"name"`
	Age       int       `json:"age" db:"age"`
	Gender    string    `json:"gender" db:"gender"`
	Role      string    `json:"role" db:"role"`
	Avatar    string    `json:"avatar" db:"avatar"`
	Verified  bool      `json:"verified" db:"verified"`
	Banned    bool      `json:"banned" db:"banned"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}
