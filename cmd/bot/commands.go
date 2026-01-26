package main

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

func (bot *Bot) RegisterCommands() error {
	commands := []*discordgo.ApplicationCommand{
		{
			Name:        "create-poll",
			Description: "Create a new poll",
		},
	}

	_, err := bot.DiscordSession.ApplicationCommandBulkOverwrite(bot.AppID, bot.GuildID, commands)
	if err != nil {
		log.Printf("Error overwriting commands: %v", err)
		return err
	}

	log.Println("Commands successfully registered.")
	return nil
}