package domain

import "time"

type RefreshSession struct {
	Id        int
	UserId    int
	Token     string
	ExpiresAt time.Time
}
