package users

import (
	"errors"
	"testing"

	"betting-discord-bot/internal/bets"
	"betting-discord-bot/internal/polls"
)

func TestCreateUser(t *testing.T) {
	t.Parallel()
	pollMemoryRepo := polls.NewMemoryRepository()
	pollService := polls.NewService(pollMemoryRepo)
	betService := bets.NewService(pollService, nil)
	userRepo := NewMemoryRepository()
	userService := NewService(userRepo, betService)

	user, err := userService.CreateUser("discord", "12345")
	if err != nil {
		t.Fatalf("CreateUser returned an unexpected error: %v", err)
	}

	if user.GetID() == "" {
		t.Error("Expected ID to be set")
	}
}

func TestGetUserByExternalID(t *testing.T) {
	t.Parallel()
	userRepo := NewMemoryRepository()
	userService := NewService(userRepo, nil)

	user, err := userService.CreateUser("discord", "12345")
	if err != nil {
		t.Fatalf("CreateUser returned an unexpected error: %v", err)
	}

	retrievedUser, err := userService.GetUserByExternalID("discord", "12345")
	if err != nil {
		t.Fatalf("GetUserByExternalID returned an unexpected error: %v", err)
	}

	if retrievedUser.GetID() != user.GetID() {
		t.Errorf("Expected ID to be '%s', got '%s'", user.GetID(), retrievedUser.GetID())
	}
}

func TestDeleteUser(t *testing.T) {
	t.Parallel()
	userRepo := NewMemoryRepository()
	userService := NewService(userRepo, nil)

	_, err := userService.CreateUser("discord", "12345")
	if err != nil {
		t.Fatalf("CreateUser returned an unexpected error: %v", err)
	}

	err = userService.DeleteUser("discord", "12345")
	if err != nil {
		t.Fatalf("DeleteUser returned an unexpected error: %v", err)
	}

	_, err = userService.GetUserByExternalID("discord", "12345")
	if err == nil {
		t.Fatalf("Expected GetUserByExternalID to return an error after deletion")
	}

	if errors.Is(err, ErrUserNotFound) {
		t.Fatalf("Expected GetUserByExternalID to return ErrUserNotFound after deletion")
	}
}
