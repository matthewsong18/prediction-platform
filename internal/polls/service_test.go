package polls

import "testing"

func setupService(t *testing.T) (PollService, func()) {
	t.Helper()

	repo := NewMemoryRepository()
	service := NewService(repo)
	teardown := func() {}

	return service, teardown
}

func TestServiceImplementations(t *testing.T) {
	t.Parallel()
	implementations := []struct {
		name  string
		setup func(t *testing.T) (PollService, func())
	}{
		{name: "service", setup: setupService},
	}

	testCases := []struct {
		name string
		run  func(t *testing.T, service PollService)
	}{
		{"it should create a poll", testCreatePoll},
		{"it should close a poll", testClosePoll},
		{"it should select an outcome", testSelectOutcome},
		{"it should get a poll by ID", testGetPollById},
		{"it should return an error for more than two options", testExactlyTwoOptions},
		{"it should return all open polls", testGetAllOpen},
	}

	for _, implementation := range implementations {
		t.Run(implementation.name, func(t *testing.T) {
			t.Parallel()
			for _, testCase := range testCases {
				t.Run(testCase.name, func(t *testing.T) {
					t.Parallel()
					service, cleanup := implementation.setup(t)
					t.Cleanup(cleanup)

					testCase.run(t, service)
				})
			}
		})
	}
}

func testGetAllOpen(t *testing.T, pollService PollService) {
	// ARRANGE: Create open and closed polls
	if _, err := pollService.CreatePoll("openPoll1", []string{"option1", "option2"}); err != nil {
		t.Fatalf("Failed to create open poll: %v", err)
	}
	if _, err := pollService.CreatePoll("openPoll2", []string{"option1", "option2"}); err != nil {
		t.Fatalf("Failed to create open poll: %v", err)
	}
	closedPoll, err := pollService.CreatePoll("closedPoll", []string{"option1", "option2"})
	if err != nil {
		t.Fatalf("Failed to create closed poll: %v", err)
	}
	err = pollService.ClosePoll(closedPoll.GetID())
	if err != nil {
		t.Fatalf("Failed to close poll: %v", err)
	}

	// ACT: Get all open polls
	polls, err := pollService.GetOpenPolls()
	if err != nil {
		t.Fatalf("GetOpenPolls returned an unexpected error: %v", err)
	}

	// ASSERT: Verify the number of open polls
	if len(polls) != 2 {
		t.Errorf("Expected 2 open polls, but got %d", len(polls))
	}

}

func testCreatePoll(t *testing.T, service PollService) {
	title := "Which team will win first map?"
	options := []string{"Team A", "Team B"}

	poll, err := service.CreatePoll(title, options)

	if err != nil {
		t.Fatalf("CreatePoll returned an unexpected error: %v", err)
	}

	if poll.GetID() == "" {
		t.Error("Expected poll ID to be set, but it was empty")
	}

	if poll.GetTitle() != title {
		t.Errorf("Expected poll title to be '%s', but got '%s'", title, poll.GetTitle())
	}

	for i, option := range options {
		if poll.GetOptions()[i] != option {
			t.Errorf("Expected option %d to be '%s', but got '%s'", i, option, poll.GetOptions()[i])
		}
	}

	if poll.GetStatus() != Open {
		t.Error("Expected poll to be open, but it was closed")
	}

}

func testExactlyTwoOptions(t *testing.T, service PollService) {
	title := "Which team will win first map?"
	options := []string{"Team A", "Team B", "Team C"}

	_, err := service.CreatePoll(title, options)

	if err == nil {
		t.Fatal("Expected CreatePoll to return an error for more than two options, but it did not")
	}

	expectedError := "poll must have exactly two options"
	if err.Error() != expectedError {
		t.Errorf("Expected error message '%s', but got '%s'", expectedError, err.Error())
	}
}

func testClosePoll(t *testing.T, service PollService) {
	poll, err := createDefaultTestPoll(service)

	if err != nil {
		t.Fatal("CreatePoll returned an unexpected error:", err)
	}
	if poll.GetStatus() != Open {
		t.Fatal("Expected poll to be open after creation, but it was closed")
	}

	if err := service.ClosePoll(poll.GetID()); err != nil {
		t.Fatal("ClosePoll returned an unexpected error:", err)
	}

	updatedPoll, updateError := service.GetPollById(poll.GetID())
	if updateError != nil {
		t.Fatal("GetPollById returned an unexpected error:", updateError)
	}

	if updatedPoll.GetStatus() != Closed {
		t.Error("Expected poll to be closed after ClosePoll, but it was still open")
	}
}

func testSelectOutcome(t *testing.T, service PollService) {
	poll, err := createDefaultTestPoll(service)
	if err != nil {
		t.Fatal("CreatePoll returned an unexpected error", err)
	}

	// Test selecting an outcome
	teamAIndex := Option1
	err = service.SelectOutcome(poll.GetID(), teamAIndex)
	if err != nil {
		t.Fatal("SelectOutcome returned an unexpected error", err)
	}

	// Get the updated poll
	poll, err = service.GetPollById(poll.GetID())
	if err != nil {
		t.Fatal("GetPollById returned an unexpected error:", err)
	}

	// Verify the outcome
	if poll.GetOutcome() != teamAIndex {
		t.Errorf("Expected selected outcome to be '%d', but got '%d'", teamAIndex, poll.GetOutcome())
	}
}

func createDefaultTestPoll(service PollService) (Poll, error) {
	title := "Which team will win first map?"
	options := []string{"Team A", "Team B"}
	poll, err := service.CreatePoll(title, options)
	return poll, err
}

func testGetPollById(t *testing.T, pollService PollService) {
	// Testing that the GetPollById method retrieves the exact poll that was created by CreatePoll instead of a copy.
	poll, err := createDefaultTestPoll(pollService)
	if err != nil {
		t.Fatal("CreatePoll returned an unexpected error:", err)
	}

	retrievedPoll, err := pollService.GetPollById(poll.GetID())
	if err != nil {
		t.Fatal("GetPollById returned an unexpected error:", err)
	}

	if retrievedPoll.GetID() != poll.GetID() {
		t.Errorf("Expected retrieved poll to be equal to created poll, but they differ")
	}
}
