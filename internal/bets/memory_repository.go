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
	userWinLossStats := make(map[string][2]int)
	for _, bet := range repo.betList {
		win := 0
		loss := 0
		switch bet.BetStatus {
		case Won:
			win = 1
		case Lost:
			loss = 1
		}
		user := bet.UserID
		if _, exists := userWinLossStats[user]; !exists {
			userWinLossStats[user] = [2]int{win, loss}
		} else {
			existingStat := userWinLossStats[user]
			userWinLossStats[user] = [2]int{existingStat[0] + win, existingStat[1] + loss}
		}
	}

	allUserStats := make([]*UserStats, 0, 1)
	for user, winLoss := range userWinLossStats {
		newUserStat := &UserStats{
			UserID:       user,
			Wins:         winLoss[0],
			Losses:       winLoss[1],
			Total:        winLoss[0] + winLoss[1],
			WinLossRatio: float64(winLoss[0]) / float64(winLoss[1]),
		}
		allUserStats = append(allUserStats, newUserStat)
	}

	return allUserStats, nil
}

var _ BetRepository = (*memoryRepository)(nil)
