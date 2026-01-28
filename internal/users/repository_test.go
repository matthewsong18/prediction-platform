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
		"it should save then get a user":            testSaveAndGet,
		"it should get a user by their external ID": testGetByExternalID,
		"it should delete a user":                   testDelete,
		"it should save the user atomically":        testSaveUserIsAtomicTransaction,
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
		Username:    "test-username",
		DisplayName: "test-display-name",
	}

	if err := repo.Save(user, "test-provider", "test-external-id"); err != nil {
		t.Fatalf("Failed to save user: %v", err)
	}

	savedUser, err := repo.GetByID(user.ID)
	if err != nil {
		t.Fatalf("Failed to get user by ID: %v", err)
	}

	if savedUser.Username != user.Username {
		t.Errorf("Expected Username %s, got %s", user.Username, savedUser.Username)
	}

	if savedUser.DisplayName != user.DisplayName {
		t.Errorf("Expected DisplayName %s, got %s", user.DisplayName, savedUser.DisplayName)
	}
}

func testGetByExternalID(t *testing.T, repo UserRepository) {
	user := &user{
		ID: "test-id",
	}

	if err := repo.Save(user, "test-provider", "test-external-id"); err != nil {
		t.Fatalf("Failed to save user: %v", err)
	}

	savedUser, err := repo.GetByExternalID("test-provider", "test-external-id")
	if err != nil {
		t.Fatalf("Failed to get user by External ID: %v", err)
	}

	if savedUser.ID != user.ID {
		t.Errorf("Expected ID %s, got %s", user.ID, savedUser.ID)
	}
}

// testDelete tests deleting a user
func testDelete(t *testing.T, repo UserRepository) {
	user := &user{
		ID: "test-id",
	}

	// Save the user first
	if err := repo.Save(user, "test-provider", "test-external-id"); err != nil {
		t.Fatalf("Failed to save user: %v", err)
	}

	// Ensure user exists before deletion
	_, err := repo.GetByExternalID("test-provider", "test-external-id")
	if err != nil {
		t.Fatalf("Failed to get user by ExternalID before deletion: %v", err)
	}

	// Delete the user
	if err := repo.Delete(user.ID); err != nil {
		t.Fatalf("Failed to delete user: %v", err)
	}

	// Assert that the user no longer exists
	_, err = repo.GetByExternalID("test-provider", "test-external-id")
	if err == nil {
		t.Fatal("Expected error when getting deleted user, got none")
	}
}

func testSaveUserIsAtomicTransaction(t *testing.T, repo UserRepository) {
	// Creating a random user to occupy the identity
	blockingUser := &user{
		ID: "blocker",
	}
	if err := repo.Save(blockingUser, "test-provider", "blocking-id"); err != nil {
		t.Fatalf("Failed to save setup blocking user: %v", err)
	}

	startingUserCount, err := repo.getUserCount()
	if err != nil {
		t.Fatalf("Failed to count users: %v", err)
	}

	// Trying to save a different user with the same identity
	testUser := &user{
		ID: "new user",
	}
	if err := repo.Save(testUser, "test-provider", "blocking-id"); err == nil {
		t.Fatal("Didn't fail on already used identity")
	}

	endingUserCount, err := repo.getUserCount()
	if err != nil {
		t.Fatalf("Failed to count users: %v", err)
	}

	if startingUserCount != endingUserCount {
		t.Fatal("SaveUser is not atomic")
	}
}
