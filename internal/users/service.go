package users

import (
	"fmt"

	"betting-discord-bot/internal/bets"

	"github.com/google/uuid"
)

type service struct {
	userRepo   UserRepository
	betService bets.BetService
}

func NewService(userRepo UserRepository, betService bets.BetService) UserService {
	return &service{
		userRepo:   userRepo,
		betService: betService,
	}
}

func (service service) CreateUser(provider, externalID string) (User, error) {
	user := &user{
		ID: uuid.NewString(),
	}

	err := service.userRepo.Save(user, provider, externalID)
	if err != nil {
		return nil, fmt.Errorf("could not save user: %w", err)
	}

	return user, nil
}

func (service service) GetUserByExternalID(provider, externalID string) (User, error) {
	user, userErr := service.userRepo.GetByExternalID(provider, externalID)
	if userErr != nil {
		return nil, userErr
	}

	return user, nil
}

func (service service) DeleteUser(provider, externalID string) error {
	user, err := service.userRepo.GetByExternalID(provider, externalID)
	if err != nil {
		return fmt.Errorf("could not find user to delete: %w", err)
	}

	err = service.userRepo.Delete(user.ID)
	if err != nil {
		return fmt.Errorf("could not delete user: %w", err)
	}

	return nil
}

func (service service) GetWinLoss(userID string) (*WinLoss, error) {
	winLoss := &WinLoss{
		Wins:   0,
		Losses: 0,
	}

	betList, betListErr := service.betService.GetBetsFromUser(userID)
	if betListErr != nil {
		return nil, fmt.Errorf("failed to get bets for user %s: %w", userID, betListErr)
	}

	for _, bet := range betList {
		switch bet.GetBetStatus() {
		case bets.Won:
			winLoss.Wins++
		case bets.Lost:
			winLoss.Losses++
		case bets.Pending:
			continue
		}
	}

	return winLoss, nil
}

var _ UserService = (*service)(nil)
