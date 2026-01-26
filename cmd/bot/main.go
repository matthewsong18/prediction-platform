package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/signal"

	"betting-discord-bot/internal/bets"
	"betting-discord-bot/internal/cryptography"
	"betting-discord-bot/internal/polls"
	"betting-discord-bot/internal/storage"
	"betting-discord-bot/internal/users"
	"encoding/hex"

	"github.com/bwmarrin/discordgo"
)

func run() (err error) {
	// Validate ENV
	config, err := LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	// Initialize Discord Session
	discordSession, err := discordgo.New("Bot " + config.Token)
	if err != nil {
		return fmt.Errorf("invalid bot parameters: %w", err)
	}

	// Setup DB
	db, initDBError := storage.InitializeDatabase(config.DBPath, config.EncryptionKey)

	if initDBError != nil {
		return fmt.Errorf("failed to initialize database: %w", initDBError)
	}

	log.Println("Database initialized successfully")

	// Init services
	pollService, betService, userService, err := initServices(db, config)
	if err != nil {
		return fmt.Errorf("failed to initialize services: %w", err)
	}

	// Setup discord bot
	if err := setupDiscordBot(discordSession, config, pollService, betService, userService); err != nil {
		return fmt.Errorf("failed to setup discord bot: %w", err)
	}

	// Bot shutdown handlers
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop
	log.Println("Graceful shutdown")

	defer func(discordSession *discordgo.Session) {
		_ = discordSession.Close()
	}(discordSession)

	defer func() {
		if closeError := db.Close(); closeError != nil {
			fmt.Println("Error closing database", closeError)
			if err == nil {
				err = closeError
			}
		}
	}()

	return nil
}

func setupDiscordBot(discordSession *discordgo.Session, config *Config, pollService polls.PollService, betService bets.BetService, userService users.UserService) error {
	bot := NewBot(discordSession, pollService, betService, userService, config.AppID, config.GuildID)

	bot.DiscordSession.AddHandler(bot.interactionHandler)

	if err := bot.RegisterCommands(); err != nil {
		return fmt.Errorf("failed to register commands: %w", err)
	}

	discordSession.AddHandler(func(discordSession *discordgo.Session, ready *discordgo.Ready) {
		log.Println("Bot is up")
	})

	if err := discordSession.Open(); err != nil {
		return fmt.Errorf("cannot open the session: %w", err)
	}
	return nil
}

func initServices(db *sql.DB, config *Config) (polls.PollService, bets.BetService, users.UserService, error) {
	// Initialize cryptography service
	keyBytes, err := hex.DecodeString(config.EncryptionKey)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to decode encryption key")
	}
	var key [32]byte
	copy(key[:], keyBytes)

	cryptoService, err := cryptography.NewService(key)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to initialize crypto service: %w", err)
	}

	pollRepo := polls.NewLibSQLRepository(db)
	pollService := polls.NewService(pollRepo)
	betRepo := bets.NewLibSQLRepository(db)
	betService := bets.NewService(pollService, betRepo)
	userRepo := users.NewLibSQLRepository(db, cryptoService)
	userService := users.NewService(userRepo, betService)
	return pollService, betService, userService, nil
}

func main() {
	if err := run(); err != nil {
		log.Fatalf("application failed to start: %v", err)
	}
}
