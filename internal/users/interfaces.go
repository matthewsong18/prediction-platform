package users

import "errors"

type UserService interface {
	// CreateUser creates a new internal user and links it to the given provider identity.
	CreateUser(identity Identity) (User, error)
	// GetUserByExternalID finds a user by their provider identity
	GetUserByExternalID(identity Identity) (User, error)
	// DeleteUser deletes the user and all associated identities.
	// For now, we still trigger this via a specific provider identity.
	DeleteUser(identity Identity) error
	GetWinLoss(userID string) (*WinLoss, error)
}

type UserRepository interface {
	Save(user *user, identity *Identity) error
	// AddIdentity links an external identity to an existing user.
	AddIdentity(userID string, identity *Identity) error
	GetByID(id string) (*user, error)
	GetByExternalID(identity *Identity) (*user, error)
	// Delete deletes the user and their identities.
	// The implementation should handle the cascade or multi-table deletion.
	Delete(userID string) error

	// Private testing methods

	getUserCount() (int, error)
}

var ErrUserNotFound = errors.New("user not found")
