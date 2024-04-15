package main

import (
	"database/sql"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/osag1e/user-account-auth-system/db/health"
	"github.com/osag1e/user-account-auth-system/internal/query"
	"github.com/osag1e/user-account-auth-system/service/handlers"
)

func initializeRouter(dbConn *sql.DB) http.Handler {
	router := chi.NewRouter()

	router.Get("/health", health.PostgreSQLHealthCheckHandler(dbConn))

	// Initializing the repositories and handlers
	userRepository := query.NewPostgresUserRepository(dbConn)
	sessionRepository := query.NewPostgresSesssionRepository(dbConn)
	session := handlers.NewSessionHandler(dbConn, userRepository, sessionRepository)
	userHandler := handlers.NewUserHandler(userRepository)

	// Defining Routes and Handlers
	// Create user and get users does not need session validation
	router.Post("/user", userHandler.HandleCreateUser)
	router.Get("/users", userHandler.HandleFetchUsers)
	router.Get("/user/{userID}", userHandler.HandleFetchUserByID)

	// Login is used to generate session
	router.Post("/login", session.Login)

	// Applying the ValidateSession middleware to routes that need session validation
	router.With(session.ValidateSession).Post("/logout", session.Logout)
	router.With(session.ValidateSession).Put("/user/{userID}", userHandler.HandleUserUpdate)
	router.With(session.ValidateSession).Delete("/user/{userID}", userHandler.HandleDeleteUser)

	return router
}
