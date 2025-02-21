package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/OsagieDG/user-account-auth-system/internal/models"
	"github.com/OsagieDG/user-account-auth-system/internal/query"
	"github.com/google/uuid"
)

// ContextKey is a type used for keys in context to avoid collisions.
type ContextKey string

const ContextKeyUserID ContextKey = "userID"

type SessionHandler struct {
	DB                *sql.DB
	userRepository    query.UserRepository
	sessionRepository query.SesionRepository
}

func NewSessionHandler(db *sql.DB, userRepository query.UserRepository, sessionRepository query.SesionRepository) *SessionHandler {
	return &SessionHandler{
		DB:                db,
		userRepository:    userRepository,
		sessionRepository: sessionRepository,
	}
}

func (s *SessionHandler) GenerateSession(userID uuid.UUID) (*models.Session, error) {
	token, err := GenerateRandomToken(64)
	if err != nil {
		return nil, err
	}
	expiresAt := time.Now().Add(time.Hour * 1) // Set session expiration time 1 hour
	session := models.Session{
		UserID:    userID,
		Token:     token,
		ExpiresAt: expiresAt,
	}
	err = s.sessionRepository.CreateSession(session)
	if err != nil {
		return nil, err
	}
	return &session, nil
}

// Login function to authenticate the user and generate a session
func (s *SessionHandler) Login(w http.ResponseWriter, r *http.Request) {
	var params loginParams
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		log.Printf("Error decoding request body: %v", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	user, err := s.userRepository.GetUserByEmail(params.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("User not found for email: %s", params.Email)
			if err := invalidCredentials(w); err != nil {
				http.Error(w, "Invalid email", http.StatusUnauthorized)
				return
			}
			return
		}
		log.Printf("Error fetching user by email: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if !models.IsValidPassword(user.EncryptedPassword, params.Password) {
		if err := invalidCredentials(w); err != nil {
			http.Error(w, "Invalid password", http.StatusUnauthorized)
			return
		}
		return
	}

	session, err := s.GenerateSession(user.ID)
	if err != nil {
		log.Printf("Error generating session: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Sets the session token as a response header or in a cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    session.Token,
		Expires:  session.ExpiresAt,
		HttpOnly: true, // HttpOnly cookie important for security
		SameSite: http.SameSiteStrictMode,
	})

	writeJSONResponse(w, http.StatusOK, map[string]string{"message": "Login successful"})
}

// Logout function to delete the user's session
func (s *SessionHandler) Logout(w http.ResponseWriter, r *http.Request) {
	// Gets the session token from the request
	token, err := r.Cookie("session_token")
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Deletes the session from the repository
	err = s.sessionRepository.DeleteSession(token.Value)
	if err != nil {
		http.Error(w, "Failed to log out", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})

	writeJSONResponse(w, http.StatusOK, map[string]string{"message": "Logged out successfully"})
}

// ValidateSession middleware to check if the user's session is valid
func (s *SessionHandler) ValidateSession(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the session token from the request
		token, err := r.Cookie("session_token")
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Retrieves the session from the repository
		session, err := s.sessionRepository.GetSessionByToken(token.Value)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Checks if the session has expired
		if time.Now().After(session.ExpiresAt) {
			// Expire the session by deleting it from the repository
			err := s.sessionRepository.DeleteSession(token.Value)
			if err != nil {
				http.Error(w, "Failed to log out", http.StatusInternalServerError)
				return
			}

			http.SetCookie(w, &http.Cookie{
				Name:     "session_token",
				Value:    "",
				Expires:  time.Now().Add(-time.Hour),
				HttpOnly: true,
				SameSite: http.SameSiteStrictMode,
			})

			http.Error(w, "Session expired", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), ContextKeyUserID, session.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
