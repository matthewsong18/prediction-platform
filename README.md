# Discord Betting Bot

A stateful Discord bot that allows users to create polls, place bets on
outcomes, and track their win-loss records on a server leaderboard.

Written in Go, this project implements **Hexagonal Architecture** (Ports and
Adapters) to ensure the core betting logic remains completely decoupled from the
Discord API and the storage layer.

## Features

- **Poll Management**: Admins can create multi-option polls using interactive
  modals.
- **Betting System**: Users place bets on open polls via dynamic dropdown menus.
- **Stat Tracking**: The system calculates and persists user win-loss ratios and
  leaderboard standings.

## Architecture

- **Domain-Driven Design**: All business logic resides in `internal/`, isolated
  from external dependencies. The core logic does not know that Discord or
  SQLite exists.
- **Dependency Injection**: Services are injected via interfaces. This allows
  the application to swap seamlessly between storage backends (e.g., In-Memory
  for unit tests, LibSQL for production) without changing business code.
- **Event-Driven Router**: A custom implementation of the Discord gateway that
  routes slash commands (`/bet`) and modal submissions (`create-poll`) to
  specific handlers.

## Engineering Highlights

- **Contract Testing**: Implemented a shared behavioral test suite. This
  guarantees that the In-Memory repository and the LibSQL repository adhere to
  the exact same contract, preventing regression when switching storage drivers.
- **Sentinel Errors**: Utilizes public, package-level error variables (e.g.,
  `ErrPollClosed`) for robust control flow, avoiding brittle string comparisons.
- **LibSQL Integration**: Integrates the LibSQL via CGo bindings for persistent
  storage. Includes custom logic for handling connection string quirks (`Opaque`
  paths) and atomic schema initialization.

## Getting Started

### Prerequisites

- Go 1.22+
- A Discord Bot Token
- (Optional) LibSQL encryption key

### Installation

1. Clone the repository

```bash
git clone https://github.com/yourusername/betting-discord-bot.git
```

2. Set environment variables

```bash
export DISCORD_TOKEN="your_token_here"
export DB_PATH="data/bets.db"
```

3. Run the bot

```bash
go run ./cmd/bot
```

4. Run tests

```bash
go test -v ./...
```
