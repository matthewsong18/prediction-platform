package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/signal"

	"betting-discord-bot/internal/bets"
	"betting-discord-bot/internal/polls"
	"betting-discord-bot/internal/storage"
	"betting-discord-bot/internal/users"

	"github.com/bwmarrin/discordgo"
)

// Bot parameters
var (
	GuildID       = os.Getenv("GUILD_ID")
	Token         = os.Getenv("TOKEN")
	AppID         = os.Getenv("APP_ID")
	DBPath        = os.Getenv("DB_PATH")
	EncryptionKey = os.Getenv("ENCRYPTION_KEY")
)

var discordSession *discordgo.Session

func init() {
	var err error
	discordSession, err = discordgo.New("Bot " + Token)
	if err != nil {
		log.Fatalf("Invalid bot parameters: %v", err)
	}
}

func run() (err error) {
	// Validate ENV
	if err := validateEnvVariables(); err != nil {
		return fmt.Errorf("failed to validate env variables: %w", err)
	}

	// Setup DB
	db, initDBError := storage.InitializeDatabase(DBPath, EncryptionKey)

	if initDBError != nil {
		return fmt.Errorf("failed to initialize database: %w", initDBError)
	}

	log.Println("Database initialized successfully")

	// Init services
	pollService, betService, userService := initServices(db)

	// Setup discord bot
	if err := setupDiscordBot(pollService, betService, userService); err != nil {
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

func setupDiscordBot(pollService polls.PollService, betService bets.BetService, userService users.UserService) error {
	bot := NewBot(discordSession, pollService, betService, userService)

	bot.DiscordSession.AddHandler(bot.interactionHandler)

	if err := bot.RegisterCommands(); err != nil {
		return fmt.Errorf("failed to register commands: %w", err)
	}

	discordSession.AddHandler(func(discordSession *discordgo.Session, ready *discordgo.Ready) {
		log.Println("Bot is up")
	})

	if err := discordSession.Open(); err != nil {
		log.Fatalf("Cannot open the session: %v", err)
	}
	return nil
}

func initServices(db *sql.DB) (polls.PollService, bets.BetService, users.UserService) {
	pollRepo := polls.NewLibSQLRepository(db)
	pollService := polls.NewService(pollRepo)
	betRepo := bets.NewLibSQLRepository(db)
	betService := bets.NewService(pollService, betRepo)
	userRepo := users.NewLibSQLRepository(db)
	userService := users.NewService(userRepo, betService)
	return pollService, betService, userService
}

func validateEnvVariables() error {
	flagErr := false
	if GuildID == "" {
		log.Println("GUILD_ID environment variable is not set")
		flagErr = true
	}
	if Token == "" {
		log.Println("TOKEN environment variable is not set")
		flagErr = true
	}
	if AppID == "" {
		log.Println("APP_ID environment variable is not set")
		flagErr = true
	}
	if DBPath == "" {
		log.Println("DB_PATH environment variable is not set")
		flagErr = true
	}

	if flagErr {
		return fmt.Errorf("invalid environment variables")
	}

	if EncryptionKey == "" {
		log.Println("ENCRYPTION_KEY environment variable is not set, using unencrypted database")
	}
	return nil
}

func main() {
	if err := run(); err != nil {
		log.Fatalf("application failed to start: %v", err)
	}
}
