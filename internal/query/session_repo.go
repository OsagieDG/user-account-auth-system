package query

import (
	"database/sql"
	"fmt"

	"github.com/OsagieDG/user-account-auth-system/internal/models"
)

type SesionRepository interface {
	CreateSession(session models.Session) error
	GetSessionByToken(token string) (*models.Session, error)
	DeleteSession(token string) error
}

type SessionPostgresRepository struct {
	DB *sql.DB
}

func NewPostgresSesssionRepository(db *sql.DB) SesionRepository {
	return &SessionPostgresRepository{DB: db}
}

func (sr *SessionPostgresRepository) CreateSession(session models.Session) error {
	_, err := sr.DB.Exec(
		"INSERT INTO auth.sessions (userid, token, expiresat) VALUES ($1, $2, $3)",
		session.UserID, session.Token, session.ExpiresAt,
	)
	return err
}

func (sr *SessionPostgresRepository) GetSessionByToken(token string) (*models.Session, error) {
	row := sr.DB.QueryRow("SELECT id, userid, token, expiresat FROM auth.sessions WHERE token = $1", token)
	var session models.Session
	if err := row.Scan(&session.ID, &session.UserID, &session.Token, &session.ExpiresAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("session not found")
		}
		return nil, err
	}
	return &session, nil
}

func (sr *SessionPostgresRepository) DeleteSession(token string) error {
	_, err := sr.DB.Exec("DELETE FROM auth.sessions WHERE token = $1", token)
	return err
}
