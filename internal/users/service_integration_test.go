package users

import (
	"testing"

	"betting-discord-bot/internal/bets"
	"github.com/google/uuid"
)

type mockBet struct {
	betKey bets.BetKey
	status bets.BetStatus
}

func (m mockBet) GetBetKey() bets.BetKey {
	return m.betKey
}

func (m mockBet) GetSelectedOptionIndex() int {
	panic("should not be called")
}

func (m mockBet) GetBetStatus() bets.BetStatus {
	return m.status
}

func createMockBet(status bets.BetStatus) bets.Bet {
	return mockBet{betKey: bets.BetKey{PollID: uuid.NewString(), UserID: "user"}, status: status}
}

type mockBetService struct {
	betsToReturn []bets.Bet
}

func (m *mockBetService) GetBetsFromUser(string) ([]bets.Bet, error) {
	return m.betsToReturn, nil
}

func (m *mockBetService) CreateBet(string, string, int) (bets.Bet, error) {
	return nil, nil
}
func (m *mockBetService) GetBet(string, string) (bets.Bet, error) { return nil, nil }
func (m *mockBetService) UpdateBetsByPollId(string) error         { return nil }

var _ bets.BetService = (*mockBetService)(nil)

func getTestBets(wins int, losses int, pending int) []bets.Bet {
	var betList []bets.Bet

	for i := 0; i < wins; i++ {
		betList = append(betList, createMockBet(bets.Won))
	}

	for i := 0; i < losses; i++ {
		betList = append(betList, createMockBet(bets.Lost))
	}

	for i := 0; i < pending; i++ {
		betList = append(betList, createMockBet(bets.Pending))
	}

	return betList
}

func TestGetUserWinLoss(t *testing.T) {
	t.Parallel()
	betsArgument := []struct {
		name    string
		betList []bets.Bet
		winLoss *WinLoss
	}{
		{
			"no bets",
			getTestBets(0, 0, 0),
			&WinLoss{Wins: 0, Losses: 0},
		},
		{
			"a winning bet",
			getTestBets(1, 0, 0),
			&WinLoss{Wins: 1, Losses: 0},
		},
		{
			"a losing bet",
			getTestBets(0, 1, 0),
			&WinLoss{Wins: 0, Losses: 1},
		},
		{
			"pending bets should not count",
			getTestBets(1, 1, 1),
			&WinLoss{1, 1},
		},
	}

	for _, tc := range betsArgument {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			testWinLoss(t, tc.betList, tc.winLoss)
		})
	}
}

func testWinLoss(t *testing.T, betsToReturn []bets.Bet, expectedWinLoss *WinLoss) {
	// ARRANGE
	mockBets := &mockBetService{
		betsToReturn: betsToReturn,
	}

	userRepo := NewMemoryRepository()
	userService := NewService(userRepo, mockBets)

	user, err := userService.CreateUser("discord", "test-discord-id")
	if err != nil {
		t.Fatalf("Setup failed: could not create user: %v", err)
	}

	// ACT
	actualWinLoss, err := userService.GetWinLoss(user.GetID())
	if err != nil {
		t.Fatalf("GetWinLoss returned an unexpected error: %v", err)
	}

	// ASSERT
	if actualWinLoss.Wins != expectedWinLoss.Wins {
		t.Errorf("Expected Wins to be %d, but got %d", expectedWinLoss.Wins, actualWinLoss.Wins)
	}

	if actualWinLoss.Losses != expectedWinLoss.Losses {
		t.Errorf("Expected Losses to be %d, but got %d", expectedWinLoss.Losses, actualWinLoss.Losses)
	}
}
