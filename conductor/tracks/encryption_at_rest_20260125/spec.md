# Specification: Privacy-Focused Data Encryption at Rest

## Overview
This track focuses on enhancing the privacy of the betting bot by ensuring that sensitive user information (Personally Identifiable Information - PII) is encrypted before being stored in the database. Following the project's "Privacy-First" principle, this implementation will focus on encrypting data at rest while maintaining a transparent plaintext experience for the application's core logic and Discord UI.

## Functional Requirements
- **Transparent Persistence Encryption:** The `users.Repository` implementation (specifically the `LibSQLRepository`) must encrypt sensitive fields before writing to the database and decrypt them upon retrieval.
- **Encrypted Fields:** The following fields in the `users` domain must be encrypted:
    - Discord User IDs
    - Usernames and Display Names
- **Key Management:** The encryption key must be a 32-byte value loaded from the `ENCRYPTION_KEY` environment variable at application startup.
- **Reversibility:** Encryption must be reversible (AES-GCM) to allow the bot to display the correct names and link back to the correct Discord accounts when needed (e.g., for leaderboards).

## Non-Functional Requirements
- **Performance:** Encryption/decryption operations should have minimal impact on the overall response time of the bot.
- **Security:** Use the `internal/cryptography` service which utilizes Go 1.24+'s `cipher.NewGCMWithRandomNonce` for secure, authenticated encryption.
- **Privacy:** Sensitive PII must not appear in plaintext in database backups or logs.

## Acceptance Criteria
- [ ] New `users.Repository` decorator or updated `LibSQLRepository` that utilizes `cryptography.Service`.
- [ ] User IDs and names are stored as encrypted blobs/strings in the LibSQL database.
- [ ] The Discord bot displays usernames and processes bets correctly using the transparent decryption layer.
- [ ] The application fails to start if `ENCRYPTION_KEY` is not provided or is invalid.
- [ ] Unit tests in `internal/users` verify that data saved is indeed encrypted in the underlying storage.

## Out of Scope
- End-to-end encryption (displaying encrypted IDs in the Discord UI).
- Encryption of non-PII betting data (e.g., bet amounts).
- Automated key rotation.
