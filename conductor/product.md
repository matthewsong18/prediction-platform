# Initial Concept
A privacy-first, self-hosted Discord betting bot designed with a decoupled logic engine (Hexagonal Architecture) to allow for future expansion to other platforms. It enables communities to engage through friendly wagers on user-created polls and tracks performance via persistent statistics and leaderboards.

# Product Definition

## Target Audience
- **Private Gaming Communities:** Small to medium-sized groups looking for a private way to track friendly bets.
- **Large Public Discord Servers:** Communities that need a reliable, self-hosted alternative to centralized bots.
- **Multi-Platform Communities:** Groups that may eventually move or expand beyond Discord, requiring a transport-agnostic backend.

## Value Proposition
- **Privacy-First & Self-Hosted:** Users maintain control over their data and bot availability by hosting the backend themselves.
- **Clean Architectural Separation:** The core logic engine is strictly decoupled from the transport layer (Discord), making the system portable and resilient.
- **Gamified Engagement:** Enhances community interaction by providing a structured way to place bets, resolve outcomes, and compete on leaderboards.

## Core Features
- **Community-Driven Polls:** Users can create and manage their own polls for others to bet on.
- **Flexible Betting Markets:** Initial focus on binary (Yes/No) outcomes, powered by a robust domain model.
- **Statistics & Leaderboards:** Persistent tracking of user performance (wins, losses, accuracy) to drive server engagement.
- **Platform Agnosticism:** A design that prioritizes porting to other communication platforms (e.g., Slack, Telegram, or Web) in the future.
