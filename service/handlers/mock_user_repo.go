package handlers

import (
	"github.com/OsagieDG/user-account-auth-system/internal/models"
	"github.com/google/uuid"
)

type MockUserRepository struct{}

func (m *MockUserRepository) InsertUser(user *models.User) (*models.User, error) {
	return &models.User{
		ID:       uuid.New(),
		Email:    user.Email,
		UserName: user.UserName,
	}, nil
}

func (m *MockUserRepository) DeleteUserByID(userID uuid.UUID) error {
	return nil
}

func (m *MockUserRepository) GetUserByEmail(email string) (*models.User, error) {
	return &models.User{
		Email: email,
	}, nil
}

func (m *MockUserRepository) GetUserByID(userID uuid.UUID) (*models.User, error) {
	return &models.User{
		ID: userID,
	}, nil
}

func (m *MockUserRepository) GetUsers() ([]models.User, error) {
	users := []models.User{
		{
			ID:       uuid.New(),
			Email:    "user1@example.com",
			UserName: "User1",
		},
		{
			ID:       uuid.New(),
			Email:    "user2@example.com",
			UserName: "User2",
		},
	}
	return users, nil
}

func (m *MockUserRepository) UpdateUserByID(userID uuid.UUID, params models.UpdateUserParams) (*models.User, error) {
	updatedUser := &models.User{
		ID:       userID,
		UserName: params.UserName,
	}
	return updatedUser, nil
}
