package users

import (
	"betting-discord-bot/internal/cryptography"
	"database/sql"
	"errors"
	"fmt"
)

type libsqlRepository struct {
	db            *sql.DB
	cryptoService cryptography.CryptoService
}

func NewLibSQLRepository(db *sql.DB, cryptoService cryptography.CryptoService) UserRepository {
	return &libsqlRepository{db, cryptoService}
}

func (repo *libsqlRepository) Save(user *user, provider, externalID string) error {
	// Encrypt sensitive data before saving
	encryptedUsername, err := repo.cryptoService.Encrypt(user.Username)
	if err != nil {
		return fmt.Errorf("failed to encrypt username: %w", err)
	}

	encryptedDisplayName, err := repo.cryptoService.Encrypt(user.DisplayName)
	if err != nil {
		return fmt.Errorf("failed to encrypt display_name: %w", err)
	}

	userQuery := `INSERT INTO users (id, username, display_name)
              VALUES (?, ?, ?)
              ON CONFLICT(id) DO UPDATE SET
              username = excluded.username,
              display_name = excluded.display_name`

	_, err = repo.db.Exec(userQuery, user.ID, encryptedUsername, encryptedDisplayName)
	if err != nil {
		return fmt.Errorf("error saving user: %w", err)
	}

	encryptedExternalID, err := repo.cryptoService.Encrypt(externalID)
	if err != nil {
		return fmt.Errorf("failed to encrypt external_id: %w", err)
	}

	externalIDHash := repo.cryptoService.GenerateBlindIndex(externalID)

	identityQuery := `INSERT INTO user_identities (provider, external_id, external_id_hash, user_id)
              VALUES (?, ?, ?, ?)
              ON CONFLICT(provider, external_id) DO UPDATE SET
              external_id_hash = excluded.external_id_hash,
              user_id = excluded.user_id`

	_, err = repo.db.Exec(identityQuery, provider, encryptedExternalID, externalIDHash, user.ID)
	if err != nil {
		return fmt.Errorf("error saving identity: %w", err)
	}

	return nil
}

func (repo *libsqlRepository) AddIdentity(userID, provider, externalID string) error {
	encryptedExternalID, err := repo.cryptoService.Encrypt(externalID)
	if err != nil {
		return fmt.Errorf("failed to encrypt external_id: %w", err)
	}

	externalIDHash := repo.cryptoService.GenerateBlindIndex(externalID)

	query := `INSERT INTO user_identities (provider, external_id, external_id_hash, user_id)
              VALUES (?, ?, ?, ?)
              ON CONFLICT(provider, external_id) DO UPDATE SET
              external_id_hash = excluded.external_id_hash,
              user_id = excluded.user_id`

	_, err = repo.db.Exec(query, provider, encryptedExternalID, externalIDHash, userID)
	if err != nil {
		return fmt.Errorf("error saving identity: %w", err)
	}

	return nil
}

func (repo *libsqlRepository) GetByID(id string) (*user, error) {
	query := `SELECT id, username, display_name FROM users WHERE id = ?`
	row := repo.db.QueryRow(query, id)

	var u user
	var encUsername, encDisplayName string
	err := row.Scan(&u.ID, &encUsername, &encDisplayName)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("error retrieving user: %w", err)
	}

	// Decrypt sensitive data
	u.Username, err = repo.cryptoService.Decrypt(encUsername)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt username: %w", err)
	}

	u.DisplayName, err = repo.cryptoService.Decrypt(encDisplayName)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt display_name: %w", err)
	}

	return &u, nil
}

func (repo *libsqlRepository) GetByExternalID(provider, externalID string) (*user, error) {
	// Search by blind index (hash)
	externalIDHash := repo.cryptoService.GenerateBlindIndex(externalID)

	query := `SELECT u.id, u.username, u.display_name, ui.external_id
              FROM users u
              JOIN user_identities ui ON u.id = ui.user_id
              WHERE ui.provider = ? AND ui.external_id_hash = ?`

	row := repo.db.QueryRow(query, provider, externalIDHash)

	var retrievedUser user
	var encryptedExternalID, encryptedUsername, encryptedDisplayName string
	err := row.Scan(&retrievedUser.ID, &encryptedUsername, &encryptedDisplayName, &encryptedExternalID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("error retrieving user by identity: %w", err)
	}

	// Verify the external ID matches (to prevent hash collisions)
	decryptedExternalID, err := repo.cryptoService.Decrypt(encryptedExternalID)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt external_id: %w", err)
	}
	if decryptedExternalID != externalID {
		return nil, ErrUserNotFound
	}

	// Decrypt sensitive data
	retrievedUser.Username, err = repo.cryptoService.Decrypt(encryptedUsername)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt username: %w", err)
	}

	retrievedUser.DisplayName, err = repo.cryptoService.Decrypt(encryptedDisplayName)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt display_name: %w", err)
	}

	return &retrievedUser, nil
}

func (repo *libsqlRepository) Delete(userID string) error {
	query := "DELETE FROM users WHERE id = ?"
	result, err := repo.db.Exec(query, userID)
	if err != nil {
		return fmt.Errorf("error deleting user %s: %w", userID, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking rows affected for user %s: %w", userID, err)
	}

	if rowsAffected == 0 {
		return ErrUserNotFound
	}

	return nil
}

var _ UserRepository = (*libsqlRepository)(nil)
