# Tech Stack

## Core Technologies
- **Programming Language:** [Go 1.25+](https://go.dev/) - Chosen for its strong concurrency primitives, performance, and simplicity, which are well-suited for a high-performance bot backend.
- **Discord Integration:** [discordgo](https://github.com/bwmarrin/discordgo) - The established library for interacting with the Discord API in Go.
- **Persistence:** [LibSQL](https://github.com/tursodatabase/go-libsql) - A source-available fork of SQLite optimized for edge and distributed applications, providing local-first speed with future scaling options.

## Architecture
- **Hexagonal Architecture (Ports and Adapters):** Strictly separates the core business logic from external delivery mechanisms (Discord) and storage (LibSQL). This ensures the domain remains testable and portable to other platforms.

## Security & Privacy
- **Encryption at Rest:** Sensitive PII (Discord IDs, Usernames) is encrypted using **AES-256-GCM** with a random nonce (Go 1.24+ `cipher.NewGCMWithRandomNonce`).
- **Searchable Encrypted Data:** Employs **Blind Indexing** (HMAC-SHA256) to allow efficient lookups of encrypted identifiers without compromising privacy.
- **Key Management:** Encryption keys are managed via environment variables to ensure separation from the codebase.

## Development & Infrastructure
- **Environment:** Container-ready (implied by self-hosting goal) and optimized for low-resource environments.
- **Testing:** Standard Go `testing` package with a focus on Contract Testing for repository implementations to ensure interchangeable storage backends.
