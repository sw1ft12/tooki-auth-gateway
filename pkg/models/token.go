package models

import "time"

type RefreshToken struct {
	Token          string
	ExpirationTime time.Time
	UserId         string
	UserAgent      string
}
