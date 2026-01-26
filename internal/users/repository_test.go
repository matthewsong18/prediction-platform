package users

import (
	"betting-discord-bot/internal/cryptography"
	"betting-discord-bot/internal/storage"
	"os"
	"strings"
	"testing"
)

func setupInMemory(t *testing.T) (UserRepository, func()) {
	t.Helper()

	repo := NewMemoryRepository()
	teardown := func() {}
	return repo, teardown
}

func setupLibSql(t *testing.T) (UserRepository, func()) {
	t.Helper()

	sanitizedTestName := strings.ReplaceAll(t.Name(), "/", "_")
	dbPath := sanitizedTestName + ".db"
	_ = os.Remove(dbPath)

	db, err := storage.InitializeDatabase(dbPath, "")
	if err != nil {
		t.Fatalf("Failed to initialize test database: %v", err)
	}

	cryptoService := setupCryptoService(t)

	repo := NewLibSQLRepository(db, cryptoService)
	teardown := func() {
		if err := db.Close(); err != nil {
			t.Fatal("failed to close database")
		}
		if err := os.Remove(dbPath); err != nil {
			t.Fatal("failed to remove database file")
		}
	}

	return repo, teardown
}

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

func TestUserRepositoryImplementations(t *testing.T) {
	t.Parallel()
	implementations := []struct {
		name  string
		setup func(t *testing.T) (UserRepository, func())
	}{
		{name: "InMemoryRepository", setup: setupInMemory},
		{name: "LibSQLRepository", setup: setupLibSql},
	}

	tests := map[string]func(t *testing.T, repo UserRepository){
		"it should save then get a user":           testSaveAndGet,
		"it should get a user by their discord ID": testGetByDiscordID,
		"it should delete a user":                  testDelete,
	}

	for _, implementation := range implementations {
		t.Run(implementation.name, func(t *testing.T) {
			t.Parallel()

			for name, test := range tests {
				t.Run(name, func(t *testing.T) {
					t.Parallel()

					repo, teardown := implementation.setup(t)
					t.Cleanup(teardown)

					test(t, repo)
				})
			}
		})
	}
}

// testSaveAndGet tests saving a user and retrieving it by ID.
func testSaveAndGet(t *testing.T, repo UserRepository) {
	user := &user{
		ID:          "test-id",
		DiscordID:   "test-discord-id",
		Username:    "test-username",
		DisplayName: "test-display-name",
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

	if savedUser.Username != user.Username {
		t.Errorf("Expected Username %s, got %s", user.Username, savedUser.Username)
	}

	if savedUser.DisplayName != user.DisplayName {
		t.Errorf("Expected DisplayName %s, got %s", user.DisplayName, savedUser.DisplayName)
	}
}

func testGetByDiscordID(t *testing.T, repo UserRepository) {
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

// testDelete tests deleting a user by DiscordID.
func testDelete(t *testing.T, repo UserRepository) {
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
