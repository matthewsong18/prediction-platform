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

func (repo *libsqlRepository) Save(user *user) error {
	// Encrypt sensitive data before saving
	encryptedDiscordID, err := repo.cryptoService.Encrypt(user.DiscordID)
	if err != nil {
		return fmt.Errorf("failed to encrypt discord_id: %w", err)
	}

	encryptedUsername, err := repo.cryptoService.Encrypt(user.Username)
	if err != nil {
		return fmt.Errorf("failed to encrypt username: %w", err)
	}

	encryptedDisplayName, err := repo.cryptoService.Encrypt(user.DisplayName)
	if err != nil {
		return fmt.Errorf("failed to encrypt display_name: %w", err)
	}

	// Generate blind index for searching
	discordIDHash := repo.cryptoService.GenerateBlindIndex(user.DiscordID)

	query := `INSERT INTO users (id, discord_id, discord_id_hash, username, display_name)
              VALUES (?, ?, ?, ?, ?)
              ON CONFLICT(id) DO UPDATE SET
              discord_id = excluded.discord_id,
              discord_id_hash = excluded.discord_id_hash,
              username = excluded.username,
              display_name = excluded.display_name`

	_, err = repo.db.Exec(query, user.ID, encryptedDiscordID, discordIDHash, encryptedUsername, encryptedDisplayName)
	if err != nil {
		return fmt.Errorf("error saving user: %w", err)
	}

	return nil
}

func (repo *libsqlRepository) GetByID(id string) (*user, error) {
	query := `SELECT id, discord_id, username, display_name FROM users WHERE id = ?`
	row := repo.db.QueryRow(query, id)

	var u user
	var encDiscordID, encUsername, encDisplayName string
	err := row.Scan(&u.ID, &encDiscordID, &encUsername, &encDisplayName)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("error retrieving user: %w", err)
	}

	// Decrypt sensitive data
	u.DiscordID, err = repo.cryptoService.Decrypt(encDiscordID)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt discord_id: %w", err)
	}

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

func (repo *libsqlRepository) GetByDiscordID(discordID string) (*user, error) {
	// Search by blind index (hash)
	discordIDHash := repo.cryptoService.GenerateBlindIndex(discordID)

	query := `SELECT id, discord_id, username, display_name FROM users WHERE discord_id_hash = ?`
	row := repo.db.QueryRow(query, discordIDHash)

	var retrievedUser user
	var encryptedDiscordID, encryptedUsername, encryptedDisplayName string
	err := row.Scan(&retrievedUser.ID, &encryptedDiscordID, &encryptedUsername, &encryptedDisplayName)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("error retrieving user: %w", err)
	}

	// Decrypt sensitive data
	retrievedUser.DiscordID, err = repo.cryptoService.Decrypt(encryptedDiscordID)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt discord_id: %w", err)
	}

	retrievedUser.Username, err = repo.cryptoService.Decrypt(encryptedUsername)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt username: %w", err)
	}

	retrievedUser.DisplayName, err = repo.cryptoService.Decrypt(encryptedDisplayName)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt display_name: %w", err)
	}

	if retrievedUser.DiscordID != discordID {
		return nil, ErrUserNotFound
	}

	return &retrievedUser, nil
}

func (repo *libsqlRepository) Delete(discordID string) error {
	// Search by blind index (hash)
	discordIDHash := repo.cryptoService.GenerateBlindIndex(discordID)

	query := "DELETE FROM users WHERE discord_id_hash = ?"
	result, err := repo.db.Exec(query, discordIDHash)
	if err != nil {
		return fmt.Errorf("error deleting user with discord_id %s: %w", discordID, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking rows affected for discord_id %s: %w", discordID, err)
	}

	if rowsAffected == 0 {
		return ErrUserNotFound
	}

	return nil
}

var _ UserRepository = (*libsqlRepository)(nil)
