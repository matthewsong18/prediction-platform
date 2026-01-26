package bets

import (
	"errors"
	"testing"

	"betting-discord-bot/internal/polls"
)

func TestCreateBet(t *testing.T) {
	t.Parallel()
	pollMemoryRepo := polls.NewMemoryRepository()
	pollService := polls.NewService(pollMemoryRepo)
	betRepo := NewMemoryRepository()
	betService := NewService(pollService, betRepo)

	poll, err := pollService.CreatePoll("Test Poll", []string{"Option 1", "Option 2"})
	if err != nil {
		t.Fatal("Failed to create poll:", err)
	}

	pollId := poll.GetID()
	userId := "12345"
	selectedOptionIndex := 0
	bet, err1 := betService.CreateBet(pollId, userId, selectedOptionIndex)

	if err1 != nil {
		t.Fatal("CreateBet returned an unexpected error:", err1)
	}

	if bet.GetBetKey().PollID != pollId {
		t.Errorf("Expected bet to be associated with poll %s, but got %s", pollId, bet.GetBetKey().PollID)
	}

	if bet.GetBetKey().UserID != userId {
		t.Errorf("Expected bet to be associated with user %s, but got %s", userId, bet.GetBetKey().UserID)
	}

	if bet.GetSelectedOptionIndex() != selectedOptionIndex {
		t.Errorf("Expected bet to select option %d, but got %d", selectedOptionIndex, bet.GetSelectedOptionIndex())
	}
}

func TestInvalidOption(t *testing.T) {
	t.Parallel()
	pollMemoryRepo := polls.NewMemoryRepository()
	pollService := polls.NewService(pollMemoryRepo)
	betService := NewService(pollService, nil)
	pollId := "12345"
	userId := "12345"
	selectedOptionIndex := -1 // Invalid index
	_, err := betService.CreateBet(pollId, userId, selectedOptionIndex)

	if err == nil {
		t.Fatal("Expected CreateBet to return an error for invalid option index, but got nil")
	}

	if err.Error() != "invalid option index" {
		t.Errorf("Expected error message 'invalid option index', but got '%s'", err.Error())
	}
}

func TestPreventingMultipleBetsPerPoll(t *testing.T) {
	t.Parallel()
	pollMemoryRepo := polls.NewMemoryRepository()
	pollService := polls.NewService(pollMemoryRepo)
	betRepo := NewMemoryRepository()
	betService := NewService(pollService, betRepo)

	poll, _ := pollService.CreatePoll("Test Poll", []string{"Option 1", "Option 2"})

	// Create the first bet for the poll
	pollId := poll.GetID()
	userId := "12345"
	selectedOptionIndex := 0

	_, _ = betService.CreateBet(pollId, userId, selectedOptionIndex)

	// Attempt to create a second bet for the same poll
	_, err := betService.CreateBet(pollId, userId, selectedOptionIndex)

	if err == nil {
		t.Fatal("Expected an error when creating a second bet for the same poll, but got nil")
	}

	if !errors.Is(err, ErrUserAlreadyBet) {
		t.Errorf("Expected error message '%s', but got '%s'", ErrUserAlreadyBet.Error(), err.Error())
	}
	if err.Error() != "user already placed a bet on this poll" {
	}
}

func TestCannotBetOnClosedPoll(t *testing.T) {
	t.Parallel()
	pollMemoryRepo := polls.NewMemoryRepository()
	pollService := polls.NewService(pollMemoryRepo)
	betService := NewService(pollService, nil)

	poll, err := pollService.CreatePoll("Test Poll", []string{"Option 1", "Option 2"})
	if err != nil {
		t.Fatal("Failed to create poll:", err)
	}

	if err := pollService.ClosePoll(poll.GetID()); err != nil {
		t.Fatal("Failed to close poll:", err)
	}

	// Attempt to create a bet on a closed poll
	_, err = betService.CreateBet(poll.GetID(), "12345", 0)
	if err == nil {
		t.Fatal("Expected an error when betting on a closed poll, but got nil")
	}

	// Check if the error message is as expected
	if !errors.Is(err, ErrPollIsClosed) {
		t.Errorf("Expected error message '%s', but got '%s'", ErrPollIsClosed.Error(), err.Error())
	}
}

func TestGetBetOutcome(t *testing.T) {
	t.Parallel()
	// Check if the bet outcome is correctly retrieved

	pollMemoryRepo := polls.NewMemoryRepository()
	pollService := polls.NewService(pollMemoryRepo)
	betRepo := NewMemoryRepository()
	betService := NewService(pollService, betRepo)
	poll, err := pollService.CreatePoll("Test Poll", []string{"Option 1", "Option 2"})
	if err != nil {
		t.Fatal("Failed to create poll:", err)
	}

	pollId := poll.GetID()
	userId := "12345"
	selectedOptionIndex := 0
	bet, err1 := betService.CreateBet(pollId, userId, selectedOptionIndex)

	if err1 != nil {
		t.Fatal("CreateBet returned an unexpected error:", err1)
	}

	if bet.GetBetStatus() != Pending {
		t.Fatalf("Expected bet status to be 'PENDING', but got '%s'", bet.GetBetStatus())
	}

	if err := pollService.ClosePoll(poll.GetID()); err != nil {
		t.Fatal("ClosePoll returned an unexpected error:", err)
	}

	err = pollService.SelectOutcome(poll.GetID(), polls.OutcomeStatus(selectedOptionIndex))
	if err != nil {
		t.Fatal("SelectOutcome returned an unexpected error:", err)
	}

	if err := betService.UpdateBetsByPollId(poll.GetID()); err != nil {
		t.Fatal("UpdateBetsByPollId returned an unexpected error:", err)
	}
	bet, err2 := betService.GetBet(pollId, userId)

	if err2 != nil {
		t.Fatal("GetBet returned an unexpected error:", err2)
	}

	if bet.GetBetStatus() != Won {
		t.Errorf("Expected bet status to be 'WON', but got '%s'", bet.GetBetStatus())
	}

	if bet.GetSelectedOptionIndex() != selectedOptionIndex {
		t.Errorf("Expected bet to select option %d, but got %d", selectedOptionIndex, bet.GetSelectedOptionIndex())
	}
}

func TestGettingUserBets(t *testing.T) {
	t.Parallel()
	pollMemoryRepo := polls.NewMemoryRepository()
	pollService := polls.NewService(pollMemoryRepo)
	betRepo := NewMemoryRepository()
	betService := NewService(pollService, betRepo)

	poll, createPollErr := pollService.CreatePoll("Test Poll", []string{"Option 1", "Option 2"})
	if createPollErr != nil {
		t.Fatal("Failed to create poll:", createPollErr)
	}

	userID := "12345"
	bet, createBetErr := betService.CreateBet(poll.GetID(), userID, 0)
	if createBetErr != nil {
		t.Fatal("Failed to create bet:", createBetErr)
	}

	bets, getBetsErr := betService.GetBetsFromUser(userID)
	if getBetsErr != nil {
		t.Fatal("Failed to get bets from user:", getBetsErr)
	}

	if len(bets) != 1 {
		t.Errorf("Expected 1 bet for user %s, but got %d", userID, len(bets))
	}

	if bets[0] != bet || bets[0].GetBetKey().PollID != poll.GetID() || bets[0].GetBetKey().UserID != userID {
		t.Errorf("Expected bet %v, but got %v", bet, bets[0])
	}
}
