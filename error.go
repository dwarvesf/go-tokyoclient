package tokyoclient

import "errors"

var (
	// ErrInvalidID is the error when the ID is invalid
	ErrInvalidID = errors.New("invalid ID")

	// ErrInvalidTeamNames is the error when the team names are invalid
	ErrInvalidTeamNames = errors.New("invalid team names")

	// ErrInvalidState is the error when the state is invalid
	ErrInvalidState = errors.New("invalid state")

	// ErrConnNotInitialized is the error when the connection is not initialized
	ErrConnNotInitialized = errors.New("connection not initialized")
)
