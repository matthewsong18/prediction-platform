package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	"betting-discord-bot/internal/bets"
	"betting-discord-bot/internal/polls"
	"betting-discord-bot/internal/users"

	"github.com/bwmarrin/discordgo"
)

func handleBet(s *discordgo.Session, i *discordgo.InteractionCreate, bot *Bot, pollID string, user users.User, optionIndex int) {
	bet, betErr := bot.BetService.CreateBet(pollID, user.GetID(), optionIndex)
	if betErr != nil {
		if errors.Is(betErr, bets.ErrUserAlreadyBet) {
			if err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "You have already bet on this poll.",
					Flags:   discordgo.MessageFlagsEphemeral,
				},
			}); err != nil {
				log.Printf("Error sending invalid action response: %v", err)
			}
		}

		if errors.Is(betErr, bets.ErrPollIsClosed) {
			if err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "This poll is closed. You cannot place a bet.",
					Flags:   discordgo.MessageFlagsEphemeral,
				},
			}); err != nil {
				log.Printf("Error sending poll is closed response: %v", err)
			}

		}

		log.Printf("Error creating bet: %v", betErr)
		return
	}

	log.Printf("Bet created: %v", bet)

	if err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Bet submitted",
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	}); err != nil {
		log.Printf("Error sending bet confirmation: %v", err)
	}
}

func handleEndPoll(s *discordgo.Session, i *discordgo.InteractionCreate, bot *Bot, pollID string) {
	if doesNotHaveManageMemberPerm(s, i) {
		return
	}

	if err := bot.PollService.ClosePoll(pollID); err != nil {
		if errors.Is(err, polls.ErrPollIsAlreadyClosed) {
			log.Printf("Poll \"%s\" is already closed", pollID)

			sendInteractionResponse(s, i, "The poll is already closed")

			return
		}
		log.Printf("Error closing poll: %v", err)
		return
	}

	sendInteractionResponse(s, i, "The poll is closed")

	log.Printf("User %s ended poll %s", i.Member.User.GlobalName, pollID)
}

func (bot *Bot) handleSelectOutcomeButton(s *discordgo.Session, i *discordgo.InteractionCreate, pollID string) {
	if doesNotHaveManageMemberPerm(s, i) {
		return
	}

	poll, pollErr := bot.PollService.GetPollById(pollID)
	if pollErr != nil {
		log.Printf("Error getting poll: %v", pollErr)
		return
	}

	if poll.GetStatus() == polls.Open {
		sendInteractionResponse(s, i, "The poll is still open. You cannot select an outcome.")
	}

	textDisplay := NewTextDisplay("Choose the outcome of the poll")

	selectOutcomeDropdown := NewStringSelect(
		"Select An Outcome",
		1,
		1,
		fmt.Sprintf("select:%s", pollID),
		[]interface{}{
			&StringOption{
				Label:       poll.GetOptions()[0],
				Value:       "1",
				Description: "Option 1",
			},
			&StringOption{
				Label:       poll.GetOptions()[1],
				Value:       "2",
				Description: "Option 2",
			},
		},
	)

	actionRow := NewActionRow([]interface{}{selectOutcomeDropdown})

	messageContainer := NewContainer(
		0xe32458,
		[]interface{}{
			textDisplay,
			actionRow,
		},
	)

	const permissions = IsComponentsV2 | MessageIsEphemeral

	message := MessageSend{
		Flags: permissions,
		Components: []interface{}{
			messageContainer,
		},
	}

	response := NewInteractionResponse(ChannelMessageWithSource, message)

	jsonMessage, jsonErr := json.Marshal(response)
	if jsonErr != nil {
		log.Printf("Error marshaling selectOutcomeDropdown: %v", jsonErr)
		return
	}

	interactionID := i.ID
	interactionToken := i.Token
	url := fmt.Sprintf("https://discord.com/api/v10/interactions/%s/%s/callback", interactionID, interactionToken)
	sendHttpRequest(url, jsonMessage)
}

func (bot *Bot) handleSelectOutcomeDropdown(s *discordgo.Session, i *discordgo.InteractionCreate) {
	customID := i.MessageComponentData().CustomID
	messageData := strings.Split(customID, ":")
	pollID := messageData[1]
	optionIndex, err := strconv.Atoi(i.MessageComponentData().Values[0])
	if err != nil {
		log.Printf("Error parsing option index: %v", err)
		return
	}

	poll, pollErr := bot.PollService.GetPollById(pollID)
	if pollErr != nil {
		log.Printf("Error getting poll: %v", pollErr)
		return
	}

	var pollOutcome polls.OutcomeStatus
	switch optionIndex {
	case 1:
		pollOutcome = polls.Option1
	case 2:
		pollOutcome = polls.Option2
	default:
		log.Panicf("Invalid option index: %d", optionIndex)
	}

	if err := bot.PollService.SelectOutcome(pollID, pollOutcome); err != nil {
		log.Printf("Error selecting outcome: %v", err)
		return
	}

	poll, pollErr = bot.PollService.GetPollById(pollID)
	if pollErr != nil {
		log.Panicf("Error getting poll: %v", pollErr)
	}

	sendInteractionResponse(s, i, "The outcome of the poll has been selected.")

	messageString := NewTextDisplay(fmt.Sprintf(
		"Outcome for **%s** between **%s** and **%s** has been decided.\n\nThe outcome is **%s**.",
		poll.GetTitle(),
		poll.GetOptions()[0],
		poll.GetOptions()[1],
		poll.GetOptions()[poll.GetOutcome()],
	))

	messageContainer := NewContainer(
		0xe32458,
		[]interface{}{
			messageString,
		},
	)

	const permissions = IsComponentsV2
	messageSend := MessageSend{
		Flags: permissions,
		Components: []interface{}{
			messageContainer,
		},
	}

	jsonMessage, jsonErr := json.Marshal(messageSend)
	if jsonErr != nil {
		log.Panicf("Error marshaling selectOutcomeDropdown:%v", jsonErr)
	}

	url := createMessageAPI(i.ChannelID)
	sendHttpRequest(url, jsonMessage)
}
