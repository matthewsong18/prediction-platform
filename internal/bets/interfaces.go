package bets

import "errors"

type BetService interface {
	CreateBet(pollID string, userID string, selectedOptionIndex int) (Bet, error)
	GetBet(pollID string, userID string) (Bet, error)
	UpdateBetsByPollId(pollID string) error
	GetBetsFromUser(userID string) ([]Bet, error)
}

type BetRepository interface {
	Save(bet *bet) error
	GetByPollIdAndUserId(pollID string, userID string) (*bet, error)
	GetBetsFromUser(userID string) ([]*bet, error)
	GetBetsByPollId(pollID string) ([]*bet, error)
	UpdateBet(bet *bet) error
	GetAllUserStats() ([]*UserStats, error)
}

// Errors related to bets

var ErrBetNotFound = errors.New("bet not found")
var ErrUserAlreadyBet = errors.New("user already bet")
var ErrPollIsClosed = errors.New("poll is closed")
