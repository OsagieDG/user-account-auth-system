package models

import "github.com/google/uuid"

type User struct {
	ID                uuid.UUID `json:"id"`
	UserName          string    `json:"username"`
	Email             string    `json:"email"`
	EncryptedPassword string    `json:"encryptedpassword"`
	IsAdmin           bool      `json:"isadmin"`
}

/*
	type UserID struct {
		ID uuid.UUID `json:"id"`
	}
*/
func NewUUID() uuid.UUID {
	return uuid.New()
}
