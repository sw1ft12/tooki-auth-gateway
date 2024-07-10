package models

import "time"

type RefreshToken struct {
	Token     string    `json:"token" db:"token"`
	ExpiresIn time.Time `json:"expires_in" db:"expires_in"`
	UserId    string    `json:"user_id" db:"user_id"`
	UserAgent string    `json:"user_agent" db:"user_agent"`
}
