package models

import (
	"time"

	"github.com/google/uuid"
)

type Session struct {
	ID        int       `json:"id"`
	UserID    uuid.UUID `json:"userid"`
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expiresat"`
}
