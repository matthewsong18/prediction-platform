package users

import "errors"

type UserService interface {
	// CreateUser creates a new internal user and links it to the given provider identity.
	CreateUser(provider, externalID string) (User, error)
	// GetUserByExternalID finds a user by their provider identity (e.g., "discord", "123").
	GetUserByExternalID(provider, externalID string) (User, error)
	// DeleteUser deletes the user and all associated identities.
	// For now, we still trigger this via a specific provider identity.
	DeleteUser(provider, externalID string) error
	GetWinLoss(userID string) (*WinLoss, error)
}

type UserRepository interface {
	Save(user *user, provider, externalID string) error
	// AddIdentity links an external identity to an existing user.
	AddIdentity(userID, provider, externalID string) error
	GetByID(id string) (*user, error)
	GetByExternalID(provider, externalID string) (*user, error)
	// Delete deletes the user and their identities.
	// The implementation should handle the cascade or multi-table deletion.
	Delete(userID string) error

	// Private testing methods

	getUserCount() (int, error)
}

var ErrUserNotFound = errors.New("user not found")
