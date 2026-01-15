// internal/bets/repository_test.go

package bets

import (
	"os"
	"strconv"
	"strings"
	"testing"

	"betting-discord-bot/internal/storage"
)

// setupLibSQL is a helper function specifically for the LibSQL implementation.
func setupLibSQL(t *testing.T) (BetRepository, func()) {
	t.Helper()

	// Sanitize the test name to create a clean, unique filename for each test run.
	sanitizedTestName := strings.ReplaceAll(t.Name(), "/", "_")
	dbPath := sanitizedTestName + ".db"

	// Proactively remove any old database file from a previous failed run.
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

// setupInMemory is a helper function for the in-memory implementation.
func setupInMemory(t *testing.T) (BetRepository, func()) {
	t.Helper()
	repo := NewMemoryRepository()
	teardown := func() {
		// No cleanup needed for the in-memory version
	}
	return repo, teardown
}

// TestBetRepositoryImplementations is the main entry point for testing all BetRepository implementations.
func TestBetRepositoryImplementations(t *testing.T) {
	t.Parallel()
	// This table defines all the implementations we want to test.
	implementations := []struct {
		name  string
		setup func(t *testing.T) (BetRepository, func())
	}{
		{name: "InMemoryRepository", setup: setupInMemory},
		{name: "LibSQLRepository", setup: setupLibSQL},
	}

	// This table defines all the behavioral tests we want to run against each implementation.
	testCases := []struct {
		name string
		run  func(t *testing.T, repo BetRepository)
	}{
		{"it should save and get a bet", testSaveAndGet},
		{"it should get all bets from a user", testGetAllBetsFromUser},
		{"it should get all bets from a poll", testGetAllBetsFromPoll},
		{"it should update a bet", testUpdateBet},
		{"it should get all user stats", testGetAllUserStats},
	}

	// Loop through each implementation and run each test against it. Did this
	// because I needed the setup/teardown to run for each individual test instead of
	// once per implementation.
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

func testSaveAndGet(t *testing.T, repo BetRepository) {
	// ARRANGE
	bet := &bet{
		PollID:              "poll123",
		UserID:              "user456",
		SelectedOptionIndex: 0,
		BetStatus:           Pending,
	}

	// ACT & ASSERT (Save)
	if err := repo.Save(bet); err != nil {
		t.Fatalf("Failed to save bet: %v", err)
	}

	// ACT & ASSERT (Get)
	retrievedBet, err := repo.GetByPollIdAndUserId(bet.PollID, bet.UserID)
	if err != nil {
		t.Fatalf("Failed to get bet by PollID and UserID: %v", err)
	}
	if retrievedBet == nil {
		t.Fatal("Retrieved bet is nil, expected a valid bet")
	}
	if retrievedBet.PollID != bet.PollID || retrievedBet.UserID != bet.UserID {
		t.Errorf("Retrieved bet does not match original: got %+v, want %+v", retrievedBet, bet)
	}
}

func testGetAllBetsFromUser(t *testing.T, repo BetRepository) {
	// ARRANGE
	userID := "user789"
	bets := []bet{
		{PollID: "poll1", UserID: userID, SelectedOptionIndex: 0, BetStatus: Pending},
		{PollID: "poll2", UserID: userID, SelectedOptionIndex: 1, BetStatus: Pending},
	}
	for _, bet := range bets {
		if err := repo.Save(&bet); err != nil {
			t.Fatalf("Failed to save bet: %v", err)
		}
	}

	// ACT
	retrievedBets, err := repo.GetBetsFromUser(userID)
	if err != nil {
		t.Fatalf("Failed to get bets from user: %v", err)
	}

	// ASSERT
	if len(retrievedBets) != len(bets) {
		t.Fatalf("Expected %d bets for user %s, got %d", len(bets), userID, len(retrievedBets))
	}
}

func testGetAllBetsFromPoll(t *testing.T, repo BetRepository) {
	// ARRANGE
	pollID := "poll456"
	bets := []bet{
		{PollID: pollID, UserID: "user1", SelectedOptionIndex: 0, BetStatus: Pending},
		{PollID: pollID, UserID: "user2", SelectedOptionIndex: 1, BetStatus: Pending},
	}
	for _, bet := range bets {
		if err := repo.Save(&bet); err != nil {
			t.Fatalf("Failed to save bet: %v", err)
		}
	}

	// ACT
	retrievedBets, err := repo.GetBetsByPollId(pollID)
	if err != nil {
		t.Fatalf("Failed to get bets by PollID: %v", err)
	}

	// ASSERT
	if len(retrievedBets) != len(bets) {
		t.Fatalf("Expected %d bets for poll %s, got %d", len(bets), pollID, len(retrievedBets))
	}
}

func testUpdateBet(t *testing.T, repo BetRepository) {
	// ARRANGE
	bet := &bet{
		PollID:    "poll123",
		UserID:    "user456",
		BetStatus: Pending,
	}
	if err := repo.Save(bet); err != nil {
		t.Fatalf("Failed to save initial bet: %v", err)
	}

	// ACT
	bet.BetStatus = Won
	if err := repo.UpdateBet(bet); err != nil {
		t.Fatalf("Failed to update bet: %v", err)
	}

	// ASSERT
	retrievedBet, err := repo.GetByPollIdAndUserId(bet.PollID, bet.UserID)
	if err != nil {
		t.Fatalf("Failed to get updated bet: %v", err)
	}
	if retrievedBet.BetStatus != Won {
		t.Errorf("Expected bet status %v, but got %v", Won, retrievedBet.BetStatus)
	}
}

func testGetAllUserStats(t *testing.T, repo BetRepository) {
	// ARRANGE
	var bets []bet
	user1ID := "user1"
	index := 0
	for range 10 {
		newBet := bet{
			PollID:              "poll" + strconv.Itoa(index),
			UserID:              user1ID,
			SelectedOptionIndex: 0,
			BetStatus:           Won,
		}
		bets = append(bets, newBet)
		index++
	}
	for range 5 {
		newBet := bet{
			PollID:              "poll" + strconv.Itoa(index),
			UserID:              user1ID,
			SelectedOptionIndex: 0,
			BetStatus:           Lost,
		}
		bets = append(bets, newBet)
		index++
	}

	if len(bets) != 15 {
		t.Fatalf("Expected 15 bets after setup but have %d instead", len(bets))
	}

	for _, bet := range bets {
		err := repo.Save(&bet)
		if err != nil {
			t.Fatalf("Failed to save bets during setup: %v", err)
		}
	}

	user1Bets, err := repo.GetBetsFromUser(user1ID)
	if err != nil {
		t.Fatalf("Failed to obtain bets after setup: %v", err)
	}
	if len(user1Bets) != 15 {
		t.Fatalf("Expected 15 bets after setup but got %d instead", len(user1Bets))
	}

	// ACT
	actualUserStats, err := repo.GetAllUserStats()
	if err != nil {
		t.Fatalf("Failed to get all user stats: %v", err)
	}

	// ASSERT
	if len(actualUserStats) != 1 {
		t.Fatalf("Expected 1 user stat but got %d", len(actualUserStats))
	}

	actualUserStat := actualUserStats[0]
	if actualUserStat.UserID != user1ID {
		t.Errorf("Expected %s but got %s", user1ID, actualUserStat.UserID)
	}
	if actualUserStat.Wins != 10 {
		t.Errorf("Expected %d but got %d", 10, actualUserStat.Wins)
	}
	if actualUserStat.Losses != 5 {
		t.Errorf("Expected %d but got %d", 5, actualUserStat.Losses)
	}
	if actualUserStat.Total != 15 {
		t.Errorf("Expected %d but got %d", 15, actualUserStat.Total)
	}
	if actualUserStat.WinLossRatio != float64(2) {
		t.Errorf("Expected %f but got %f", float64(2), actualUserStat.WinLossRatio)
	}
}
