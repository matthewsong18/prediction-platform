package users

import (
	"betting-discord-bot/internal/cryptography"
	"database/sql"
	"os"
	"testing"

	"betting-discord-bot/internal/storage"
)

// setupCryptoService creates a cryptography service for testing.
func setupCryptoService(t *testing.T) cryptography.CryptoService {
	t.Helper()
	var key [32]byte // zero key is fine for tests
	service, err := cryptography.NewService(key)
	if err != nil {
		t.Fatalf("failed to create crypto service: %v", err)
	}
	return service
}

// Creates a temporary database for testing user libsqlRepository.
func setupTestDB(t *testing.T) (*sql.DB, func()) {
	t.Helper()

	dbPath := t.Name() + ".db"
	db, err := storage.InitializeDatabase(dbPath, "")
	if err != nil {
		t.Fatalf("Failed to initialize test database: %v", err)
	}

	teardown := func() {
		err := db.Close()
		if err != nil {
			t.Log("Failed to close test database:", err)
		}

		err = os.Remove(dbPath)
		if err != nil {
			t.Log("Failed to remove test database:", err)
		}
	}

	return db, teardown
}

// TestSaveAndGet tests saving a user and retrieving it by ID.
func TestSaveAndGet(t *testing.T) {
	db, teardown := setupTestDB(t)
	t.Cleanup(teardown)
	crypto := setupCryptoService(t)

	repo := NewLibSQLRepository(db, crypto)

	user := &user{
		ID:        "test-id",
		DiscordID: "test-discord-id",
	}

	if err := repo.Save(user); err != nil {
		t.Fatalf("Failed to save user: %v", err)
	}

	savedUser, err := repo.GetByID(user.ID)
	if err != nil {
		t.Fatalf("Failed to get user by ID: %v", err)
	}

	if savedUser.DiscordID != user.DiscordID {
		t.Errorf("Expected DiscordID %s, got %s", user.DiscordID, savedUser.DiscordID)
	}
}

func TestGetByDiscordID(t *testing.T) {
	db, teardown := setupTestDB(t)
	t.Cleanup(teardown)
	crypto := setupCryptoService(t)

	repo := NewLibSQLRepository(db, crypto)

	user := &user{
		ID:        "test-id",
		DiscordID: "test-discord-id",
	}

	if err := repo.Save(user); err != nil {
		t.Fatalf("Failed to save user: %v", err)
	}

	savedUser, err := repo.GetByDiscordID(user.DiscordID)
	if err != nil {
		t.Fatalf("Failed to get user by Discord ID: %v", err)
	}

	if savedUser.ID != user.ID {
		t.Errorf("Expected ID %s, got %s", user.ID, savedUser.ID)
	}
}

// TestDelete tests deleting a user by DiscordID.
func TestDelete(t *testing.T) {
	db, teardown := setupTestDB(t)
	t.Cleanup(teardown)
	crypto := setupCryptoService(t)

	repo := NewLibSQLRepository(db, crypto)

	user := &user{
		ID:        "test-id",
		DiscordID: "test-discord-id",
	}

	// Save the user first
	if err := repo.Save(user); err != nil {
		t.Fatalf("Failed to save user: %v", err)
	}

	// Ensure user exists before deletion
	_, err := repo.GetByDiscordID(user.DiscordID)
	if err != nil {
		t.Fatalf("Failed to get user by DiscordID before deletion: %v", err)
	}

	// Delete the user
	if err := repo.Delete(user.DiscordID); err != nil {
		t.Fatalf("Failed to delete user: %v", err)
	}

	// Assert that the user no longer exists
	_, err = repo.GetByDiscordID(user.DiscordID)
	if err == nil {
		t.Fatal("Expected error when getting deleted user, got none")
	}
}

// TestEncryptionTransparency verifies that data is stored as ciphertext but retrieved as plaintext.
func TestEncryptionTransparency(t *testing.T) {
	db, teardown := setupTestDB(t)
	t.Cleanup(teardown)
	crypto := setupCryptoService(t)

	repo := NewLibSQLRepository(db, crypto)

	originalDiscordID := "real-discord-id"
	originalUsername := "testuser"
	originalDisplayName := "Test User"
	u := &user{
		ID:          "user-1",
		DiscordID:   originalDiscordID,
		Username:    originalUsername,
		DisplayName: originalDisplayName,
	}

	// 1. Save user via repository
	if err := repo.Save(u); err != nil {
		t.Fatalf("Failed to save user: %v", err)
	}

	// 2. Query the database directly to verify ciphertext
	var dbDiscordID, dbUsername, dbDisplayName string
	err := db.QueryRow("SELECT discord_id, username, display_name FROM users WHERE id = ?", u.ID).Scan(&dbDiscordID, &dbUsername, &dbDisplayName)
	if err != nil {
		t.Fatalf("Direct DB query failed: %v", err)
	}

	if dbDiscordID == originalDiscordID {
		t.Errorf("CRITICAL: DiscordID stored in DB is plaintext! Expected ciphertext, got %q", dbDiscordID)
	}
	if dbUsername == originalUsername {
		t.Errorf("CRITICAL: Username stored in DB is plaintext! Expected ciphertext, got %q", dbUsername)
	}
	if dbDisplayName == originalDisplayName {
		t.Errorf("CRITICAL: DisplayName stored in DB is plaintext! Expected ciphertext, got %q", dbDisplayName)
	}

	// 3. Retrieve via repository to verify transparent decryption
	retrieved, err := repo.GetByID(u.ID)
	if err != nil {
		t.Fatalf("Repository retrieval failed: %v", err)
	}

	if retrieved.DiscordID != originalDiscordID {
		t.Errorf("Transparent decryption for DiscordID failed. Expected %q, got %q", originalDiscordID, retrieved.DiscordID)
	}
	if retrieved.Username != originalUsername {
		t.Errorf("Transparent decryption for Username failed. Expected %q, got %q", originalUsername, retrieved.Username)
	}
	if retrieved.DisplayName != originalDisplayName {
		t.Errorf("Transparent decryption for DisplayName failed. Expected %q, got %q", originalDisplayName, retrieved.DisplayName)
	}
}

