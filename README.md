# Technical Design Document: Scalable Prediction Backend

![Go Version](https://img.shields.io/badge/go-1.22%2B-blue)
![Architecture](https://img.shields.io/badge/architecture-hexagonal-orange)
![Status](https://img.shields.io/badge/status-wip-yellow)

## 1. Context

**The Problem:** Running a simple poll or friendly wager inside a community
currently requires using centralized bots that are often expensive, unreliable,
or invasive with user data.

**The Solution:** A privacy-first, self-hosted backend that strips away the
bloat. It provides a clean logic engine for tracking bets and outcomes, allowing
communities to plug it into Discord (or any other app) without relying on a
third-party vendor.

## 2. System Architecture

The system implements Hexagonal Architecture (Ports and Adapters) to strictly
enforce separation of concerns between the core business domain and external
infrastructure.

### 2.1. Architectural Pattern: Ports & Adapters

The application core is isolated from external concerns (transport, persistence)
via interface-based contracts.

- **Primary Ports (Inbound):** Strictly defines the domain boundary. The Core
  has zero knowledge of the Discord implementations.
- **Secondary Ports (Outbound):** Define the dependencies required by the
  domain.
  - Examples: `Save(bet)`, `GetOpenPolls()`.
- **Adapters:**
  - **Driving Adapter:** The `cmd/bot` package functions as a Discord-specific
    implementation, translating Discord Interaction Events into domain commands.
  - **Driven Adapter:** Each domain package (`internal/bets`, `internal/polls`,
    `internal/users`) contains its own specific LibSQL implementation of the
    repository interface, handling persistence to SQLite/LibSQL.

### 2.2. Design Justification & Trade-offs

| Decision                  | Justification                                                                                                            |
| :------------------------ | :----------------------------------------------------------------------------------------------------------------------- |
| **Interface Agnosticism** | Ensures the backend can support multiple simultaneous frontends (e.g., Web, CLI) without modification to business logic. |
| **Dependency Injection**  | Facilitates rigorous unit testing and runtime adapter swapping (e.g., switching storage backends).                       |

### 2.3. Component Diagram

WIP placeholder

## 3. Testing Strategy

The quality assurance strategy prioritizes Functional Correctness and
Refactoring Resilience.

### 3.1. Blackbox Functional Testing

Tests operate strictly against the public API of the domain services. Internal
state and private methods are not mocked or inspected.

- **Objective:** Ensure that for a given input state, the system produces the
  compliant output state.
- **Benefit:** Enables aggressive refactoring of internal algorithms without
  triggering false-positive test failures.

### 3.2. Contract Testing (Storage)

To ensure interchangeable persistence layers, a shared compliance test suite
validates all repository implementations.

- **Mechanism:** Both `InMemoryRepository` (used for fast unit tests) and
  `LibSQLRepository` (production) run against the exact same test cases.
- **Guarantee:** Verifies that swapping the database driver preserves semantic
  correctness.

## 4. Development Status

**Current Status:** Active Development (WIP)

The core domain logic and Discord adapter are functional. Focus is currently on
stability and expanding the domain features.

## 5. Operations

### Prerequisites

- Go 1.22+
- Discord Bot Token
- LibSQL / SQLite compatible database file

### Quick Start

```bash
# 1. Clone the repository
git clone https://github.com/matthewsong18/prediction-platform.git

# 2. Configure Environment
export TOKEN="your_discord_token"
export GUILD_ID="your_guild_id"
export APP_ID="your_app_id"
export DB_PATH="./data/bets.db"

# 3. Run the application
go run ./cmd/bot
```
