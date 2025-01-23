package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/osag1e/user-account-auth-system/internal/models"
)

func TestHandleCreateUser(t *testing.T) {
	mockRepo := &MockUserRepository{}
	userHandler := NewUserHandler(mockRepo)

	userParams := models.CreateUserParams{
		Email:    "osagie@gg.uk",
		UserName: "Osagie",
		Password: "Password",
	}

	jsonData, err := json.Marshal(userParams)
	if err != nil {
		t.Fatalf("Failed to marshal user params: %v", err)
	}

	req, err := http.NewRequest("POST", "/create-user", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	resRecorder := httptest.NewRecorder()

	userHandler.HandleCreateUser(resRecorder, req)

	if resRecorder.Code != http.StatusCreated {
		t.Errorf("Expected status code %d, got %d", http.StatusCreated, resRecorder.Code)
	}

	var createdUser models.User
	err = json.NewDecoder(resRecorder.Body).Decode(&createdUser)
	if err != nil {
		t.Fatalf("Failed to decode response body: %v", err)
	}
}
