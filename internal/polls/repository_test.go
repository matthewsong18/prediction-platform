package polls

import (
	"os"
	"strings"
	"testing"

	"betting-discord-bot/internal/storage"

	"github.com/google/uuid"
)

func setupLibSQL(t *testing.T) (PollRepository, func()) {
	t.Helper()

	// Sanitize the test name to create a clean, unique filename for each test run.
	sanitizedTestName := strings.ReplaceAll(t.Name(), "/", "_")
	dbPath := sanitizedTestName + ".db"

	// Remove any old database file from a previous failed run.
	_ = os.Remove(dbPath)

	db, err := storage.InitializeDatabase(dbPath, "")
	if err != nil {
		t.Fatalf("Failed to initialize test database: %v", err)
	}

	repo := NewLibSQLRepository(db)

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

func setupInMemory(t *testing.T) (PollRepository, func()) {
	t.Helper()

	repo := NewMemoryRepository()
	teardown := func() {
		// No cleanup needed for the in-memory version
	}
	return repo, teardown
}

func TestPollRepositoryImplementations(t *testing.T) {
	t.Parallel()
	implementations := []struct {
		name  string
		setup func(t *testing.T) (PollRepository, func())
	}{
		{name: "InMemoryRepository", setup: setupInMemory},
		{name: "LibSQLRepository", setup: setupLibSQL},
	}

	testCases := []struct {
		name string
		run  func(t *testing.T, repo PollRepository)
	}{
		{"it should save and retrieve the poll", testSaveAndReceive},
		{"it should update the poll", testUpdate},
		{"it should delete the poll", testDelete},
		{"it should return all open polls", testGetAllOpenInRepo},
	}

	for _, impl := range implementations {
		t.Run(impl.name, func(t *testing.T) {
			t.Parallel()

			for _, tc := range testCases {
				t.Run(tc.name, func(t *testing.T) {
					t.Parallel()

					repo, cleanup := impl.setup(t)
					t.Cleanup(cleanup)

					// Run the actual test logic.
					tc.run(t, repo)
				})
			}
		})
	}
}

func testSaveAndReceive(t *testing.T, repo PollRepository) {
	// ARRANGE: Create a new poll to save
	pollToSave := &poll{
		ID:    uuid.New().String(),
		Title: "First Poll",
		Options: []string{
			"Option 1",
			"Option 2",
		},
		Outcome: Pending,
		Status:  Open,
	}

	// ACT: Save the poll
	err := repo.Save(pollToSave)
	if err != nil {
		t.Fatalf("Save() returned an unexpected error: %v", err)
	}

	// ACT: Get the poll back
	retrievedPoll, err := repo.GetById(pollToSave.ID)
	if err != nil {
		t.Fatalf("GetById() returned an unexpected error: %v", err)
	}

	// ASSERT
	if retrievedPoll.ID != pollToSave.ID {
		t.Errorf("Expected poll ID %s, but got %s", pollToSave.ID, retrievedPoll.ID)
	}
	if retrievedPoll.Title != pollToSave.Title {
		t.Errorf("Expected poll title %q, but got %q", pollToSave.Title, retrievedPoll.Title)
	}
	if retrievedPoll.Status != pollToSave.Status {
		t.Errorf("Expected poll status %v, but got %v", pollToSave.Status, retrievedPoll.Status)
	}
	if retrievedPoll.Outcome != pollToSave.Outcome {
		t.Errorf("Expected poll outcome %v, but got %v", pollToSave.Outcome, retrievedPoll.Outcome)
	}
}

func testUpdate(t *testing.T, repo PollRepository) {
	pollToUpdate := &poll{
		ID:     uuid.NewString(),
		Title:  "Poll to Update",
		Status: Open,
		Options: []string{
			"Option 1",
			"Option 2",
		},
		Outcome: Pending,
	}

	// ACT: Save the poll first
	if err := repo.Save(pollToUpdate); err != nil {
		t.Fatalf("Save() returned an unexpected error: %v", err)
	}

	// ACT: Update the poll
	pollToUpdate.Title = "Updated Poll Title"
	pollToUpdate.Status = Closed
	pollToUpdate.Outcome = Option2
	pollToUpdate.Options[0] = "Replaced"
	if err := repo.Update(pollToUpdate); err != nil {
		t.Fatalf("Update() returned an unexpected error: %v", err)
	}

	// ACT: Retrieve the updated poll
	retrievedPoll, err := repo.GetById(pollToUpdate.ID)
	if err != nil {
		t.Fatalf("GetById() returned an unexpected error: %v", err)
	}

	// ASSERT
	if retrievedPoll.Title != "Updated Poll Title" {
		t.Errorf("Expected updated poll title %q, but got %q", "Updated Poll Title", retrievedPoll.Title)
	}
	if retrievedPoll.Status != pollToUpdate.Status {
		t.Errorf("Expected poll status %v, but got %v", pollToUpdate.Status, retrievedPoll.Status)
	}
	if retrievedPoll.Outcome != pollToUpdate.Outcome {
		t.Errorf("Expected poll outcome %v, but got %v", pollToUpdate.Outcome, retrievedPoll.Outcome)
	}
	if retrievedPoll.Options[0] != "Replaced" {
		t.Errorf("Expected updated option %q, but got %q", "Replaced", retrievedPoll.Options[0])
	}
}

func testDelete(t *testing.T, repo PollRepository) {
	// ARRANGE: Create a new poll to delete
	pollToDelete := &poll{
		ID:     uuid.NewString(),
		Title:  "Poll to Delete",
		Status: Open,
		Options: []string{
			"Option 1",
			"Option 2",
		},
		Outcome: 0,
	}

	// ACT: Save the poll first
	if err := repo.Save(pollToDelete); err != nil {
		t.Fatalf("Save() returned an unexpected error: %v", err)
	}

	// ACT: Retrieve the poll to ensure it exists before deletion
	_, err := repo.GetById(pollToDelete.ID)
	if err != nil {
		t.Fatalf("GetById() returned an unexpected error: %v", err)
	}

	// ACT: Delete the poll
	if err := repo.Delete(pollToDelete.ID); err != nil {
		t.Fatalf("Delete() returned an unexpected error: %v", err)
	}

	// ASSERT: Try to retrieve the deleted poll
	_, err = repo.GetById(pollToDelete.ID)
	if err == nil {
		t.Fatal("Expected error when retrieving deleted poll, but got nil")
	}
}

func testGetAllOpenInRepo(t *testing.T, repo PollRepository) {
	// ARRANGE: Create open and closed polls
	if err := repo.Save(&poll{
		ID:    uuid.NewString(),
		Title: "open poll 1",
		Options: []string{
			"Option 1",
			"Option 2",
		},
		Status:  Open,
		Outcome: 0,
	}); err != nil {
		t.Fatalf("Save() returned an unexpected error: %v", err)
	}
	if err := repo.Save(&poll{
		ID:    uuid.NewString(),
		Title: "open poll 2",
		Options: []string{
			"Option 1",
			"Option 2",
		},
		Status:  Open,
		Outcome: 0,
	}); err != nil {
		t.Fatalf("Save() returned an unexpected error: %v", err)
	}
	if err := repo.Save(&poll{
		ID:    uuid.NewString(),
		Title: "closed poll 1",
		Options: []string{
			"Option 1",
			"Option 2",
		},
		Status:  Closed,
		Outcome: 0,
	}); err != nil {
		t.Fatalf("Save() returned an unexpected error: %v", err)
	}

	// ACT: Get Open Polls
	openPolls, err := repo.GetOpenPolls()
	if err != nil {
		t.Fatalf("GetOpenPolls() returned an unexpected error: %v", err)
	}

	// ASSERT
	if len(openPolls) != 2 {
		t.Errorf("Expected 2 open polls, but got %d", len(openPolls))
	}

	for _, poll := range openPolls {
		if poll.Status != Open {
			t.Errorf("Expected poll status %v, but got %v", Open, poll.Status)
		}
	}
}
