package query

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/OsagieDG/user-account-auth-system/internal/models"
	"github.com/google/uuid"
)

type UserRepository interface {
	InsertUser(user *models.User) (*models.User, error)
	GetUserByEmail(email string) (*models.User, error)
	GetUserByID(userID uuid.UUID) (*models.User, error)
	GetUsers() ([]models.User, error)
	UpdateUserByID(userID uuid.UUID, params models.UpdateUserParams) (*models.User, error)
	DeleteUserByID(userID uuid.UUID) error
}

type UserPostgresRepository struct {
	DB *sql.DB
}

func NewPostgresUserRepository(db *sql.DB) UserRepository {
	return &UserPostgresRepository{DB: db}
}

func (ur *UserPostgresRepository) InsertUser(user *models.User) (*models.User, error) {
	userID := models.NewUUID()

	_, err := ur.DB.Exec(
		"INSERT INTO auth.users (id, username, email, encryptedpassword, isadmin) VALUES ($1, $2, $3, $4, $5)",
		userID, user.UserName, user.Email, user.EncryptedPassword, false,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to insert user into database: %w", err)
	}

	return user, nil
}

func (ur *UserPostgresRepository) UpdateUserByID(userID uuid.UUID, params models.UpdateUserParams) (*models.User, error) {
	tx, err := ur.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
		}
	}()

	_, err = tx.Exec(
		"UPDATE auth.users SET username = $1 WHERE id = $2",
		params.UserName, userID,
	)
	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return ur.GetUserByID(userID)
}

func (ur *UserPostgresRepository) GetUserByEmail(email string) (*models.User, error) {
	row := ur.DB.QueryRow("SELECT id, username, email, encryptedpassword, isadmin FROM auth.users WHERE email = $1", email)

	var user models.User
	err := row.Scan(&user.ID, &user.UserName, &user.Email, &user.EncryptedPassword, &user.IsAdmin)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return &user, nil
}

func (ur *UserPostgresRepository) GetUserByID(userID uuid.UUID) (*models.User, error) {
	row := ur.DB.QueryRow("SELECT id, username, email, isadmin FROM auth.users WHERE id = $1", userID)

	var user models.User
	if err := row.Scan(&user.ID, &user.UserName, &user.Email, &user.IsAdmin); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user with ID %s not found", userID.String())
		}
		return nil, err
	}

	return &user, nil
}

func (ur *UserPostgresRepository) GetUsers() ([]models.User, error) {
	rows, err := ur.DB.Query("SELECT id, username, email, isadmin FROM auth.users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.UserName, &user.Email, &user.IsAdmin); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func (ur *UserPostgresRepository) DeleteUserByID(userID uuid.UUID) error {
	tx, err := ur.DB.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
		}
	}()

	_, err = tx.Exec("DELETE FROM auth.sessions WHERE userid = $1", userID)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	_, err = tx.Exec("DELETE FROM auth.users WHERE id = $1", userID)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	return tx.Commit()
}
