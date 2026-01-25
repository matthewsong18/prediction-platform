# Implementation Plan: Privacy-Focused Data Encryption at Rest

This plan outlines the steps to implement transparent encryption for sensitive user data in the `users` repository using the `cryptography` service.

## Phase 1: Infrastructure & Configuration
- [ ] Task: Configure `ENCRYPTION_KEY` loading
    - [ ] Update `cmd/bot/main.go` or a configuration utility to read `ENCRYPTION_KEY` from environment variables.
    - [ ] Validate that the key is a valid 32-byte hex-encoded string.
    - [ ] Ensure the application fails fast with a clear error message if the key is missing or invalid.
- [ ] Task: Inject `cryptography.Service` into `users.LibSQLRepository`
    - [ ] Update `internal/users/interfaces.go` or factory methods to accept `cryptography.CryptoService`.
    - [ ] Update dependency injection in `cmd/bot/main.go` to wire up the crypto service to the user repository.
- [ ] Task: Conductor - User Manual Verification 'Phase 1: Infrastructure & Configuration' (Protocol in workflow.md)

## Phase 2: Transparent Encryption Logic
- [ ] Task: Implement Encryption/Decryption in `LibSQLRepository`
    - [ ] **Red Phase:** Write failing tests in `internal/users/libsql_repository_test.go` that check if data is encrypted in the database but returned as plaintext.
    - [ ] **Green Phase:** Modify `Save` and `Find` methods in `internal/users/libsql_repository.go` to use `cryptoService.Encrypt` and `cryptoService.Decrypt` for `DiscordID`, `Username`, and `DisplayName`.
    - [ ] Verify that internal domain logic remains unchanged and tests still pass.
- [ ] Task: Conductor - User Manual Verification 'Phase 2: Transparent Encryption Logic' (Protocol in workflow.md)

## Phase 3: Migration & Verification
- [ ] Task: Verify data at rest security
    - [ ] Write a script or temporary test that directly queries the LibSQL database (bypassing the repository's decryption) to confirm that PII fields contain ciphertext.
- [ ] Task: Conductor - User Manual Verification 'Phase 3: Migration & Verification' (Protocol in workflow.md)
