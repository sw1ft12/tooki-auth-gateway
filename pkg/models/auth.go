package models

import "time"

type RegisterUserDto struct {
	Email    string `json:"email" db:"email" validate:"required" example:"test@email.com"`
	Login    string `json:"login" db:"login" validate:"required" example:"sw1ft12"`
	Password string `json:"password" db:"password" validate:"required" example:"assag23214"`
	Name     string `json:"name" db:"name" validate:"required" example:"Артёмчик Zиновьев"`
	Age      int    `json:"age" db:"age" validate:"required" example:"5"`
	Gender   string `json:"gender" enums:"Male, Female, Unknown" db:"gender" validate:"required, oneof=Male Female Unknown" example:"Female"`
}

type RegisterResponse struct {
	Id        string    `json:"id" db:"id" example:"8cbabbe9-5fff-4dbe-a77e-104bf4e63dbe"`
	Email     string    `json:"email" db:"email" example:"test@gmail.com"`
	Name      string    `json:"name" db:"name" example:"Зиновьев Артём"`
	Role      string    `json:"role" enums:"User, Admin, SuperAdmin" db:"role" example:"USER"`
	CreatedAt time.Time `json:"created_at" db:"created_at" example:"2024-03-02"`
	Verified  bool      `json:"verified" db:"verified"`
	Banned    bool      `json:"banned" db:"banned"`
}

type LoginUserDto struct {
	Login    string `json:"login" db:"login" validate:"required"`
	Password string `json:"password" db:"password" validate:"required"`
}

type LoginResponse struct {
	AccessToken string `json:"token" db:"token"`
	Id          string `json:"id" db:"id"`
	Name        string `json:"name" db:"name"`
	Email       string `json:"email" db:"email"`
	Gender      string `json:"gender" enums:"Male, Female, Unknown" db:"gender"`
	Role        string `json:"role" enums:"User, Admin, SuperAdmin" db:"role"`
}
