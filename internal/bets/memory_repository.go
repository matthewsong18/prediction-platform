package bets

import "errors"

type memoryRepository struct {
	betList map[BetKey]*bet
}

func NewMemoryRepository() BetRepository {
	return &memoryRepository{
		betList: make(map[BetKey]*bet),
	}
}

func (repo memoryRepository) Save(bet *bet) error {
	key := BetKey{bet.PollID, bet.UserID}
	if _, exists := repo.betList[key]; exists {
		return errors.New("user already placed a bet on this poll")
	}

	repo.betList[key] = bet
	return nil
}

func (repo memoryRepository) GetByPollIdAndUserId(pollID string, userID string) (*bet, error) {
	key := BetKey{PollID: pollID, UserID: userID}
	if bet, exists := repo.betList[key]; exists {
		return bet, nil
	}
	return nil, errors.New("bet not found for the given poll and user")
}

func (repo memoryRepository) GetBetsFromUser(userID string) ([]*bet, error) {
	var bets []*bet
	for key, bet := range repo.betList {
		if key.UserID == userID {
			bets = append(bets, bet)
		}
	}
	return bets, nil
}

func (repo memoryRepository) GetBetsByPollId(pollID string) ([]*bet, error) {
	var bets []*bet
	for key, bet := range repo.betList {
		if key.PollID == pollID {
			bets = append(bets, bet)
		}
	}
	return bets, nil
}

func (repo memoryRepository) UpdateBet(bet *bet) error {
	key := BetKey{bet.PollID, bet.UserID}
	if _, exists := repo.betList[key]; !exists {
		return errors.New("bet not found for the given poll and user")
	}

	repo.betList[key] = bet
	return nil
}

func (repo *memoryRepository) GetAllUserStats() ([]*UserStats, error) {
	return nil, errors.New("unimplemented")
}

var _ BetRepository = (*memoryRepository)(nil)
