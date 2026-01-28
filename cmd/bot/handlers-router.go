package main

import (
	"errors"
	"log"
	"strconv"
	"strings"

	"betting-discord-bot/internal/users"

	"github.com/bwmarrin/discordgo"
)

func (bot *Bot) interactionHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	switch i.Type {
	case discordgo.InteractionApplicationCommand:
		bot.handleSlashCommands(s, i)
	case discordgo.InteractionModalSubmit:
		// This is a modal submission
		bot.handleModals(s, i)
	case discordgo.InteractionMessageComponent:
		bot.handleMessageComponents(s, i)
	default:
		log.Printf("Unknown interaction type received: %v", i.Type)
	}
}

func (bot *Bot) handleSlashCommands(s *discordgo.Session, i *discordgo.InteractionCreate) {
	commandName := i.ApplicationCommandData().Name
	switch commandName {
	case "create-poll":
		bot.handleCreatePollCommand(s, i)
	default:
		log.Printf("Unknown slash command received: %s", commandName)
	}
}

func (bot *Bot) handleModals(s *discordgo.Session, i *discordgo.InteractionCreate) {
	customID := i.ModalSubmitData().CustomID
	switch customID {
	case "poll_modal":
		bot.handlePollModalSubmit(s, i)
	default:
		log.Printf("Unknown modal submission received: %s", customID)
	}
}

func (bot *Bot) handleMessageComponents(s *discordgo.Session, i *discordgo.InteractionCreate) {
	customID := i.MessageComponentData().CustomID
	messageData := strings.Split(customID, ":")
	switch messageData[0] {
	case "bet":
		log.Println("Routing bet interaction")
		bot.handleButtonPresses(s, i)
	case "select":
		log.Println("Routing select interaction")
		bot.handleSelectOutcomeDropdown(s, i)
	default:
		log.Printf("Unknown interaction type received: %v", messageData[0])
	}
}

func (bot *Bot) handleButtonPresses(s *discordgo.Session, i *discordgo.InteractionCreate) {
	customID := i.MessageComponentData().CustomID
	betData := strings.Split(customID, ":")

	if len(betData) != 3 {
		log.Printf("Invalid custom ID received: %s", customID)
		return
	}

	pollID := betData[1]
	optionIndex, err := strconv.Atoi(betData[2])
	if err != nil {
		log.Printf("failed to convert option index to int: %s", betData[2])
		return
	}

	userDiscordID := i.Member.User.ID

	identity := users.Identity{
		Provider:   "discord",
		ExternalID: userDiscordID,
	}

	user, getUserErr := bot.UserService.GetUserByExternalID(identity)
	if getUserErr != nil {
		if errors.Is(getUserErr, users.ErrUserNotFound) {
			var createUserErr error
			user, createUserErr = bot.UserService.CreateUser(identity)
			if createUserErr != nil {
				log.Printf("Error creating user: %v", createUserErr)
				return
			}
		} else {
			log.Printf("Error getting user: %v", getUserErr)
			return
		}
	}

	if optionIndex < 2 {
		handleBet(s, i, bot, pollID, user, optionIndex)
		return
	}

	if optionIndex == 2 {
		handleEndPoll(s, i, bot, pollID)
	}

	if optionIndex == 3 {
		bot.handleSelectOutcomeButton(s, i, pollID)
	}
}
