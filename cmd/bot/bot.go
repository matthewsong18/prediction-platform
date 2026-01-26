package main

import (
	"betting-discord-bot/internal/bets"
	"betting-discord-bot/internal/polls"
	"betting-discord-bot/internal/users"
	"github.com/bwmarrin/discordgo"
)

type Bot struct {
	DiscordSession *discordgo.Session
	PollService    polls.PollService
	BetService     bets.BetService
	UserService    users.UserService
	AppID          string
	GuildID        string
}

func NewBot(session *discordgo.Session, pollService polls.PollService, betService bets.BetService, userService users.UserService, appID, guildID string) *Bot {
	return &Bot{
		DiscordSession: session,
		PollService:    pollService,
		BetService:     betService,
		UserService:    userService,
		AppID:          appID,
		GuildID:        guildID,
	}
}